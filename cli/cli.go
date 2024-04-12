package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strings"
	"syscall"

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
	Paths     string
	Names     string
	Recursive bool
	Retries   int

	// for destination: env
	EnvOutput        bool
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
	fs.BoolVar(&flags.Recursive, "recursive", false, "retrieve values recursively")
	fs.IntVar(&flags.Retries, "retries", 0, "number of times of retry")

	fs.BoolVar(&flags.EnvOutput, "env", false, "export values as environment variables")
	fs.StringVar(&flags.EnvPrefix, "env-prefix", "", "prefix for environment variables")
	fs.BoolVar(&flags.EnvUseEntirePath, "env-entire-path", false, "use entire parameter path for name of environment variables\ndisabled: /path/to/value -> VALUE\nenabled: /path/to/value -> PATH_TO_VALUE")
	fs.StringVar(&flags.EnvPrefix, "prefix", "", "alias for -env-prefix")

	// for destination: file
	fs.Var(&flags.FileTargets, "file", strings.Join([]string{
		"write values into file. multiple flags are allowed.",
		"format: Name=VALUE_NAME,Dest=FILE_PATH,Mode=FILE_MODE[,Gid=FILE_GROUP_ID][,Uid=FILE_USER_ID]",
		"        write value of VALUE_NAME into FILE_PATH with FILE_MODE.",
		"        if no Gid and Uid, current user's Gid and Uid will be used.",
	}, "\n"))

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
		Recursive: flags.Recursive,
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

	if flags.EnvOutput {
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	if err := ssmwrap.Run(ctx, options, ssm, dests); err != nil {
		if errors.Is(err, context.Canceled) {
			fmt.Fprintf(os.Stderr, "Interrupted\n")
		} else {
			fmt.Fprintf(os.Stderr, "Erorr occurred: %s", err)
		}

		return 1
	}

	return 0
}
