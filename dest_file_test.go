package ssmwrap

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDestinationFileOutputSuccessAll(t *testing.T) {
	targets := []FileTarget{
		// parameter exists -> file its content is parameter value will be created
		{
			Name: "foo",
			Dest: makeTempfileName(t),
			Mode: 0600,
		},

		// parameter not exists -> empty file will be created
		{
			Name: "bar",
			Dest: makeTempfileName(t),
			Mode: 0600,
		},
	}

	parameters := map[string]string{
		"foo": "hogefuga",
	}

	dest := DestinationFile{
		Targets: targets,
	}

	err := dest.Output(parameters)

	if err != nil {
		t.Errorf("failed to output: %s", err)
	}

	for _, target := range targets {
		t.Logf("'%s' to %s", parameters[target.Name], target.Dest)

		f, err := os.Open(target.Dest)
		if err != nil {
			t.Errorf("failed to open target file: %s", err)
		}

		body, err := io.ReadAll(f)
		if err != nil {
			t.Errorf("failed to read body from file: %s", err)
		}

		if value := string(body); value != parameters[target.Name] {
			t.Errorf("unexpected body: %s", body)
		}
	}
}

func TestDestinationFileOutputFailedToWrite(t *testing.T) {
	tempDirPath, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Errorf("failed to create tempdir: %s", err)
	}
	defer os.Remove(tempDirPath)

	t.Log(tempDirPath)

	dest := DestinationFile{
		Targets: []FileTarget{
			{
				Name: "foo",
				Dest: tempDirPath, // directory!!
				Mode: 0600,
			},
		},
	}

	err = dest.Output(map[string]string{})
	if err == nil {
		t.Errorf("should be error: %s", err)
	}

	if !strings.HasPrefix(err.Error(), "failed to write to file") {
		t.Errorf("unexpected error: %s", err)
	}
}

func makeTempfileName(t *testing.T) string {
	tmpdir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Errorf("failed to create temp dir: %s", err)
	}

	return filepath.Join(tmpdir, "out")
}
