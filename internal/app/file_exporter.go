package app

import (
	"fmt"
	"io/fs"
	"os"
)

type FileExporter struct {
	Path string
	Mode fs.FileMode
	Uid  int
	Gid  int
}

func NewFileExporter(path string) *FileExporter {
	return &FileExporter{
		Path: path,
		Mode: 0644,
		Uid:  os.Getuid(),
		Gid:  os.Getgid(),
	}
}

func (e FileExporter) Export(v string) error {
	err := os.WriteFile(e.Path, []byte(v), e.Mode)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", e.Path, err)
	}

	uid := e.Uid
	if uid == 0 {
		uid = os.Getuid()
	}

	gid := e.Gid
	if gid == 0 {
		gid = os.Getgid()
	}

	err = os.Chown(e.Path, uid, gid)
	if err != nil {
		return fmt.Errorf("failed to chown file %s: %w", e.Path, err)
	}

	return nil
}
