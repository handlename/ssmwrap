package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/handlename/ssmwrap"
)

var version string

func main() {
	var (
		paths       string
		prefix      string
		versionFlag bool
	)

	flag.StringVar(&paths, "paths", "/", "comma separated parameter paths")
	flag.StringVar(&prefix, "prefix", "", "prefix for environment variables")
	flag.BoolVar(&versionFlag, "version", false, "display version")
	flag.Parse()

	if versionFlag {
		fmt.Printf("ssmwrap v%s\n", version)
		os.Exit(0)
	}

	if err := ssmwrap.Run(strings.Split(paths, ","), prefix, flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
