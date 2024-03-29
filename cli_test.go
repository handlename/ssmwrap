package ssmwrap

import (
	"flag"
	"path/filepath"
	"strings"
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

func TestFileTargetsParseFlag(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		expected *FileTarget
	}{
		{
			name: "valid",
			in:   "Name=/test/foo,Path=foo.txt,Mode=0600,Uid=1000,Gid=1000",
			expected: &FileTarget{
				Name: "/test/foo",
				Path: mustAbsPath(t, "foo.txt"),
				Mode: 0600,
				Uid:  1000,
				Gid:  1000,
			},
		},
	}

	targets := FileTargetFlags{}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			res, err := targets.parseFlag(test.in)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if res == nil {
				t.Fatalf("unexpected nil result")
			}

			if diff := cmp.Diff(test.expected, res); diff != "" {
				t.Errorf("unexpected result: %s", diff)
			}
		})
	}
}

func TestFileTargetsParseFlagRetunsError(t *testing.T) {
	tests := []struct {
		name string
		in   string
		err  string
	}{
		{
			name: "invalid form",
			in:   "foobar",
			err:  "invalid format",
		},
		// how to cause error by filepath.Abs?
		// {
		// 	name: "invalid Path",
		// },
		{
			name: "invalid Mode",
			in:   "Name=/test/foo,Path=foo.txt,Mode=foo,Uid=1000,Gid=1000",
			err:  "invalid Mode",
		},
		{
			name: "invalid Uid",
			in:   "Name=/test/foo,Path=foo.txt,Mode=0600,Uid=bar,Gid=1000",
			err:  "invalid Uid",
		},
		{
			name: "invalid Gid",
			in:   "Name=/test/foo,Path=foo.txt,Mode=0600,Uid=1000,Gid=buzz",
			err:  "invalid Gid",
		},
	}

	targets := FileTargetFlags{}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			res, err := targets.parseFlag(test.in)
			if err == nil {
				t.Fatalf("expected error but got nil")
			}

			if !strings.Contains(err.Error(), test.err) {
				t.Errorf("error must be %s but %s", test.err, err)
			}

			if res != nil {
				t.Fatalf("unexpected nil result")
			}
		})
	}
}

func TestParseFlag(t *testing.T) {
	flagEnvPrefix := "SSMWRAP_TEST_"

	tests := []struct {
		name     string
		flags    []string
		envs     map[string]string
		expected *CLIFlags
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
			expected: &CLIFlags{
				VersionFlag: false,

				Paths:           "/foo/",
				Names:           "foo,bar",
				RecursiveFlag:   true,
				NoRecursiveFlag: false,
				Retries:         3,

				EnvOutputFlag:    true,
				EnvNoOutputFlag:  false,
				EnvPrefix:        "TEST_",
				EnvUseEntirePath: false,

				FileTargets: FileTargetFlags{
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
			expected: &CLIFlags{
				VersionFlag: false,

				Paths:           "/foo/",
				Names:           "foo,bar",
				RecursiveFlag:   true,
				NoRecursiveFlag: false,
				Retries:         3,

				EnvOutputFlag:    true,
				EnvNoOutputFlag:  false,
				EnvPrefix:        "TEST_",
				EnvUseEntirePath: false,

				FileTargets: FileTargetFlags{
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
			expected: &CLIFlags{
				VersionFlag: false,

				Paths:           "/foo/",
				Names:           "foo,bar",
				RecursiveFlag:   false,
				NoRecursiveFlag: false,
				Retries:         3,

				EnvOutputFlag:    false,
				EnvNoOutputFlag:  false,
				EnvPrefix:        "",
				EnvUseEntirePath: false,

				FileTargets: FileTargetFlags{
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

			parsedFlags, _, err := parseCLIFlags(test.flags, flagEnvPrefix)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(test.expected, parsedFlags); diff != "" {
				t.Errorf("unexpected result: %s", diff)
			}
		})
	}
}
