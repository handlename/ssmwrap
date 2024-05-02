package cli

import (
	"strings"
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
			value: "path=/path/to/param,type=env,prefix=PREFIX_,entirepath=true",
			want: app.Rule{
				ParameterRule: app.ParameterRule{
					Path:  "/path/to/param",
					Level: app.ParameterLevelStrict,
				},
				DestinationRule: app.DestinationRule{
					Type: app.DestinationTypeEnv,
					To:   "",
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

			if len(f.Rules) != 1 {
				t.Fatalf("unexpected length: %d", len(f.Rules))
			}

			if diff := cmp.Diff(tt.want, f.Rules[0]); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func TestRuleFlagsSetReturnsError(t *testing.T) {
	tests := []struct {
		title string
		value string
		err   string
	}{
		{
			title: "type: required",
			value: "path=/path/to/param",
			err:   "invalid `type`",
		},
		{
			title: "path: invalid suffix",
			value: "path=/path/to/param/**,type=env",
			err:   "invalid `path` format",
		},
		{
			title: "path: invalid prefix",
			value: "path=path/to/param,type=env",
			err:   "invalid `path` format",
		},
		{
			title: "path: end with `/*` is not allowed for `type=file`",
			value: "path=/path/to/param/*,type=file,to=/path/to/file",
			err:   "not allowed for `type=file`",
		},
		{
			title: "path: end with `/**/*` not allowed for `type=file`",
			value: "path=/path/to/param/**/*,type=file,to=/path/to/file",
			err:   "not allowed for `type=file`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			var f RuleFlags
			err := f.Set(tt.value)
			if err == nil {
				t.Fatalf("should be error")
			}

			if !strings.Contains(err.Error(), tt.err) {
				t.Errorf("unexpected error: '%s' is not contains '%s'", err, tt.err)
			}
		})
	}
}

func TestRuleFlagsCheckOptionsCombinations(t *testing.T) {
	tests := []struct {
		title    string
		destType app.DestinationType
		opts     map[string]string
		err      string
	}{
		{
			title:    "prefix: only for `type=env`",
			destType: app.DestinationTypeFile,
			opts: map[string]string{
				"to":     "/path/to/file",
				"prefix": "PREFIX_",
			},
			err: "`prefix` is only allowed for `type=env`",
		},
		{
			title:    "entirepath: only for `type=env`",
			destType: app.DestinationTypeFile,
			opts: map[string]string{
				"to":         "/path/to/file",
				"entirepath": "true",
			},
			err: "`entirepath` is only allowed for `type=env`",
		},
		{
			title:    "entirepath: exclusive with `to`",
			destType: app.DestinationTypeEnv,
			opts: map[string]string{
				"to":         "/path/to/file",
				"entirepath": "true",
			},
			err: "can't use `to` with `entirepath`",
		},
		{
			title:    "mode: only for `type=file`",
			destType: app.DestinationTypeEnv,
			opts: map[string]string{
				"mode": "0644",
			},
			err: "is only allowed for `type=file`",
		},
		{
			title:    "uid: only for `type=file`",
			destType: app.DestinationTypeEnv,
			opts: map[string]string{
				"uid": "1000",
			},
			err: "is only allowed for `type=file`",
		},
		{
			title:    "gid: only for `type=file`",
			destType: app.DestinationTypeEnv,
			opts: map[string]string{
				"gid": "1000",
			},
			err: "is only allowed for `type=file`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			var f RuleFlags
			err := f.checkOptionsCombinations(tt.destType, tt.opts)
			if err == nil {
				t.Fatalf("should be error")
			}

			if !strings.Contains(err.Error(), tt.err) {
				t.Errorf("unexpected error: '%s' is not contains '%s'", err, tt.err)
			}
		})
	}
}
