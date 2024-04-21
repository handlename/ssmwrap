package cli

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/handlename/ssmwrap/internal/app"
)

func TestRuleFlagsSetSuccess(t *testing.T) {
	tests := []struct {
		title string
		value string
		want  app.Rule
	}{
		{
			title: "type env (strict)",
			value: "path=/path/to/param,type=env",
			want: app.Rule{
				ParameterRule: app.ParameterRule{
					Path:  "/path/to/param",
					Level: app.ParameterLevelStrict,
				},
				DestinationRule: app.DestinationRule{
					Type: app.DestinationTypeEnv,
					To:   "",
					TypeEnvOptions: &app.DestinationTypeEnvOptions{
						Prefix:     "",
						EntirePath: false,
					},
				},
			},
		},
		{
			title: "type env (under)",
			value: "path=/path/under/*,type=env",
			want: app.Rule{
				ParameterRule: app.ParameterRule{
					Path:  "/path/under/",
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
			},
		},
		{
			title: "type env (all)",
			value: "path=/path/all/**/*,type=env",
			want: app.Rule{
				ParameterRule: app.ParameterRule{
					Path:  "/path/all/",
					Level: app.ParameterLevelAll,
				},
				DestinationRule: app.DestinationRule{
					Type: app.DestinationTypeEnv,
					To:   "",
					TypeEnvOptions: &app.DestinationTypeEnvOptions{
						Prefix:     "",
						EntirePath: false,
					},
				},
			},
		},
		{
			title: "type env with options",
			value: "path=/path/to/param,type=env,to=MY_PARAM,prefix=PREFIX_,entirepath=true",
			want: app.Rule{
				ParameterRule: app.ParameterRule{
					Path:  "/path/to/param",
					Level: app.ParameterLevelStrict,
				},
				DestinationRule: app.DestinationRule{
					Type: app.DestinationTypeEnv,
					To:   "MY_PARAM",
					TypeEnvOptions: &app.DestinationTypeEnvOptions{
						Prefix:     "PREFIX_",
						EntirePath: true,
					},
				},
			},
		},
		{
			title: "type file",
			value: "path=/path/to/param,type=file,to=/path/to/file",
			want: app.Rule{
				ParameterRule: app.ParameterRule{
					Path:  "/path/to/param",
					Level: app.ParameterLevelStrict,
				},
				DestinationRule: app.DestinationRule{
					Type: app.DestinationTypeFile,
					To:   "/path/to/file",
					TypeFileOptions: &app.DestinationTypeFileOptions{
						Mode: 0,
						Uid:  0,
						Gid:  0,
					},
				},
			},
		},
		{
			title: "type file with options",
			value: "path=/path/to/param,type=file,to=/path/to/file,mode=0644,uid=1000,gid=1000",
			want: app.Rule{
				ParameterRule: app.ParameterRule{
					Path:  "/path/to/param",
					Level: app.ParameterLevelStrict,
				},
				DestinationRule: app.DestinationRule{
					Type: app.DestinationTypeFile,
					To:   "/path/to/file",
					TypeFileOptions: &app.DestinationTypeFileOptions{
						Mode: 0644,
						Uid:  1000,
						Gid:  1000,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			var f RuleFlags
			err := f.Set(tt.value)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if len(f) != 1 {
				t.Fatalf("unexpected length: %d", len(f))
			}

			if diff := cmp.Diff(tt.want, f[0]); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}
