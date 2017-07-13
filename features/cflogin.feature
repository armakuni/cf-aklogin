Feature: CF Login tool
  I want a tool that will allow me to switch between CF environments with a
  single command. We already have this for ComicRelief, but this would be a
  good thing to generalise.

Scenario: I can log into a CF, setting org and space
  Given I have a YAML file "foo.yml":
    """
      foo:
        target: run.pivotal.io
        username: adrian.fedoreanu@armakuni.com
        password: Pennies!20
        org: adrian-fedoreanu-armakuni
        space: development
    """
  When I run "cflogin -f foo.yml foo"
  Then I should be logged into the "run.pivotal.io" CF as "adrian.fedoreanu@armakuni.com"
  And my selected org/space should be "adrian-fedoreanu-armakuni"/"development"

Scenario: I can log into a CF, with no org or space
  Given I have a YAML file "foo.yml":
    """
      foo:
        target: run.pivotal.io
        username: adrian.fedoreanu@armakuni.com
        password: Pennies!20
    """
  When I run "cflogin -f foo.yml foo"
  Then I should be logged into the "run.pivotal.io" CF as "adrian.fedoreanu@armakuni.com"
  And my selected org/space should be undetermined

Scenario: I can log into a CF with a global YAML file
  Given I have a YAML file "~/.cflogin.yml":
    """
      foo:
        target: run.pivotal.io
        username: adrian.fedoreanu@armakuni.com
        password: Pennies!20
    """
  When I run "cflogin foo"
  Then I should be logged into the "run.pivotal.io" CF as "adrian.fedoreanu@armakuni.com"
  And my selected org/space should be undetermined
