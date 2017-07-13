Feature: CF Login tool
  I want a tool that will allow me to switch between CF environments with a
  single command. We already have this for ComicRelief, but this would be a
  good thing to generalise.

  Scenario: I can log into a CF, setting org and space
    Given I have a YML file "foo.yml":
      """
        foo:
          target: api.run.pivotal.io
          username: adrian.fedoreanu@armakuni.com
          password: Pennies!20
          org: adrian-fedoreanu-armakuni
          space: development
      """
    When I run "cf aklogin -f foo.yml foo"
    Then I should be logged into "api.run.pivotal.io" CF as "adrian.fedoreanu@armakuni.com"
    And my selected org/space should be "adrian-fedoreanu-armakuni"/"development"

  Scenario: I can log into a CF, with no org or space
    Given I have a YML file "foo.yml":
      """
        foo:
          target: api.run.pivotal.io
          username: adrian.fedoreanu@armakuni.com
          password: Pennies!20
      """
    When I run "cf aklogin -f foo.yml foo"
    Then I should be logged into "api.run.pivotal.io" CF as "adrian.fedoreanu@armakuni.com"
    And my selected org/space should auto-assigned

  Scenario: I can log into a CF with a global YML file
    Given I have a YML file "~/.cflogin.yml":
      """
        foo:
          target: api.run.pivotal.io
          username: adrian.fedoreanu@armakuni.com
          password: Pennies!20
      """
    When I run "cf aklogin foo"
    Then I should be logged into "api.run.pivotal.io" CF as "adrian.fedoreanu@armakuni.com"
    And my selected org/space should auto-assigned
