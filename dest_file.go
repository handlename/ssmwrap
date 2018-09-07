package ssmwrap

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type DestinationFile struct {
	Targets []FileTarget
}

type FileTarget struct {
	Name string
	Path string
	Mode os.FileMode
}

func (d DestinationFile) Name() string {
	return "File"
}

func (d DestinationFile) Output(parameters map[string]string) error {
	for _, target := range d.Targets {
		err := ioutil.WriteFile(target.Path, []byte(parameters[target.Name]), target.Mode)
		if err != nil {
			return errors.Wrapf(err, "failed to write to file %s", target.Path)
		}
	}

	return nil
}
