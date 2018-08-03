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
		retries     int
		versionFlag bool
	)

	flag.StringVar(&paths, "paths", "/", "comma separated parameter paths")
	flag.StringVar(&prefix, "prefix", "", "prefix for environment variables")
	flag.IntVar(&retries, "retries", 0, "number of times of retry")
	flag.BoolVar(&versionFlag, "version", false, "display version")
	flag.VisitAll(func(f *flag.Flag) {
		// read flag values also from environment variable e.g. SSMWRAP_PATHS
		if s := os.Getenv("SSMWRAP_" + strings.ToUpper(f.Name)); s != "" {
			f.Value.Set(s)
		}
	})
	flag.Parse()

	if versionFlag {
		fmt.Printf("ssmwrap v%s\n", version)
		os.Exit(0)
	}

	options := ssmwrap.Options{
		Paths:   strings.Split(paths, ","),
		Prefix:  prefix,
		Retries: retries,
		Command: flag.Args(),
	}

	if err := ssmwrap.Run(options); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
