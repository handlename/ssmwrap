package cli

import (
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/handlename/ssmwrap/internal/app"
	"github.com/handlename/ssmwrap/internal/cli"
	"github.com/samber/lo"
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

	envRules := lo.Times(3, func(i int) app.Rule {
		return app.Rule{
			ParameterRule: app.ParameterRule{
				Path:  fmt.Sprintf("/path/to/env/param%d/", i),
				Level: app.ParameterLevelUnder,
			},
			DestinationRule: app.DestinationRule{
				Type: app.DestinationTypeEnv,
				To:   "",
				TypeEnvOptions: &app.DestinationTypeEnvOptions{
					Prefix:     "",
					EntirePath: false,
				},
			},
		}
	})

	fileRules := lo.Times(3, func(i int) app.Rule {
		return app.Rule{
			ParameterRule: app.ParameterRule{
				Path:  fmt.Sprintf("/path/to/file/param%d/", i),
				Level: app.ParameterLevelStrict,
			},
			DestinationRule: app.DestinationRule{
				Type: app.DestinationTypeFile,
				To:   fmt.Sprintf("/path/to/file%d", i),
				TypeFileOptions: &app.DestinationTypeFileOptions{
					Mode: 0,
					Uid:  0,
					Gid:  0,
				},
			},
		}
	})

	tests := []struct {
		name     string
		flags    []string
		envs     map[string]string
		expected *Flags
	}{
		{
			name: "valid: flags",
			flags: []string{
				"-retries", "3",
				"-rule", envRules[0].String(),
				"-rule", fileRules[0].String(),
				"-env", envRules[1].String(),
				"-env", envRules[2].String(),
				"-file", fileRules[1].String(),
				"-file", fileRules[2].String(),
			},
			expected: &Flags{
				VersionFlag: false,
				Retries:     3,
				RuleFlags: cli.RuleFlags{
					Rules: []app.Rule{
						envRules[0],
						fileRules[0],
					},
				},
				EnvFlags: cli.EnvFlags{
					RuleFlags: cli.RuleFlags{
						Rules: []app.Rule{
							envRules[1],
							envRules[2],
						},
					},
				},
				FileFlags: cli.FileFlags{
					RuleFlags: cli.RuleFlags{
						Rules: []app.Rule{
							fileRules[1],
							fileRules[2],
						},
					},
				},
			},
		},
		{
			name: "valid: envs",
			envs: map[string]string{
				flagEnvPrefix + "RULE_1": envRules[0].String(),
				flagEnvPrefix + "RULE_2": fileRules[0].String(),
				flagEnvPrefix + "ENV_1":  envRules[1].String(),
				flagEnvPrefix + "ENV_2":  envRules[2].String(),
				flagEnvPrefix + "FILE_1": fileRules[1].String(),
				flagEnvPrefix + "FILE_2": fileRules[2].String(),
			},
			expected: &Flags{
				VersionFlag: false,
				Retries:     0,
				RuleFlags: cli.RuleFlags{
					Rules: []app.Rule{
						envRules[0],
						fileRules[0],
					},
				},
				EnvFlags: cli.EnvFlags{
					RuleFlags: cli.RuleFlags{
						Rules: []app.Rule{
							envRules[1],
							envRules[2],
						},
					},
				},
				FileFlags: cli.FileFlags{
					RuleFlags: cli.RuleFlags{
						Rules: []app.Rule{
							fileRules[1],
							fileRules[2],
						},
					},
				},
			},
		},
		{
			name: "valid: flags & envs",
			flags: []string{
				// will overwrite env
				"-retries", "3",

				// will be merged
				"-rule", fileRules[0].String(),
				"-env", envRules[1].String(),
				"-file", fileRules[1].String(),
			},
			envs: map[string]string{
				// will be overwritten by flag
				flagEnvPrefix + "RETRIES": "5",

				// will be merged
				flagEnvPrefix + "RULE": envRules[0].String(),
				flagEnvPrefix + "ENV":  envRules[2].String(),
				flagEnvPrefix + "FILE": fileRules[2].String(),
			},
			expected: &Flags{
				VersionFlag: false,
				Retries:     3,
				RuleFlags: cli.RuleFlags{
					Rules: []app.Rule{
						envRules[0],
						fileRules[0],
					},
				},
				EnvFlags: cli.EnvFlags{
					RuleFlags: cli.RuleFlags{
						Rules: []app.Rule{
							envRules[2],
							envRules[1],
						},
					},
				},
				FileFlags: cli.FileFlags{
					RuleFlags: cli.RuleFlags{
						Rules: []app.Rule{
							fileRules[2],
							fileRules[1],
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset flags
			flag.CommandLine = flag.NewFlagSet("ssmwrap", flag.ExitOnError)

			if tt.envs != nil {
				for k, v := range tt.envs {
					t.Setenv(k, v)
				}
			}

			parsedFlags, _, err := parseFlags(tt.flags, flagEnvPrefix)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			// test.expected.FixOrder()
			if diff := cmp.Diff(tt.expected, parsedFlags); diff != "" {
				t.Errorf("unexpected result: %s", diff)
			}
		})
	}
}
