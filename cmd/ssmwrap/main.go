package main

import (
	"os"

	"github.com/handlename/ssmwrap"
)

var version string

const FlagEnvPrefix = "SSMWRAP_"

func main() {
	os.Exit(ssmwrap.RunCLI(version, FlagEnvPrefix))
}
