# CF Login tool

A tool that will allow switching between CF environments with a single command. 


### Install

[Download](https://bitbucket.org/armakuni/cf-aklogin/downloads/) the [plugin](https://bitbucket.org/armakuni/cf-aklogin/downloads/cf-aklogin) and run:
 
    $ cf install-plugin ~/Downloads/cf-aklogin
    
_Note: If you get `persmission denied`, run `chmod +x ~/Downloads/cf-aklogin`._

### Usage

#### Login

Create `~/.cflogin.yml`:
    
    include: //optional
    - ~/bar.yml
    foo:
      target: api.run.pivotal.io
      username: <username>
      password: <password>
      org: <org>
      space: <space>

_Note: leave password/org/space blank for `os.Stdin` input._

And then run:
    
    $ cf aklogin foo 

Or with your own `foo.yml`:

    include: //optional
    - ~/bar.yml
    foo:
      target: api.run.pivotal.io
      username: <username>
      password: <password>
      org: <org> // optional
      space: <space> // optional

    $ cf aklogin -f foo.yml bar

#### List

    $ cf aklogin -h
    NAME:
       aklogin - CF login via profiles
    
    USAGE:
       cf aklogin [options] <profile>
    
    OPTIONS:
       --filename       YML config file path
       --list           List available profiles

    $ cf aklogin --list
    Available profiles:
    0. ak
    1. bar    
    2. foo
    Select profile: _
        
### Tests

    $ godog

_Note: the plugin needs to be installed for the tests to succeed!_
