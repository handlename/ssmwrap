package main

import (
	"os"

	"github.com/handlename/ssmwrap/v2/cli"
)

const version = "2.1.0"

const FlagEnvPrefix = "SSMWRAP_"

func main() {
	os.Exit(int(cli.Run(version, FlagEnvPrefix)))
}
