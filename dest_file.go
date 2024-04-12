package ssmwrap

import (
	"fmt"
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
		err := os.WriteFile(target.Dest, []byte(parameters[target.Name]), target.Mode)
		if err != nil {
			return fmt.Errorf("failed to write to file %s: %w", target.Dest, err)
		}

		uid := target.Uid
		if uid == 0 {
			uid = os.Getuid()
		}

		gid := target.Gid
		if gid == 0 {
			gid = os.Getgid()
		}

		err = os.Chown(target.Dest, uid, gid)
		if err != nil {
			return fmt.Errorf("failed to chown file %s: %w", target.Dest, err)
		}
	}

	return nil
}

type FileTarget struct {
	Name string
	Dest string
	Mode os.FileMode
	Uid  int
	Gid  int
}

func (t FileTarget) IsSatisfied() error {
	if t.Name == "" {
		return fmt.Errorf("Name is required")
	}

	if t.Dest == "" {
		return fmt.Errorf("Path is required")
	}

	return nil
}

func (t FileTarget) String() string {
	return fmt.Sprintf("Name=%s,Path=%s,Mode=%d,Uid=%d,Gid=%d", t.Name, t.Dest, t.Mode, t.Uid, t.Gid)
}
