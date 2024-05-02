package app

import "os"

type EnvExporter struct {
	Name string
}

func NewEnvExporter(name string) *EnvExporter {
	return &EnvExporter{
		Name: name,
	}
}

func (e EnvExporter) Address() string {
	return e.Name
}

func (e EnvExporter) Export(value string) error {
	return os.Setenv(e.Name, value)
}
