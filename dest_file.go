package ssmwrap

import (
	"fmt"
	"io/ioutil"
	"os"
)

// DestinationFile is an implementation of Destination interface.
type DestinationFile struct {
	Targets []FileTarget
}

// TODO: func NewFileTarget(...)

func (d DestinationFile) Name() string {
	return "File"
}

func (d DestinationFile) Output(parameters map[string]string) error {
	for _, target := range d.Targets {
		err := ioutil.WriteFile(target.Path, []byte(parameters[target.Name]), target.Mode)
		if err != nil {
			return fmt.Errorf("failed to write to file %s: %w", target.Path, err)
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
			return fmt.Errorf("failed to chown file %s: %w", target.Path, err)
		}
	}

	return nil
}

type FileTarget struct {
	Name string
	Path string
	Mode os.FileMode
	Uid  int
	Gid  int
}

func (t FileTarget) IsSatisfied() error {
	if t.Name == "" {
		return fmt.Errorf("Name is required")
	}

	if t.Path == "" {
		return fmt.Errorf("Path is required")
	}

	return nil
}

func (t FileTarget) String() string {
	return fmt.Sprintf("Name=%s,Path=%s,Mode=%d,Uid=%d,Gid=%d", t.Name, t.Path, t.Mode, t.Uid, t.Gid)
}
