package cli

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/handlename/ssmwrap"
)

func TestFileTargetsParseFlag(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		expected *ssmwrap.FileTarget
	}{
		{
			name: "valid",
			in:   "Name=/test/foo,Dest=foo.txt,Mode=0600,Uid=1000,Gid=1000",
			expected: &ssmwrap.FileTarget{
				Name: "/test/foo",
				Dest: mustAbsPath(t, "foo.txt"),
				Mode: 0600,
				Uid:  1000,
				Gid:  1000,
			},
		},
	}

	targets := FileFlags{}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			res, err := targets.parseFlag(test.in)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
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

func TestFileTargetsParseFlagRetunsError(t *testing.T) {
	tests := []struct {
		name string
		in   string
		err  string
	}{
		{
			name: "invalid form",
			in:   "foobar",
			err:  "invalid format",
		},
		// how to cause error by filepath.Abs?
		// {
		// 	name: "invalid Dest",
		// },
		{
			name: "invalid Mode",
			in:   "Name=/test/foo,Dest=foo.txt,Mode=foo,Uid=1000,Gid=1000",
			err:  "invalid Mode",
		},
		{
			name: "invalid Uid",
			in:   "Name=/test/foo,Dest=foo.txt,Mode=0600,Uid=bar,Gid=1000",
			err:  "invalid Uid",
		},
		{
			name: "invalid Gid",
			in:   "Name=/test/foo,Dest=foo.txt,Mode=0600,Uid=1000,Gid=buzz",
			err:  "invalid Gid",
		},
	}

	targets := FileFlags{}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			res, err := targets.parseFlag(test.in)
			if err == nil {
				t.Fatalf("expected error but got nil")
			}

			if !strings.Contains(err.Error(), test.err) {
				t.Errorf("error must be %s but %s", test.err, err)
			}

			if res != nil {
				t.Fatalf("unexpected nil result")
			}
		})
	}
}
