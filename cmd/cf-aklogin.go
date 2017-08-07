package main

import (
	"code.cloudfoundry.org/cli/plugin"

	"github.com/armakuni/cf-aklogin"
)

func main() {
	plugin.Start(new(aklogin.CFPlugin))
}
