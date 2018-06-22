package ssmwrap

import (
	"reflect"
	"testing"
)

func TestPrepareEnvVars(t *testing.T) {
	patterns := []struct {
		Title    string
		Input    map[string]string
		Expected []string
	}{
		{
			Title: "to upper",
			Input: map[string]string{
				"/foo": "bar",
			},
			Expected: []string{
				"FOO=bar",
			},
		},
		{
			Title: "deep path",
			Input: map[string]string{
				"/d/e/e/p/path": "deep",
			},
			Expected: []string{
				"PATH=deep",
			},
		},
	}

	for _, pattern := range patterns {
		t.Log(pattern.Title)

		envVars := prepareEnvVars(pattern.Input)

		if !reflect.DeepEqual(envVars, pattern.Expected) {
			t.Errorf("unexpected envVars: %+v", envVars)
		}
	}
}
