package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"code.cloudfoundry.org/cli/cf/flags"
	"code.cloudfoundry.org/cli/plugin"

	"github.com/olebedev/config"
)

const defaultYML = "~/.cflogin.yml"

type akLoginPlugin struct{}

// Profile matches a YML profile
type Profile struct {
	Target, Username, Password, Org, Space string
}

func (ak *akLoginPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "aklogin",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 2,
			Build: 8,
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
						"-version":  "Print version",
					},
				},
			},
		},
	}
}

func (ak *akLoginPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	switch args[0] {
	case "aklogin":
		fc, err := parseArguments(args)
		if check(err) {
			return
		}

		if fc.IsSet("version") {
			fmt.Printf("%d.%d.%d\n",
				ak.GetMetadata().Version.Major,
				ak.GetMetadata().Version.Minor,
				ak.GetMetadata().Version.Build)
			return
		}

		yml, err := globalYML(fc.String("filename"))
		if check(err) {
			return
		}

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
			fmt.Println("Please specify a profile.")
			return
		}
		fmt.Printf("Using profile: '%s'\n", profile)

		activeProfile, err := yml.Get(profile)
		if err != nil {
			fmt.Println("Profile not found.")
			return
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
	if err != nil {
		return nil, err
	}
	cfg, err := config.ParseYamlBytes(yml)
	if err != nil {
		return nil, err
	}

	include, err := cfg.Get("include")
	if err == nil {
		includes, _ := include.List("")
		for _, path := range includes {
			iyml, err := ioutil.ReadFile(expandTilde(path.(string)))
			if err != nil {
				return nil, err
			}
			yml = append(append(yml, 0x0a), iyml...) // 0x0a == "\n"
		}
	}
	return config.ParseYamlBytes(yml)
}

func login(cliConn plugin.CliConnection, p *Profile) error {
	output, err := cliConn.CliCommandWithoutTerminalOutput("login",
		"-a", p.Target,
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
	fc.NewStringFlagWithDefault("filename", "f", "YML config file path", expandTilde(defaultYML))
	fc.NewBoolFlag("list", "l", "List available profiles")
	fc.NewBoolFlag("version", "v", "Print version")
	return fc, fc.Parse(args...)
}

func expandTilde(filename string) string {
	if filename[:2] == "~/" {
		filename = filepath.Join(os.Getenv("HOME"), filename[2:])
	}
	return filename
}

func check(err error) (ok bool) {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return
}

func main() {
	plugin.Start(new(akLoginPlugin))
}
