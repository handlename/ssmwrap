package main

import (
	"os"

	"github.com/handlename/ssmwrap/v2/cli"
)

var version string

const FlagEnvPrefix = "SSMWRAP_"

func main() {
	os.Exit(int(cli.Run(version, FlagEnvPrefix)))
}
