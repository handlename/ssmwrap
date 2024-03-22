package ssmwrap

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func mustAbsPath(t *testing.T, path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		t.Fatalf("failed to get absolute path: %s", err)
	}

	return abs
}

func TestFileTargetsParseFlag(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		expected *FileTarget
		err      string
	}{
		{
			name: "valid",
			in:   "Name=/test/foo,Path=foo.txt,Mode=0600,Uid=1000,Gid=1000",
			expected: &FileTarget{
				Name: "/test/foo",
				Path: mustAbsPath(t, "foo.txt"),
				Mode: 0600,
				Uid:  1000,
				Gid:  1000,
			},
		},
		{
			name: "invalid form",
			in:   "foobar",
			err:  "invalid format",
		},
		// how to cause error by filepath.Abs?
		// {
		// 	name: "invalid Path",
		// },
		{
			name: "invalid Mode",
			in:   "Name=/test/foo,Path=foo.txt,Mode=foo,Uid=1000,Gid=1000",
			err:  "invalid Mode",
		},
		{
			name: "invalid Uid",
			in:   "Name=/test/foo,Path=foo.txt,Mode=0600,Uid=bar,Gid=1000",
			err:  "invalid Uid",
		},
		{
			name: "invalid Gid",
			in:   "Name=/test/foo,Path=foo.txt,Mode=0600,Uid=1000,Gid=buzz",
			err:  "invalid Gid",
		},
	}

	targets := FileTargetFlags{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := targets.parseFlag(test.in)
			if err != nil {
				if test.err != "" {
					if strings.Contains(err.Error(), test.err) {
						// expected error, ok
						return
					} else {
						t.Errorf("error must be %s but %s", test.err, err)
					}
				} else {
					t.Errorf("unexpected error: %s", err)
				}
			}

			if res == nil {
				t.Fatalf("unexpected nil result")
			}

			if diff := cmp.Diff(test.expected, res); diff != "" {
				t.Errorf("unexpected result: %s", diff)
			}
		})
	}
}


	}
}
