package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/handlename/ssmwrap"
)

func main() {
	var (
		paths  string
		prefix string
	)

	flag.StringVar(&paths, "paths", "/", "comma separated parameter paths")
	flag.StringVar(&prefix, "prefix", "", "prefix for environment variables")
	flag.Parse()

	if err := ssmwrap.Run(strings.Split(paths, ","), prefix, flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
