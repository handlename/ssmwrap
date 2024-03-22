package ssmwrap

import (
	"strings"
	"testing"
)

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
				Path: "foo.txt",
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

	targets := FileTargets{}

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

			if res.Name != test.expected.Name {
				t.Errorf("unexpected Name: %s", res.Name)
			}

			absPath, err := targets.parsePath(test.expected.Path)
			if err != nil {
				t.Errorf("invalid expected path %s: %s", test.expected.Path, err)
			}
			if res.Path != absPath {
				t.Errorf("unexpected Path: %s", res.Path)
			}

			if res.Mode != test.expected.Mode {
				t.Errorf("unexpected Mode: %d", res.Mode)
			}

			if res.Uid != test.expected.Uid {
				t.Errorf("unexpected Uid: %d", res.Uid)
			}

			if res.Gid != test.expected.Gid {
				t.Errorf("unexpected Gid: %d", res.Gid)
			}
		})
	}
}
