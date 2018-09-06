package ssmwrap

import (
	"os"
	"reflect"
	"testing"
)

func TestFormatParametersAsEnvVars(t *testing.T) {
	patterns := []struct {
		Title           string
		InputParameters map[string]string
		InputPrefix     string
		Expected        []string
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

		envVars := formatParametersAsEnvVars(pattern.InputParameters, pattern.InputPrefix)

		if !reflect.DeepEqual(envVars, pattern.Expected) {
			t.Errorf("unexpected envVars: %+v", envVars)
		}
	}
}

func TestExport(t *testing.T) {
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

	for _, pattern := range patterns {
		t.Log(pattern.Title)

		err := export(pattern.EnvVars)
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
