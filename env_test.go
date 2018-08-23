package ssmwrap

import (
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
			Title: "preserve lower case",
			InputParameters: map[string]string{
				"/foo": "bar",
			},
			Expected: []string{
				"foo=bar",
			},
		},
		{
			Title: "deep path",
			InputParameters: map[string]string{
				"/d/e/e/p/PATH": "deep",
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
				"MY_name=john",
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
