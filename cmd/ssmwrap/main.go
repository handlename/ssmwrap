package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/handlename/ssmwrap"
)

var version string

type FileTargets []ssmwrap.FileTarget

func (ts *FileTargets) String() string {
	s := ""

	for _, t := range *ts {
		s += t.String()
	}

	return s
}

func (ts *FileTargets) Set(value string) error {
	parsed := parseFileTarget(value)

	if parsed["Name"] == "" {
		return fmt.Errorf("Name required")
	}

	if parsed["Path"] == "" {
		return fmt.Errorf("Path required")
	}

	if parsed["Mode"] == "" {
		return fmt.Errorf("Mode required")
	}

	// expand path
	path, err := filepath.Abs(parsed["Path"])
	if err != nil {
		return fmt.Errorf("Invalid Path")
	}

	// convert `Mode` to os.FileMode
	mode, err := strconv.ParseUint(parsed["Mode"], 8, 32)
	if err != nil {
		return fmt.Errorf("invalid Mode")
	}

	target := ssmwrap.FileTarget{
		Name: parsed["Name"],
		Path: path,
		Mode: os.FileMode(mode),
	}

	if parsed["Uid"] != "" {
		uid, err := strconv.Atoi(parsed["Uid"])
		if err != nil {
			return fmt.Errorf("invalid Uid")
		}

		target.Uid = uid
	}

	if parsed["Gid"] != "" {
		gid, err := strconv.Atoi(parsed["Gid"])
		if err != nil {
			return fmt.Errorf("invalid Gid")
		}

		target.Gid = gid
	}

	*ts = append(*ts, target)

	return nil
}

func parseFileTarget(value string) map[string]string {
	parts := strings.Split(value, ",")
	parsed := map[string]string{}

	for _, part := range parts {
		param := strings.SplitN(part, "=", 2)
		parsed[param[0]] = param[1]
	}

	return parsed
}

func main() {
	var (
		paths   string
		retries int

		envOutputFlag   bool
		envNoOutputFlag bool
		envPrefix       string

		fileTargets FileTargets

		versionFlag bool
	)

	flag.StringVar(&paths, "paths", "/", "comma separated parameter paths")
	flag.IntVar(&retries, "retries", 0, "number of times of retry")

	flag.BoolVar(&envOutputFlag, "env", true, "export values as environment variables")
	flag.BoolVar(&envNoOutputFlag, "no-env", false, "disable export to environment variables")
	flag.StringVar(&envPrefix, "env-prefix", "", "prefix for environment variables")
	flag.StringVar(&envPrefix, "prefix", "", "alias for -env-prefix")

	flag.Var(&fileTargets, "file", "write values as file")

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

	options := ssmwrap.RunOptions{
		Paths:   strings.Split(paths, ","),
		Retries: retries,
		Command: flag.Args(),
	}

	ssm := ssmwrap.DefaultSSMConnector{}
	dests := []ssmwrap.Destination{}

	if !envNoOutputFlag {
		dests = append(dests, ssmwrap.DestinationEnv{
			Prefix: envPrefix,
		})
	}

	if 0 < len(fileTargets) {
		dests = append(dests, ssmwrap.DestinationFile{
			Targets: fileTargets,
		})
	}

	if err := ssmwrap.Run(options, ssm, dests); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
	}
}
