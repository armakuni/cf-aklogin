package cflogin

import (
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

func iHaveAYAMLFile(arg1 string, arg2 *gherkin.DocString) error {
	return godog.ErrPending
}

func iRun(arg1 string) error {
	return godog.ErrPending
}

func iShouldBeLoggedIntoTheCFAs(arg1, arg2 string) error {
	return godog.ErrPending
}

func mySelectedOrgspaceShouldBeDevelopment(arg1, arg2 string) error {
	return godog.ErrPending
}

func mySelectedOrgspaceShouldBeUndetermined() error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I have a YAML file "([^"]*)":$`, iHaveAYAMLFile)
	s.Step(`^I run "([^"]*)"$`, iRun)
	s.Step(`^I should be logged into the "([^"]*)" CF as "([^"]*)"$`, iShouldBeLoggedIntoTheCFAs)
	s.Step(`^my selected org\/space should be "([^"]*)"\/"([^"]*)"$`, mySelectedOrgspaceShouldBeDevelopment)
	s.Step(`^my selected org\/space should be undetermined$`, mySelectedOrgspaceShouldBeUndetermined)
}
