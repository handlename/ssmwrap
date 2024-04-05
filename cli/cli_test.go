package cli

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func mustAbsPath(t *testing.T, path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		t.Fatalf("failed to get absolute path: %s", err)
	}

	return abs
}

func TestParseFlag(t *testing.T) {
	flagEnvPrefix := "SSMWRAP_TEST_"

	tests := []struct {
		name     string
		flags    []string
		envs     map[string]string
		expected *Flags
	}{
		{
			name: "valid: flags",
			flags: []string{
				"-env",
				"-paths", "/foo/",
				"-names", "foo,bar",
				"-recursive",
				"-retries", "3",
				"-env-prefix", "TEST_",
				"-file", "Name=/foo/,Path=foo.txt,Mode=0600,Uid=1000,Gid=1000",
			},
			expected: &Flags{
				VersionFlag: false,

				Paths:     "/foo/",
				Names:     "foo,bar",
				Recursive: true,
				Retries:   3,

				EnvOutput:        true,
				EnvPrefix:        "TEST_",
				EnvUseEntirePath: false,

				FileTargets: FileFlags{
					{
						Name: "/foo/",
						Path: mustAbsPath(t, "foo.txt"),
						Mode: 0600,
						Uid:  1000,
						Gid:  1000,
					},
				},
			},
		},
		{
			name: "valid: envs",
			envs: map[string]string{
				flagEnvPrefix + "ENV":       "1",
				flagEnvPrefix + "PATHS":     "/foo/",
				flagEnvPrefix + "NAMES":     "foo,bar",
				flagEnvPrefix + "RECURSIVE": "1",
				flagEnvPrefix + "RETRIES":   "3",
				flagEnvPrefix + "PREFIX":    "TEST_",
				flagEnvPrefix + "FILE_1":    "Name=/foo/,Path=foo.txt,Mode=0600,Uid=1000,Gid=1000",
				flagEnvPrefix + "FILE_2":    "Name=/bar/,Path=bar.txt,Mode=0600,Uid=2000,Gid=2000",
			},
			expected: &Flags{
				VersionFlag: false,

				Paths:     "/foo/",
				Names:     "foo,bar",
				Recursive: true,
				Retries:   3,

				EnvOutput:        true,
				EnvPrefix:        "TEST_",
				EnvUseEntirePath: false,

				FileTargets: FileFlags{
					{
						Name: "/foo/",
						Path: mustAbsPath(t, "foo.txt"),
						Mode: 0600,
						Uid:  1000,
						Gid:  1000,
					},
					{
						Name: "/bar/",
						Path: mustAbsPath(t, "bar.txt"),
						Mode: 0600,
						Uid:  2000,
						Gid:  2000,
					},
				},
			},
		},
		{
			name: "valid: flags & envs",
			envs: map[string]string{
				// only by envs
				flagEnvPrefix + "NAMES": "foo,bar",

				// will be overwriten by flags
				flagEnvPrefix + "PATHS": "/bar/",

				// multiple envs will be merged
				flagEnvPrefix + "FILE": "Name=/foo/,Path=foo.txt,Mode=0600,Uid=1000,Gid=1000",
			},
			flags: []string{
				// only by flags
				"-retries", "3",

				// flags overwrites env
				"-paths", "/foo/",

				// multiple flags will be merged
				"-file", "Name=/bar/,Path=bar.txt,Mode=0600,Uid=2000,Gid=2000",
			},
			expected: &Flags{
				VersionFlag: false,

				Paths:     "/foo/",
				Names:     "foo,bar",
				Recursive: false,
				Retries:   3,

				EnvOutput:        false,
				EnvPrefix:        "",
				EnvUseEntirePath: false,

				FileTargets: FileFlags{
					{
						Name: "/foo/",
						Path: mustAbsPath(t, "foo.txt"),
						Mode: 0600,
						Uid:  1000,
						Gid:  1000,
					},
					{
						Name: "/bar/",
						Path: mustAbsPath(t, "bar.txt"),
						Mode: 0600,
						Uid:  2000,
						Gid:  2000,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			// reset flags
			flag.CommandLine = flag.NewFlagSet("ssmwrap", flag.ExitOnError)

			if test.envs != nil {
				for k, v := range test.envs {
					t.Setenv(k, v)
				}
			}

			parsedFlags, _, err := parseFlags(test.flags, flagEnvPrefix)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(test.expected, parsedFlags); diff != "" {
				t.Errorf("unexpected result: %s", diff)
			}
		})
	}
}
