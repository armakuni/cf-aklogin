package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

type feature struct {
	*received
}

type received struct {
	target, user, org, space string
	login                    bool
}

func (r *received) String() string {
	return fmt.Sprintf("login: %t, target: %s, user: %s, org: %s, space: %s",
		r.login, r.target, r.user, r.org, r.space)
}
func (f *feature) iHaveAYMLFile(filename string, contents *gherkin.DocString) error {
	if filename[:2] == "~/" {
		filename = filepath.Join(os.Getenv("HOME"), filename[2:])
	}
	return ioutil.WriteFile(filename, []byte(contents.Content), 0644)
}

func (f *feature) iRun(commands string) error {
	command := strings.Split(commands, " ")
	cmd := exec.Command(command[0], command[1:]...)
	bytes, err := cmd.Output()
	f.received, err = parseCmdOutput(bytes)
	return err
}

func (f *feature) iShouldBeLoggedIntoCFAs(target, username string) error {
	err := assertEquals(f.received.login, true)
	err = assertEquals(f.received.target, target)
	err = assertEquals(f.received.user, username)
	return err
}

func (f *feature) mySelectedOrgspaceShouldBeDevelopment(org, space string) error {
	err := assertEquals(f.received.org, org)
	err = assertEquals(f.received.space, space)
	return err
}

func (f *feature) mySelectedOrgspaceShouldAutoassigned() error {
	err := assertNotEquals(f.received.org, "")
	err = assertNotEquals(f.received.space, "")
	return err
}

func assertEquals(actual, expected interface{}) error {
	if expected != actual {
		return fmt.Errorf("Expected %s, but got %s", expected, actual)
	}
	return nil
}

func assertNotEquals(actual, expected interface{}) error {
	if expected == actual {
		return fmt.Errorf("Expected %s, but got %s", expected, actual)
	}
	return nil
}

func parseCmdOutput(b []byte) (*received, error) {
	s := string(b)
	if strings.Contains(s, "is not a registered command") {
		return nil, errors.New("Install the plugin first.")
	}

	return &received{
		login:  regexp.MustCompile(`(?m)^Authenticating.+\nOK$`).Match(b),
		target: regexp.MustCompile(`http?s://([\w.]+)`).FindStringSubmatch(s)[1],
		user:   regexp.MustCompile(`User:\s+([\w.@]+)`).FindStringSubmatch(s)[1],
		org:    regexp.MustCompile(`Org:\s+([\w.-]+)`).FindStringSubmatch(s)[1],
		space:  regexp.MustCompile(`Space:\s+([\w.-]+)`).FindStringSubmatch(s)[1],
	}, nil
}

func FeatureContext(s *godog.Suite) {
	f := new(feature)

	s.AfterScenario(func(interface{}, error) {
		os.Remove("foo.yml")
	})

	s.Step(`^I have a YML file "([^"]*)":$`, f.iHaveAYMLFile)
	s.Step(`^I run "([^"]*)"$`, f.iRun)
	s.Step(`^I should be logged into "([^"]*)" CF as "([^"]*)"$`, f.iShouldBeLoggedIntoCFAs)
	s.Step(`^my selected org\/space should be "([^"]*)"\/"([^"]*)"$`, f.mySelectedOrgspaceShouldBeDevelopment)
	s.Step(`^my selected org\/space should auto-assigned$`, f.mySelectedOrgspaceShouldAutoassigned)
}
