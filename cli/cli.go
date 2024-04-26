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
	"github.com/handlename/ssmwrap/internal/app"
	"github.com/handlename/ssmwrap/internal/cli"
)

type ExitStatus int

const (
	ExitStatusOK    ExitStatus = 0
	ExitStatusError ExitStatus = 1
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
	Retries     int

	RuleFlags cli.RuleFlags
	EnvFlags  cli.EnvFlags
	FileFlags cli.FileFlags
}

func parseFlags(args []string, flagEnvPrefix string) (*Flags, []string, error) {
	flags := &Flags{}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	fs.BoolVar(&flags.VersionFlag, "version", false, "Display version and exit")
	fs.IntVar(&flags.Retries, "retries", 0, "Number of times of retry. Default is 0")
	fs.Var(&flags.RuleFlags, "rule", strings.Join([]string{
		"Set rule for exporting values. multiple flags are allowed.",
		"format: path=...,type={env,file},to=...[,entirepath={true,false}][,prefix=...][,mode=...][,gid=...][,uid=...]",
		"parameters:",
		"        path: [required]",
		"              Path of parameter store.",
		"              If `path` ends with no-slash character, only the value of the path will be exported.",
		"              If `path` ends with `/**/*`, all values under the path will be exported.",
		"              If `path` ends with `/*`, only top level values under the path will be exported.",
		"        type: [required]",
		"              Destination type. `env` or `file`.",
		"          to: [required for `type=file`]",
		"              Destination path.",
		"              If `type=env`, `to` is name of exported environment variable.",
		"              If `type=env`, but `to` is not set, `path` will be used as name of exported environment variable.",
		"              If `type=file`, `to` is path of file to write.",
		"  entirepath: [optional, only for `type=env`]",
		"              Export entire path as environment variable.",
		"              If `entirepath=true`, all values under the path will be exported. (/path/to/param -> PATH_TO_PARAM)",
		"              If `entirepath=false`, only top level values under the path will be exported. (/path/to/param -> PARAM)",
		"      prefix: [optional, only for `type=env`]",
		"              Prefix for exported environment variable.",
		"        mode: [optional, only for `type=file`]",
		"              File mode. Default is 0644.",
		"         gid: [optional, only for `type=file`]",
		"              Group ID of file. Default is current user's Gid.",
		"         uid: [optional, only for `type=file`]",
		"              User ID of file. Default is current user's Uid.",
	}, "\n"))
	fs.Var(&flags.EnvFlags, "env", "Alias of `rule` flag with `type=env`.")
	fs.Var(&flags.FileFlags, "file", "Alias of `rule` flag with `type=file`.")

	// Read flag values also from environment variable e.g. {flagEnvPrefix}_PATHS
	// Environment variables will be overwritten by flags.
	// Multiple values will be merged.
	fs.VisitAll(func(f *flag.Flag) {
		multiple := false

		switch f.Name {
		case "rule", "env", "file":
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
func Run(version string, flagEnvPrefix string) ExitStatus {
	ssmwrap.InitLogger()

	flags, restArgs, err := parseFlags(os.Args[1:], flagEnvPrefix)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return ExitStatusError
	}

	if flags.VersionFlag {
		fmt.Printf("ssmwrap v%s\n", version)
		return ExitStatusOK
	}

	options := ssmwrap.RunOptions{
		Recursive: flags.Recursive,
		Retries:   flags.Retries,
		Command:   restArgs,
	}
	if len(options.Command) == 0 {
		fmt.Fprintln(os.Stderr, "command required in arguments")
		return ExitStatusError
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

		return ExitStatusError
	}

	return ExitStatusOK
}
