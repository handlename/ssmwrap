package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/handlename/ssmwrap/v2"
)

func main() {
	var (
		path          string
		prefix        string
		useentirepath bool
	)

	flag.StringVar(&path, "path", "", "path to export")
	flag.StringVar(&prefix, "prefix", "EXAMPLE_", "prefix for exported environment variable")
	flag.BoolVar(&useentirepath, "entirepath", false, "use entire path as environment variables name")
	flag.Parse()

	ctx := context.Background()

	rules := []ssmwrap.ExportRule{
		{
			Path:          path,
			Prefix:        prefix,
			UseEntirePath: useentirepath,
		},
	}

	if err := ssmwrap.Export(ctx, rules, ssmwrap.ExportOptions{}); err != nil {
		fmt.Fprintf(os.Stderr, "failed to export parameters: %v", err)
		os.Exit(1)
	}

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, prefix) {
			continue
		}

		fmt.Println(env)
	}
}
