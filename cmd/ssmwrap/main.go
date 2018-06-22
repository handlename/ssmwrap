package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/handlename/ssmwrap"
)

func main() {
	var paths string

	flag.StringVar(&paths, "paths", "/", "comma separated parameter paths")
	flag.Parse()

	if err := ssmwrap.Run(strings.Split(paths, ","), flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
