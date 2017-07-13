# CF Login tool

A tool that will allow switching between CF environments with a single command. 


### Install

[Download](https://bitbucket.org/armakuni/cf-aklogin/downloads/) the [plugin](https://bitbucket.org/armakuni/cf-aklogin/downloads/cf-aklogin) and run:
 
    $ cf install-plugin ~/Downloads/cf-aklogin
    
_Note: If you get `persmission denied`, run `chmod +x ~/Downloads/cf-aklogin`._

### Usage

Create `~/.cflogin.yml`:
    
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

    bar:
      target: api.run.pivotal.io
      username: <username>
      password: <password>
      org: <org> // optional
      space: <space> // optional

    $ cf aklogin -f foo.yml bar

### Tests

    $ godog

_Note: the plugin needs to be installed for the tests to succeed!_
