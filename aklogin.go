package main

import (
	"fmt"
	"os"

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
			Minor: 0,
			Build: 0,
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
						"filename": "YML config file path",
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
		if err != nil {
			exit1(err.Error())
		}

		var profile string
		if len(fc.Args()) > 1 {
			profile = fc.Args()[1]
		}
		if profile == "" {
			exit1("Please specify a profile.")
		}
		fmt.Printf("Using profile: '%s'\n", profile)

		yml := fc.String("filename")
		cfg, err := config.ParseYamlFile(yml)
		if err != nil {
			exit1(err.Error())
		}

		activeConfig, err := cfg.Get(profile)
		if err != nil {
			exit1(fmt.Sprintf("Profile not found using '%s'.", yml))
		}

		target, err := activeConfig.String("target")
		if err != nil {
			exit1(err.Error())
		}

		username, err := activeConfig.String("username")
		if err != nil {
			exit1(err.Error())
		}

		// optional
		password, _ := activeConfig.String("password")
		org, _ := activeConfig.String("org")
		space, _ := activeConfig.String("space")

		p := &Profile{Target: target, Username: username, Password: password, Org: org, Space: space}

		err = login(cliConnection, p)
		if err != nil {
			exit1(err.Error())
		}

	case "CLI-MESSAGE-UNINSTALL":
		fmt.Println("Thanks for using the aklogin plugin.")
	}
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
	return fc, fc.Parse(args...)
}

func exit1(err string) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	defaultYML = fmt.Sprintf("%s/.cflogin.yml", os.Getenv("HOME"))
	plugin.Start(new(AKLoginPlugin))
}
