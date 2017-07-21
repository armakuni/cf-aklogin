package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/cloudfoundry/cli/cf/flags"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/olebedev/config"
)

var defaultYML string

type AKLoginPlugin struct{}

type Profile struct {
	Target, Username, Password, Org, Space string
}

func (ak *AKLoginPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "aklogin",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 2,
			Build: 6,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "aklogin",
				HelpText: "CF login via profiles",
				UsageDetails: plugin.Usage{
					Usage: "cf aklogin [options] <profile>",
					Options: map[string]string{
						"-filename": "YML config file path",
						"-list":     "List available profiles",
					},
				},
			},
		},
	}
}

func (ak *AKLoginPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	switch args[0] {
	case "aklogin":
		fc, err := parseArguments(args)
		check(err)

		yml, err := globalYML(fc.String("filename"))
		check(err)

		var profile string
		if len(fc.Args()) > 1 {
			profile = fc.Args()[1]
		}

		if fc.IsSet("list") {
			profilesMap, _ := yml.Map("")
			delete(profilesMap, "include")

			var inputProfile int
			profiles := make([]string, inputProfile, len(profilesMap))
			for p := range profilesMap {
				profiles = append(profiles, p)
			}
			sort.Strings(profiles)

			fmt.Println("Available profiles:")
			for i, p := range profiles {
				fmt.Printf("%d. %s\n", i, p)
			}

			fmt.Printf("Select profile: ")
			fmt.Scanf("%d", &inputProfile)

			if inputProfile > len(profiles) {
				inputProfile = 0
			}
			profile = profiles[inputProfile]
		}

		if profile == "" {
			exit1("Please specify a profile.")
		}
		fmt.Printf("Using profile: '%s'\n", profile)

		activeProfile, err := yml.Get(profile)
		if err != nil {
			exit1("Profile not found.")
		}

		target, err := activeProfile.String("target")
		check(err)

		username, err := activeProfile.String("username")
		check(err)

		// optional
		password, _ := activeProfile.String("password")
		org, _ := activeProfile.String("org")
		space, _ := activeProfile.String("space")

		p := &Profile{Target: target, Username: username, Password: password, Org: org, Space: space}

		err = login(cliConnection, p)
		check(err)

	case "CLI-MESSAGE-UNINSTALL":
		fmt.Println("Thanks for using the aklogin plugin.")
	}
}

func globalYML(filename string) (*config.Config, error) {
	yml, err := ioutil.ReadFile(filename)
	check(err)
	cfg, err := config.ParseYamlBytes(yml)
	check(err)

	include, err := cfg.Get("include")
	if err == nil {
		includes, _ := include.List("")
		for _, path := range includes {
			iyml, err := ioutil.ReadFile(expandTilde(path.(string)))
			if err != nil {
				fmt.Println(err)
				continue
			}
			yml = append(append(yml, 0x0a), iyml...) // 0x0a == "\n"
		}
	}
	return config.ParseYamlBytes(yml)
}

func login(cliConn plugin.CliConnection, p *Profile) error {
	output, err := cliConn.CliCommandWithoutTerminalOutput("login", "-a", p.Target,
		"-u", p.Username,
		"-p", p.Password,
		"-o", p.Org,
		"-s", p.Space)
	if err != nil {
		return err
	}
	for _, v := range output {
		fmt.Println(v)
	}
	return nil
}

func parseArguments(args []string) (flags.FlagContext, error) {
	fc := flags.New()
	fc.NewStringFlagWithDefault("filename", "f", "YML config file path", defaultYML)
	fc.NewBoolFlag("list", "l", "List available profiles")
	return fc, fc.Parse(args...)
}

func expandTilde(filename string) string {
	if filename[:2] == "~/" {
		filename = filepath.Join(os.Getenv("HOME"), filename[2:])
	}
	return filename
}

func check(err error) {
	if err != nil {
		exit1(err.Error())
	}
}

func exit1(err string) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	defaultYML = expandTilde("~/.cflogin.yml")
	plugin.Start(new(AKLoginPlugin))
}
