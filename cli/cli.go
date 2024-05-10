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

	"github.com/handlename/ssmwrap/v2"
	"github.com/handlename/ssmwrap/v2/internal/app"
	"github.com/handlename/ssmwrap/v2/internal/cli"
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
		"format: path=...,type={env,file}[,to=...][,entirepath={true,false}][,prefix=...][,mode=...][,gid=...][,uid=...]",
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
		"              Export entire path as environment variables name.",
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

	command := restArgs
	if (0 < len(command)) && (command[0] == "--") {
		command = command[1:]
	}
	if len(command) == 0 {
		fmt.Fprintln(os.Stderr, "command required in arguments")
		return ExitStatusError
	}

	rules := []app.Rule{}
	rules = append(rules, flags.RuleFlags.Rules...)
	rules = append(rules, flags.EnvFlags.Rules...)
	rules = append(rules, flags.FileFlags.Rules...)
	if len(rules) == 0 {
		fmt.Fprintf(os.Stderr, "At least one rule required\n")
		return ExitStatusError
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	sw := app.NewSSMWrap()
	if flags.Retries != 0 {
		sw.Retries = flags.Retries
	}

	if err := sw.Run(ctx, rules, command); err != nil {
		if errors.Is(err, context.Canceled) {
			fmt.Fprintf(os.Stderr, "Interrupted\n")
		} else {
			fmt.Fprintf(os.Stderr, "Error occurred: %s\n", err)
		}

		return ExitStatusError
	}

	return ExitStatusOK
}
