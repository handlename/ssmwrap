package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/handlename/ssmwrap"
)

type FileFlags []ssmwrap.FileTarget

func (ts *FileFlags) String() string {
	s := ""

	for _, t := range *ts {
		s += t.String()
	}

	return s
}

func (ts *FileFlags) Set(value string) error {
	target, err := ts.parseFlag(value)
	if err != nil {
		return err
	}

	if err := target.IsSatisfied(); err != nil {
		return fmt.Errorf("file parameter is not satisfied: %s", err)
	}

	*ts = append(*ts, *target)

	return nil
}

func (ts FileFlags) parseFlag(value string) (*ssmwrap.FileTarget, error) {
	target := &ssmwrap.FileTarget{}
	parts := strings.Split(value, ",")

	for _, part := range parts {
		param := strings.Split(part, "=")
		if len(param) != 2 {
			return nil, fmt.Errorf("invalid format")
		}

		key := param[0]
		value := param[1]

		switch key {
		case "Name":
			target.Name = value
		case "Dest":
			dest, err := ts.parseDest(value)
			if err != nil {
				return nil, fmt.Errorf("invalid Dest: %s", err)
			}
			target.Dest = dest
		case "Mode":
			mode, err := ts.parseMode(value)
			if err != nil {
				return nil, fmt.Errorf("invalid Mode: %s", err)
			}
			target.Mode = mode
		case "Uid":
			uid, err := ts.parseUid(value)
			if err != nil {
				return nil, fmt.Errorf("invalid Uid: %s", err)
			}
			target.Uid = uid
		case "Gid":
			gid, err := ts.parseGid(value)
			if err != nil {
				return nil, fmt.Errorf("invalid Gid: %s", err)
			}
			target.Gid = gid
		default:
			return nil, fmt.Errorf("unknown parameter: %s", key)
		}
	}

	return target, nil
}

func (ts FileFlags) parseDest(value string) (string, error) {
	// expand destination path
	dest, err := filepath.Abs(value)
	if err != nil {
		return "", fmt.Errorf("Invalid Dest")
	}

	return dest, nil
}

func (ts FileFlags) parseMode(value string) (os.FileMode, error) {
	// convert `Mode` to os.FileMode
	mode, err := strconv.ParseUint(value, 8, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid Mode")
	}

	return os.FileMode(mode), nil
}

func (ts FileFlags) parseGid(value string) (int, error) {
	gid, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid Gid")
	}

	return gid, nil
}

func (ts FileFlags) parseUid(value string) (int, error) {
	uid, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid Uid")
	}

	return uid, nil
}
