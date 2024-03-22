package ssmwrap

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type FileTargets []FileTarget

func (ts *FileTargets) String() string {
	s := ""

	for _, t := range *ts {
		s += t.String()
	}

	return s
}

func (ts *FileTargets) Set(value string) error {
	target, err := ts.parseFlag(value)
	if err != nil {
		return err
	}

	if err := target.IsSatisfied(); err != nil {
		return fmt.Errorf("file parameter is not satisfied: %s", err)
	}

	*ts = append(*ts, *target)

	return nil
}

func (ts FileTargets) parseFlag(value string) (*FileTarget, error) {
	target := &FileTarget{}
	parts := strings.Split(value, ",")

	for _, part := range parts {
		param := strings.Split(part, "=")
		if len(param) != 2 {
			return nil, fmt.Errorf("invalid format")
		}

		key := param[0]
		value := param[1]

		switch key {
		case "Name":
			target.Name = value
		case "Path":
			path, err := ts.checkPath(value)
			if err != nil {
				return nil, fmt.Errorf("invalid Path: %s", err)
			}
			target.Path = path
		case "Mode":
			mode, err := ts.checkMode(value)
			if err != nil {
				return nil, fmt.Errorf("invalid Mode: %s", err)
			}
			target.Mode = mode
		case "Uid":
			uid, err := ts.checkUid(value)
			if err != nil {
				return nil, fmt.Errorf("invalid Uid: %s", err)
			}
			target.Uid = uid
		case "Gid":
			gid, err := ts.checkGid(value)
			if err != nil {
				return nil, fmt.Errorf("invalid Gid: %s", err)
			}
			target.Gid = gid
		default:
			return nil, fmt.Errorf("unknown parameter: %s", key)
		}
	}

	return target, nil
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

func cliFlagViaEnv(prefix string, multiple bool) []string {
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

// RunCLI runs ssmwrap as a CLI, returns exit code.
func RunCLI(version string) int {
	var (
		paths           string
		names           string
		recursiveFlag   bool
		noRecursiveFlag bool
		retries         int

		envOutputFlag    bool
		envNoOutputFlag  bool
		envPrefix        string
		envUseEntirePath bool

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
	flag.BoolVar(&envUseEntirePath, "env-entire-path", false, "use entire parameter path for name of environment variables\ndisabled: /path/to/value -> VALUE\nenabled: /path/to/value -> PATH_TO_VALUE")
	flag.StringVar(&envPrefix, "prefix", "", "alias for -env-prefix")

	flag.Var(&fileTargets, "file", "write values as file\nformat:  Name=VALUE_NAME,Path=FILE_PATH,Mode=FILE_MODE,Gid=FILE_GROUP_ID,Uid=FILE_USER_ID\nexample: Name=/foo/bar,Path=/etc/bar,Mode=600,Gid=123,Uid=456")

	flag.BoolVar(&versionFlag, "version", false, "display version")
	flag.VisitAll(func(f *flag.Flag) {
		// read flag values also from environment variable e.g. SSMWRAP_PATHS

		multiple := false

		if f.Name == "file" {
			multiple = true
		}

		for _, value := range cliFlagViaEnv("SSMWRAP_"+strings.ToUpper(f.Name), multiple) {
			f.Value.Set(value)
		}
	})
	flag.Parse()

	if versionFlag {
		fmt.Printf("ssmwrap v%s\n", version)
		return 0
	}

	options := RunOptions{
		Recursive: !noRecursiveFlag,
		Retries:   retries,
		Command:   flag.Args(),
	}
	if len(options.Command) == 0 {
		fmt.Fprintln(os.Stderr, "command required in arguments")
		return 2
	}

	if paths != "" {
		options.Paths = strings.Split(paths, ",")
	}
	if names != "" {
		options.Names = strings.Split(names, ",")
	}

	ssm := DefaultSSMConnector{}
	dests := []Destination{}

	if !envNoOutputFlag {
		dests = append(dests, DestinationEnv{
			Prefix:        envPrefix,
			UseEntirePath: envUseEntirePath,
		})
	}

	if 0 < len(fileTargets) {
		dests = append(dests, DestinationFile{
			Targets: fileTargets,
		})
		for _, t := range fileTargets {
			options.Names = append(options.Names, t.Name)
		}
	}

	ctx := context.TODO()

	if err := Run(ctx, options, ssm, dests); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return 1
	}

	return 0
}
