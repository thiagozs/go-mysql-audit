package main

import (
	"github.com/mercadobitcoin/go-proxy-audit/cmd"
)

var Version = "beta"
var Build = "dev"

func main() {
	cmd.Version = Version
	cmd.Build = Build
	cmd.Execute()
}
