package aklogin

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"

	"code.cloudfoundry.org/cli/plugin/models"
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"code.cloudfoundry.org/cli/util/testhelpers/io"
)

type pluginFeature struct {
	output string

	fakeCliConnection *pluginfakes.FakeCliConnection
	CFPlugin          *CFPlugin
}

func (p *pluginFeature) iHaveAYMLFile(filename string, contents *gherkin.DocString) error {
	return ioutil.WriteFile(expandTilde(filename), []byte(contents.Content), 0644)
}

func (p *pluginFeature) iRunCf(command string) error {
	out := io.CaptureOutput(func() {
		p.CFPlugin.Run(p.fakeCliConnection, strings.Split(command, " "))
	})
	p.output = strings.Join(out, "\n")
	return nil
}

func (p *pluginFeature) iShouldBeLoggedIntoCFAs(target, username string) error {
	loggedIn, _ := p.fakeCliConnection.IsLoggedIn()
	return assertEq(loggedIn, true)
}

func (p *pluginFeature) mySelectedOrgspaceShouldBeDevelopment(org, space string) (err error) {
	currentOrg, _ := p.fakeCliConnection.GetCurrentOrg()
	err = assertEq(currentOrg.Name, org)
	if err != nil {
		return
	}
	currentSpace, _ := p.fakeCliConnection.GetCurrentSpace()
	return assertEq(currentSpace.Name, space)
}

func (p *pluginFeature) mySelectedOrgspaceShouldAutoassigned() (err error) {
	currentOrg, _ := p.fakeCliConnection.GetCurrentOrg()
	err = assertNotEq(currentOrg.Name, "")
	if err != nil {
		return
	}
	currentSpace, _ := p.fakeCliConnection.GetCurrentSpace()
	return assertNotEq(currentSpace.Name, "")
}

func (p *pluginFeature) theCFAKPluginIsInstalled() error {
	return nil // Plugin faked
}

func (p *pluginFeature) theOutputShouldBe(expected *gherkin.DocString) error {
	return assertEqStr(p.output, expected.Content)
}

func assertEq(got, exp interface{}) error {
	if !reflect.DeepEqual(got, exp) {
		return fmt.Errorf("Wanted '%v'; Got '%v'", exp, got)
	}
	return nil
}

func assertNotEq(got, exp interface{}) error {
	if exp == got {
		return fmt.Errorf("Wanted '%v'; Got '%v'", exp, got)
	}
	return nil
}

func assertEqStr(got, exp string) error {
	if exp != strings.TrimRight(got, "\n") {
		return fmt.Errorf("Wanted '%s', but got '%s'", exp, got)
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	p := &pluginFeature{
		fakeCliConnection: new(pluginfakes.FakeCliConnection),
		CFPlugin:          new(CFPlugin),
	}

	s.BeforeSuite(func() {
		p.fakeCliConnection.IsLoggedInReturns(true, nil)

		organization := plugin_models.Organization{
			OrganizationFields: plugin_models.OrganizationFields{Name: "adrian-fedoreanu-armakuni"}}
		p.fakeCliConnection.GetCurrentOrgReturns(organization, nil)

		space := plugin_models.Space{SpaceFields: plugin_models.SpaceFields{Name: "development"}}
		p.fakeCliConnection.GetCurrentSpaceReturns(space, nil)
	})

	s.AfterScenario(func(interface{}, error) {
		os.Remove(expandTilde("~/bar.yml"))
		os.Remove("foo.yml")
		os.Remove("invalid_foo.yml")
		os.Remove("fake_2.yml")
	})

	s.Step(`^I have a YML file "([^"]*)":$`, p.iHaveAYMLFile)
	s.Step(`^The cf-aklogin plugin is installed$`, p.theCFAKPluginIsInstalled)
	s.Step(`^I run cf "([^"]*)"$`, p.iRunCf)
	s.Step(`^I should be logged into "([^"]*)" CF as "([^"]*)"$`, p.iShouldBeLoggedIntoCFAs)
	s.Step(`^my selected org\/space should be "([^"]*)"\/"([^"]*)"$`, p.mySelectedOrgspaceShouldBeDevelopment)
	s.Step(`^my selected org\/space should auto-assigned$`, p.mySelectedOrgspaceShouldAutoassigned)
	s.Step(`^the output should be:$`, p.theOutputShouldBe)
}

func TestMain(m *testing.M) {
	format := "progress"
	for _, arg := range os.Args[1:] {
		if arg == "-test.v=true" { // go test transforms -v option
			format = "pretty"
			break
		}
	}
	status := godog.RunWithOptions("godog", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:    format,
		Paths:     []string{"features"},
		Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
