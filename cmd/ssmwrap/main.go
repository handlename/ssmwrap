package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
	parts := strings.Split(value, ",")

	target := ssmwrap.FileTarget{}

	for _, part := range parts {
		param := strings.SplitN(part, "=", 2)
		key := param[0]
		value := param[1]

		switch key {
		case "Name":
			target.Name = value
		case "Path":
			path, err := ts.checkPath(value)
			if err != nil {
				return fmt.Errorf("invalid Path: %s", err)
			}
			target.Path = path
		case "Mode":
			mode, err := ts.checkMode(value)
			if err != nil {
				return fmt.Errorf("invalid Mode: %s", err)
			}
			target.Mode = mode
		case "Uid":
			uid, err := ts.checkUid(value)
			if err != nil {
				return fmt.Errorf("invalid Uid: %s", err)
			}
			target.Uid = uid
		case "Gid":
			gid, err := ts.checkGid(value)
			if err != nil {
				return fmt.Errorf("invalid Gid: %s", err)
			}
			target.Gid = gid
		default:
			return fmt.Errorf("unknown parameter: %s", key)
		}
	}

	err := target.IsSatisfied()
	if err != nil {
		return fmt.Errorf("file parameter is not satisfied: %s", err)
	}

	*ts = append(*ts, target)

	return nil
}

func (ts FileTargets) checkPath(value string) (string, error) {
	// expand path
	path, err := filepath.Abs(value)
	if err != nil {
		return "", fmt.Errorf("Invalid Path")
	}

	return path, nil
}

func (ts FileTargets) checkMode(value string) (os.FileMode, error) {
	// convert `Mode` to os.FileMode
	mode, err := strconv.ParseUint(value, 8, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid Mode")
	}

	return os.FileMode(mode), nil
}

func (ts FileTargets) checkGid(value string) (int, error) {
	gid, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid Gid")
	}

	return gid, nil
}

func (ts FileTargets) checkUid(value string) (int, error) {
	uid, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid Uid")
	}

	return uid, nil
}

func flagViaEnv(prefix string, multiple bool) []string {
	if multiple {
		values := []string{}

		// read values named like `SSMWRAP_FOO_123`
		re := regexp.MustCompile("^" + prefix + `(_\d+)?$`)

		for _, env := range os.Environ() {
			parts := strings.SplitN(env, "=", 2)

			if re.FindString(parts[0]) != "" {
				values = append(values, parts[1])
			}
		}

		return values
	}

	return []string{os.Getenv(prefix)}
}

func main() {
	var (
		paths           string
		names           string
		recursiveFlag   bool
		noRecursiveFlag bool
		retries         int

		envOutputFlag   bool
		envNoOutputFlag bool
		envPrefix       string

		fileTargets FileTargets

		versionFlag bool
	)

	flag.StringVar(&paths, "paths", "", "comma separated parameter paths")
	flag.StringVar(&names, "names", "", "comma separated parameter names")
	flag.BoolVar(&recursiveFlag, "recursive", true, "retrieve values recursively")
	flag.BoolVar(&noRecursiveFlag, "no-recursive", false, "retrieve values just under -paths only")
	flag.IntVar(&retries, "retries", 0, "number of times of retry")

	flag.BoolVar(&envOutputFlag, "env", true, "export values as environment variables")
	flag.BoolVar(&envNoOutputFlag, "no-env", false, "disable export to environment variables")
	flag.StringVar(&envPrefix, "env-prefix", "", "prefix for environment variables")
	flag.StringVar(&envPrefix, "prefix", "", "alias for -env-prefix")

	flag.Var(&fileTargets, "file", "write values as file\nformat:  Name=VALUE_NAME,Path=FILE_PATH,Mode=FILE_MODE,Gid=FILE_GROUP_ID,Uid=FILE_USER_ID\nexample: Name=/foo/bar,Path=/etc/bar,Mode=600,Gid=123,Uid=456")

	flag.BoolVar(&versionFlag, "version", false, "display version")
	flag.VisitAll(func(f *flag.Flag) {
		// read flag values also from environment variable e.g. SSMWRAP_PATHS

		multiple := false

		if f.Name == "file" {
			multiple = true
		}

		for _, value := range flagViaEnv("SSMWRAP_"+strings.ToUpper(f.Name), multiple) {
			f.Value.Set(value)
		}
	})
	flag.Parse()

	if versionFlag {
		fmt.Printf("ssmwrap v%s\n", version)
		os.Exit(0)
	}

	options := ssmwrap.RunOptions{
		Recursive: !noRecursiveFlag,
		Retries:   retries,
		Command:   flag.Args(),
	}
	if len(options.Command) == 0 {
		fmt.Fprintln(os.Stderr, "command required in arguments")
		os.Exit(2)
	}

	if paths != "" {
		options.Paths = strings.Split(paths, ",")
	}
	if names != "" {
		options.Names = strings.Split(names, ",")
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
		for _, t := range fileTargets {
			options.Names = append(options.Names, t.Name)
		}
	}

	if err := ssmwrap.Run(options, ssm, dests); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
}
