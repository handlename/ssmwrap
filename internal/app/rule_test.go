package app

import "testing"

func TestRuleString(t *testing.T) {
	tests := []struct {
		title string
		rule  Rule
		want  string
	}{
		{
			title: "type env",
			rule: Rule{
				ParameterRule: ParameterRule{
					Path:  "/path/to/param/",
					Level: ParameterLevelAll,
				},
				DestinationRule: DestinationRule{
					Type: DestinationTypeEnv,
					To:   "",
					TypeEnvOptions: &DestinationTypeEnvOptions{
						Prefix:     "TEST_",
						EntirePath: true,
					},
				},
			},
			want: "path=/path/to/param/**/*,type=env,prefix=TEST_,entirepath=true",
		},
		{
			title: "type file",
			rule: Rule{
				ParameterRule: ParameterRule{
					Path:  "/path/to/param",
					Level: ParameterLevelStrict,
				},
				DestinationRule: DestinationRule{
					Type: DestinationTypeFile,
					To:   "/path/to/file",
					TypeFileOptions: &DestinationTypeFileOptions{
						Mode: 0644,
						Uid:  1000,
						Gid:  2000,
					},
				},
			},
			want: "path=/path/to/param,type=file,to=/path/to/file,mode=0644,uid=1000,gid=2000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			got := tt.rule.String()
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestRuleBuildEnvName(t *testing.T) {
	tests := []struct {
		title      string
		path       string
		prefix     string
		entirePath bool
		want       string
	}{
		{
			title:      "no options",
			path:       "/path/to/param",
			prefix:     "",
			entirePath: false,
			want:       "PARAM",
		},
		{
			title:      "with prefix",
			path:       "/path/to/param",
			prefix:     "TEST_",
			entirePath: false,
			want:       "TEST_PARAM",
		},
		{
			title:      "entire path",
			path:       "/path/to/param",
			prefix:     "",
			entirePath: true,
			want:       "PATH_TO_PARAM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			rule := Rule{
				DestinationRule: DestinationRule{
					TypeEnvOptions: &DestinationTypeEnvOptions{
						Prefix:     tt.prefix,
						EntirePath: tt.entirePath,
					},
				},
			}

			got := rule.buildEnvName(tt.path)
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}
