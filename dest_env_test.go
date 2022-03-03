package ssmwrap

import (
	"os"
	"reflect"
	"testing"
)

func TestDestinationEnvFormatParametersAsEnvVars(t *testing.T) {
	patterns := []struct {
		Title              string
		InputParameters    map[string]string
		InputPrefix        string
		InputUseEntirePath bool
		Expected           []string
	}{
		{
			Title: "to upper",
			InputParameters: map[string]string{
				"/foo": "bar",
			},
			Expected: []string{
				"FOO=bar",
			},
		},
		{
			Title: "deep path",
			InputParameters: map[string]string{
				"/d/e/e/p/path": "deep",
			},
			Expected: []string{
				"PATH=deep",
			},
		},
		{
			Title: "deep path, enable entire path",
			InputParameters: map[string]string{
				"/d/e/e/p/path": "deep",
			},
			InputUseEntirePath: true,
			Expected: []string{
				"D_E_E_P_PATH=deep",
			},
		},
		{
			Title: "add prefix",
			InputParameters: map[string]string{
				"/common/name": "john",
			},
			InputPrefix: "MY_",
			Expected: []string{
				"MY_NAME=john",
			},
		},
		{
			Title: "prefix will be upper, too",
			InputParameters: map[string]string{
				"/common/title": "event",
			},
			InputPrefix: "my_",
			Expected: []string{
				"MY_TITLE=event",
			},
		},
	}

	for _, pattern := range patterns {
		t.Log(pattern.Title)

		dest := DestinationEnv{
			Prefix:        pattern.InputPrefix,
			UseEntirePath: pattern.InputUseEntirePath,
		}

		envVars := dest.formatParametersAsEnvVars(pattern.InputParameters)

		if !reflect.DeepEqual(envVars, pattern.Expected) {
			t.Errorf("unexpected envVars\n\tgot:%+v\n\texpected:%+v", envVars, pattern.Expected)
		}
	}
}

func TestDestinationEnvExport(t *testing.T) {
	patterns := []struct {
		Title    string
		EnvVars  []string
		Expected map[string]string
	}{
		{
			Title:   "simple",
			EnvVars: []string{"FOO=foo", "BAR=bar=baz"},
			Expected: map[string]string{
				"FOO": "foo",
				"BAR": "bar=baz",
			},
		},
	}

	dest := DestinationEnv{}

	for _, pattern := range patterns {
		t.Log(pattern.Title)

		err := dest.export(pattern.EnvVars)
		if err != nil {
			t.Error(err)
		}

		for key, value := range pattern.Expected {
			if env := os.Getenv(key); env != value {
				t.Errorf("unexpected env %s=%s (expected %s)", key, env, value)
			}
		}
	}
}
