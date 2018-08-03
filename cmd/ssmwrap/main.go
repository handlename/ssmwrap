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
		paths   string
		retries int

		envOutputFlag   bool
		envNoOutputFlag bool
		envPrefix       string

		versionFlag bool
	)

	flag.StringVar(&paths, "paths", "/", "comma separated parameter paths")
	flag.IntVar(&retries, "retries", 0, "number of times of retry")

	flag.BoolVar(&envOutputFlag, "env", true, "export values as environment variables")
	flag.BoolVar(&envNoOutputFlag, "no-env", false, "disable export to environment variables")
	flag.StringVar(&envPrefix, "env-prefix", "", "prefix for environment variables")
	flag.StringVar(&envPrefix, "prefix", "", "alias for -env-prefix")

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
		Paths:     strings.Split(paths, ","),
		EnvOutput: !envNoOutputFlag,
		EnvPrefix: envPrefix,
		Retries:   retries,
		Command:   flag.Args(),
	}

	if err := ssmwrap.Run(options); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
