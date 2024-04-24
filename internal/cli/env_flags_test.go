package cli

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/handlename/ssmwrap/internal/app"
)

func TestEnvFlagsSuccess(t *testing.T) {
	var f EnvFlags
	err := f.Set("path=/path/to/param,entirepath=true,prefix=PREFIX_")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(f.Rules) != 1 {
		t.Fatalf("unexpected length: %d", len(f.Rules))
	}

	want := app.Rule{
		ParameterRule: app.ParameterRule{
			Path:  "/path/to/param",
			Level: 0,
		},
		DestinationRule: app.DestinationRule{
			Type: app.DestinationTypeEnv,
			To:   "",
			TypeEnvOptions: &app.DestinationTypeEnvOptions{
				Prefix:     "PREFIX_",
				EntirePath: true,
			},
		},
	}

	if diff := cmp.Diff(want, f.Rules[0]); diff != "" {
		t.Fatalf("unexpected diff: %s", diff)
	}
}
