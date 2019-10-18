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
    When I run cf "aklogin -f foo.yml foo"
    Then I should be logged into "api.run.pivotal.io" CF as "adrian.fedoreanu@armakuni.com"
    And my selected org/space should be "adrian-fedoreanu-armakuni"/"development"

  Scenario: I can login with SSO
    Given I have a YML file "foo.yml":
      """
      foo:
        target: api.run.pivotal.io
        sso: true
      """
    When I run cf "aklogin -f foo.yml foo"
    Then I should be logged into "api.run.pivotal.io" CF as "adrian.fedoreanu@armakuni.com"
    And my selected org/space should auto-assigned

  Scenario: I can log into a CF, with no org or space
    Given I have a YML file "foo.yml":
      """
      foo:
        target: api.run.pivotal.io
        username: adrian.fedoreanu@armakuni.com
        password: Pennies!20
      """
    When I run cf "aklogin -f foo.yml foo"
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
    When I run cf "aklogin foo"
    Then I should be logged into "api.run.pivotal.io" CF as "adrian.fedoreanu@armakuni.com"
    And my selected org/space should auto-assigned

  Scenario: I can include multiple YML files into the global YML file
    and I can log in into a CF with sub-profiles
    Given I have a YML file "~/.cflogin.yml":
      """
    include:
    - ~/bar.yml
    foo:
      target: api.run.pivotal.io
      username: adrian.fedoreanu@armakuni.com
      password: Pennies!20
      """
    And I have a YML file "~/bar.yml":
      """
    bar:
      target: api.run.pivotal.io
      username: adrian.fedoreanu@armakuni.com
      password: Pennies!20
      """
    When I run cf "aklogin bar"
    Then I should be logged into "api.run.pivotal.io" CF as "adrian.fedoreanu@armakuni.com"
    And my selected org/space should auto-assigned

  Scenario: I can list all available profiles with includes
    Given I have a YML file "~/.cflogin.yml":
      """
      include:
      - ~/bar.yml
      foo:
        target: api.run.pivotal.io
        username: adrian.fedoreanu@armakuni.com
        password: Pennies!20
      """
    And I have a YML file "~/bar.yml":
      """
      bar:
        target: api.run.pivotal.io
        username: adrian.fedoreanu@armakuni.com
        password: Pennies!20
      """
    When I run cf "aklogin --list"
    Then the output should be:
      """
      Available profiles:
      0. bar
      1. foo
      Select profile: Using profile: 'bar'
      """

  Scenario: Invalid profile input
    Given I have a YML file "foo.yml":
      """
      """
    When I run cf "aklogin -f foo.yml -l"
    Then the output should be:
      """
    Available profiles:
    Select profile: Invalid profile.
      """

  Scenario: I can print the version
    Given The cf-aklogin plugin is installed
    When I run cf "aklogin --version"
    Then the output should be:
      """
      1.3.0
      """

  Scenario: Unspecified profile
    Given I have a YML file "foo.yml":
      """
      foo:
        target: api.run.pivotal.io
        username: adrian.fedoreanu@armakuni.com
        password: Pennies!20
      """
    When I run cf "aklogin -f foo.yml"
    Then the output should be:
      """
      Please specify a profile.
      """

  Scenario: Unspecified file
    Given I have a YML file "foo.yml":
      """
      foo:
        target: api.run.pivotal.io
        username: adrian.fedoreanu@armakuni.com
        password: Pennies!20
      """
    When I run cf "aklogin -f foo"
    Then the output should be:
      """
      open foo: no such file or directory
      """

  Scenario: Non-existing profile
    Given I have a YML file "foo.yml":
      """
      foo:
        target: api.run.pivotal.io
        username: adrian.fedoreanu@armakuni.com
        password: Pennies!20
      """
    When I run cf "aklogin -f foo.yml made-up-profile"
    Then the output should be:
      """
      Using profile: 'made-up-profile'
      Profile not found.
      """

  Scenario: Uninstall plugin
    Given The cf-aklogin plugin is installed
    When I run cf "CLI-MESSAGE-UNINSTALL"
    Then the output should be:
      """
      Thanks for using the aklogin plugin.
      """

  Scenario: Invalid argument
    Given The cf-aklogin plugin is installed
    When I run cf "aklogin -xx"
    Then the output should be:
      """
      Invalid flag: -xx
      """

  Scenario: Missing target
    Given I have a YML file "foo.yml":
      """
      foo:
        username: adrian.fedoreanu@armakuni.com
        password: Pennies!20
      """
    When I run cf "aklogin -f foo.yml foo"
    Then the output should be:
      """
      Using profile: 'foo'
      Nonexistent map key at "target"
      """

  Scenario: Missing username
    Given I have a YML file "foo.yml":
      """
      foo:
        target: api.run.pivotal.io
        password: Pennies!20
      """
    When I run cf "aklogin -f foo.yml foo"
    Then the output should be:
      """
      Using profile: 'foo'
      Nonexistent map key at "username"
      """

  Scenario: Invalid YML
    Given I have a YML file "invalid_foo.yml":
      """
      0x0a:
        target: api.run.pivotal.io
        username: adrian.fedoreanu@armakuni.com
        password: Pennies!20
      """
    When I run cf "aklogin -f invalid_foo.yml foo"
    Then the output should be:
      """
      Unsupported map key: 10
      """

  Scenario: Invalid include YML
    Given I have a YML file "foo.yml":
      """
      include:
      - fake.yml
      foo:
        target: api.run.pivotal.io
        username: adrian.fedoreanu@armakuni.com
        password: Pennies!20
      """
    And I have a YML file "fake_2.yml":
      """
      bar:
        target: api.run.pivotal.io
        username: adrian.fedoreanu@armakuni.com
        password: Pennies!20
      """
    When I run cf "aklogin -f foo.yml foo"
    Then the output should be:
      """
      open fake.yml: no such file or directory
      """
