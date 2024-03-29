package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/handlename/ssmwrap"
)

func flagViaEnv(prefix string, multiple bool) []string {
	if multiple {
		values := []string{}

		// read values named like `{prefix}_FOO_123`
		re := regexp.MustCompile("^" + prefix + `(_\d+)?$`)
		keys := []string{}

		for _, env := range os.Environ() {
			parts := strings.SplitN(env, "=", 2)

			if re.FindString(parts[0]) != "" {
				keys = append(keys, parts[0])
			}
		}

		// sort keys for test stability
		sort.Strings(keys)

		for _, key := range keys {
			values = append(values, os.Getenv(key))
		}

		return values
	}

	return []string{os.Getenv(prefix)}
}

type Flags struct {
	VersionFlag bool

	// general
	Paths           string
	Names           string
	RecursiveFlag   bool
	NoRecursiveFlag bool
	Retries         int

	// for destination: env
	EnvOutputFlag    bool
	EnvNoOutputFlag  bool
	EnvPrefix        string
	EnvUseEntirePath bool

	// for destination: file
	FileTargets FileFlags
}

func parseFlags(args []string, flagEnvPrefix string) (*Flags, []string, error) {
	flags := &Flags{}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	fs.BoolVar(&flags.VersionFlag, "version", false, "display version")

	fs.StringVar(&flags.Paths, "paths", "", "comma separated parameter paths")
	fs.StringVar(&flags.Names, "names", "", "comma separated parameter names")
	fs.BoolVar(&flags.RecursiveFlag, "recursive", true, "retrieve values recursively")
	fs.BoolVar(&flags.NoRecursiveFlag, "no-recursive", false, "retrieve values just under -paths only")
	fs.IntVar(&flags.Retries, "retries", 0, "number of times of retry")

	fs.BoolVar(&flags.EnvOutputFlag, "env", true, "export values as environment variables")
	fs.BoolVar(&flags.EnvNoOutputFlag, "no-env", false, "disable export to environment variables")
	fs.StringVar(&flags.EnvPrefix, "env-prefix", "", "prefix for environment variables")
	fs.BoolVar(&flags.EnvUseEntirePath, "env-entire-path", false, "use entire parameter path for name of environment variables\ndisabled: /path/to/value -> VALUE\nenabled: /path/to/value -> PATH_TO_VALUE")
	fs.StringVar(&flags.EnvPrefix, "prefix", "", "alias for -env-prefix")

	fs.Var(&flags.FileTargets, "file", "write values into file\nformat:  Name=VALUE_NAME,Path=FILE_PATH,Mode=FILE_MODE,Gid=FILE_GROUP_ID,Uid=FILE_USER_ID\nexample: Name=/foo/bar,Path=/etc/bar,Mode=600,Gid=123,Uid=456")

	// Read flag values also from environment variable e.g. {flagEnvPrefix}_PATHS
	// Environment variables will be overwritten by flags.
	// Multiple values will be merged.
	fs.VisitAll(func(f *flag.Flag) {
		multiple := false

		if f.Name == "file" {
			multiple = true
		}

		envName := strings.ToUpper(f.Name)
		envName = strings.ReplaceAll(envName, "-", "_")
		envName = flagEnvPrefix + envName

		for _, value := range flagViaEnv(envName, multiple) {
			f.Value.Set(value)
		}
	})

	if err := fs.Parse(args); err != nil {
		return nil, nil, err
	}

	return flags, fs.Args(), nil
}

// Run runs ssmwrap as a CLI, returns exit code.
func Run(version string, flagEnvPrefix string) int {
	flags, restArgs, err := parseFlags(os.Args[1:], flagEnvPrefix)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return 2
	}

	if flags.VersionFlag {
		fmt.Printf("ssmwrap v%s\n", version)
		return 0
	}

	options := ssmwrap.RunOptions{
		Recursive: !flags.NoRecursiveFlag,
		Retries:   flags.Retries,
		Command:   restArgs,
	}
	if len(options.Command) == 0 {
		fmt.Fprintln(os.Stderr, "command required in arguments")
		return 2
	}

	if flags.Paths != "" {
		options.Paths = strings.Split(flags.Paths, ",")
	}
	if flags.Names != "" {
		options.Names = strings.Split(flags.Names, ",")
	}

	ssm := ssmwrap.DefaultSSMConnector{}
	dests := []ssmwrap.Destination{}

	if !flags.EnvNoOutputFlag {
		dests = append(dests, ssmwrap.DestinationEnv{
			Prefix:        flags.EnvPrefix,
			UseEntirePath: flags.EnvUseEntirePath,
		})
	}

	if 0 < len(flags.FileTargets) {
		dests = append(dests, ssmwrap.DestinationFile{
			Targets: flags.FileTargets,
		})
		for _, t := range flags.FileTargets {
			options.Names = append(options.Names, t.Name)
		}
	}

	ctx := context.TODO()

	if err := ssmwrap.Run(ctx, options, ssm, dests); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return 1
	}

	return 0
}
