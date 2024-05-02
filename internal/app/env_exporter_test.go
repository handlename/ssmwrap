package app

import (
	"os"
	"strings"
	"testing"
)

func TestEnvExporterExportSuccess(t *testing.T) {
	tests := []struct {
		title string
		init  func() *EnvExporter
	}{
		{
			title: "normal",
			init: func() *EnvExporter {
				return NewEnvExporter("TEST")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			cleaner := EnvCleaner{}
			cleaner.Clean()
			defer cleaner.Restore()

			ex := tt.init()
			if ex == nil {
				t.Errorf("failed to generate EnvExporter")
			}

			value := tt.title
			if err := ex.Export(value); err != nil {
				t.Errorf("failed to export: %s", err)
			}

			if env := os.Getenv(ex.Name); env != value {
				t.Errorf("unexpected env %s=%s (expected %s)", ex.Name, env, value)
			}
		})
	}
}

func TestEnvExporterExportReturnsError(t *testing.T) {
	tests := []struct {
		title string
		init  func() *EnvExporter
		err   string
	}{
		{
			title: "name contains `=`",
			init: func() *EnvExporter {
				return NewEnvExporter("LEFT=RIGHT")
			},
			err: "invalid argument",
		},
		{
			title: "name is empty string",
			init: func() *EnvExporter {
				return NewEnvExporter("")
			},
			err: "invalid argument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			cleaner := EnvCleaner{}
			cleaner.Clean()
			defer cleaner.Restore()

			ex := tt.init()
			if ex == nil {
				t.Errorf("failed to generate EnvExporter")
			}

			value := tt.title
			err := ex.Export(value)
			if err == nil {
				t.Fatal("expected error")
			}

			if !strings.Contains(err.Error(), tt.err) {
				t.Errorf("unexpected error: %s", err)
			}
		})
	}
}
