package cli

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/handlename/ssmwrap/internal/app"
)

func TestFileFlagsSuccess(t *testing.T) {
	var f FileFlags
	err := f.Set("path=/path/to/param,to=/path/to/file,mode=0644,uid=1000,gid=1000")
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
			Type: app.DestinationTypeFile,
			To:   "/path/to/file",
			TypeFileOptions: &app.DestinationTypeFileOptions{
				Mode: 0644,
				Uid:  1000,
				Gid:  1000,
			},
		},
	}

	if diff := cmp.Diff(want, f.Rules[0]); diff != "" {
		t.Fatalf("unexpected diff: %s", diff)
	}
}
