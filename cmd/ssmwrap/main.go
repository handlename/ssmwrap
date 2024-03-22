package main

import (
	"os"

	"github.com/handlename/ssmwrap"
)

var version string

func main() {
	os.Exit(ssmwrap.RunCLI(version))
}
