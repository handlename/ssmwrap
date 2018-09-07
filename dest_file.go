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
	Uid  int
	Gid  int
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

		uid := target.Uid
		if uid == 0 {
			uid = os.Getuid()
		}

		gid := target.Gid
		if gid == 0 {
			gid = os.Getgid()
		}

		err = os.Chown(target.Path, uid, gid)
		if err != nil {
			return errors.Wrapf(err, "failed to chown file %s", target.Path)
		}
	}

	return nil
}
