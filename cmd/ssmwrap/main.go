package main

import (
	"os"

	"github.com/handlename/ssmwrap/cli"
)

var version string

const FlagEnvPrefix = "SSMWRAP_"

func main() {
	os.Exit(cli.Run(version, FlagEnvPrefix))
}
