package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

type feature struct {
	*received
	cmdOutput []byte
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
	return ioutil.WriteFile(expandTilde(filename), []byte(contents.Content), 0644)
}

func (f *feature) iRun(commands string) error {
	command := strings.Split(commands, " ")
	cmd := exec.Command(command[0], command[1:]...)
	bytes, err := cmd.Output()
	f.cmdOutput = bytes
	return err
}

func (f *feature) iShouldBeLoggedIntoCFAs(target, username string) (err error) {
	f.received, err = parseCmdOutput(f.cmdOutput)
	err = assertEquals(f.received.login, true)
	err = assertEquals(f.received.target, target)
	err = assertEquals(f.received.user, username)
	return
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

func (f *feature) theOutputShouldBe(expected *gherkin.DocString) error {
	return assertEqualsStr(string(f.cmdOutput), expected.Content)
}

func assertEquals(actual, expected interface{}) error {
	if expected != actual {
		return fmt.Errorf("Expected '%s', but got '%s'", expected, actual)
	}
	return nil
}

func assertEqualsStr(actual, expected string) error {
	if expected != strings.TrimRight(actual, "\n") {
		return fmt.Errorf("Expected '%s', but got '%s'", expected, actual)
	}
	return nil
}

func assertNotEquals(actual, expected interface{}) error {
	if expected == actual {
		return fmt.Errorf("Expected '%s', but got '%s'", expected, actual)
	}
	return nil
}

func parseCmdOutput(b []byte) (*received, error) {
	s := string(b)
	if strings.Contains(s, "Profile not found.") {
		return nil, errors.New("Profile not found.")
	}

	if strings.Contains(s, "is not a registered command") {
		return nil, errors.New("Install the plugin first.")
	}

	return &received{
		login:  regexp.MustCompile(`(?m)^Authenticating.+\nOK$`).Match(b),
		target: extractGroupIfMatch(`http?s://([\w.]+)`, s),
		user:   extractGroupIfMatch(`User:\s+([\w.@]+)`, s),
		org:    extractGroupIfMatch(`Org:\s+([\w.-]+)`, s),
		space:  extractGroupIfMatch(`Space:\s+([\w.-]+)`, s),
	}, nil
}

func extractGroupIfMatch(regex, src string) string {
	matches := regexp.MustCompile(regex).FindStringSubmatch(src)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

func FeatureContext(s *godog.Suite) {
	f := new(feature)

	s.AfterScenario(func(interface{}, error) {
		os.Remove("foo.yml")
		os.Remove(expandTilde("~/bar.yml"))
	})

	s.Step(`^I have a YML file "([^"]*)":$`, f.iHaveAYMLFile)
	s.Step(`^I run "([^"]*)"$`, f.iRun)
	s.Step(`^I should be logged into "([^"]*)" CF as "([^"]*)"$`, f.iShouldBeLoggedIntoCFAs)
	s.Step(`^my selected org\/space should be "([^"]*)"\/"([^"]*)"$`, f.mySelectedOrgspaceShouldBeDevelopment)
	s.Step(`^my selected org\/space should auto-assigned$`, f.mySelectedOrgspaceShouldAutoassigned)
	s.Step(`^the output should be:$`, f.theOutputShouldBe)
}

func TestMain(m *testing.M) {
	status := godog.RunWithOptions("godog", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:    "progress",
		Paths:     []string{"features"},
		Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
