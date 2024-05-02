package app

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileExporterExportSuccess(t *testing.T) {
	makeTempfileName := func(t *testing.T) string {
		tmpdir, err := os.MkdirTemp("", "")
		if err != nil {
			t.Errorf("failed to create temp dir: %s", err)
		}

		return filepath.Join(tmpdir, "out")
	}

	tests := []struct {
		title string
		init  func(path string) *FileExporter
	}{
		{
			title: "default",
			init: func(path string) *FileExporter {
				return NewFileExporter(path)
			},
		},
		{
			title: "with mode",
			init: func(path string) *FileExporter {
				ex := NewFileExporter(path)
				ex.Mode = 0600
				return ex
			},
		},
	}

	for _, tt := range tests {
		tempfile := makeTempfileName(t)
		defer os.Remove(tempfile)

		ex := tt.init(tempfile)
		if ex == nil {
			t.Errorf("failed to generate FileExporter")
		}

		value := tt.title
		if err := ex.Export(value); err != nil {
			t.Errorf("failed to export: %s", err)
		}

		f, err := os.Open(ex.Path)
		if err != nil {
			t.Errorf("failed to open destination file: %s", err)
		}

		body, err := io.ReadAll(f)
		if err != nil {
			t.Errorf("failed to read body from file: %s", err)
		}

		if value := string(body); value != value {
			t.Errorf("unexpected body: %s != %s", body, value)
		}
	}
}

func TestFileExporterExportFailedToWrite(t *testing.T) {
	tempDirPath, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Errorf("failed to create tempdir: %s", err)
	}
	defer os.Remove(tempDirPath)

	ex := NewFileExporter(tempDirPath) // directory!!

	err = ex.Export("foo")
	if err == nil {
		t.Errorf("should be error: %s", err)
	}

	if !strings.HasPrefix(err.Error(), "failed to write to file") {
		t.Errorf("unexpected error: %s", err)
	}
}
