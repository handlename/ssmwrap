package app

import (
	"fmt"
	"regexp"
	"strings"
)

type ParameterLevel int

const (
	// ParameterLevelStrict means the path will be searched strictly.
	ParameterLevelStrict ParameterLevel = 0

	// ParameterLevelUnder means the path will be searched just under the path.
	ParameterLevelUnder ParameterLevel = 1

	// ParameterLevelAll means the path will be searched under the path recursively.
	ParameterLevelAll ParameterLevel = 2
)

var validPathRegexp = regexp.MustCompile(`^/[-_/a-zA-Z0-9]+((/\**)?/\*)?$`)

type ParameterRule struct {
	// Path is the target path on SSM Parameter Store.
	Path string

	// Level means how deep the path should be searched.
	Level ParameterLevel
}

// NewParameterRule creates a new ParameterRule.
// The path should be a valid path format.
// If the path ends with `/*`, the level will be `ParameterLevelUnder`.
// If the path ends with `/**/*`, the level will be `ParameterLevelAll`.
// Otherwise, the level will be `ParameterLevelStrict`.
func NewParameterRule(path string) (*ParameterRule, error) {
	if !validPathRegexp.MatchString(path) {
		return nil, fmt.Errorf("invalid `path` format")
	}

	if strings.HasSuffix(path, "/**/*") {
		return &ParameterRule{
			Path:  path[:len(path)-4],
			Level: ParameterLevelAll,
		}, nil
	}

	if strings.HasSuffix(path, "/*") {
		return &ParameterRule{
			Path:  path[:len(path)-1],
			Level: ParameterLevelUnder,
		}, nil
	}

	return &ParameterRule{
		Path:  path,
		Level: ParameterLevelStrict,
	}, nil
}

func (r ParameterRule) String() string {
	s := r.Path

	switch r.Level {
	case ParameterLevelStrict:
		// do nothing
	case ParameterLevelUnder:
		s += "*"
	case ParameterLevelAll:
		s += "**/*"
	}

	return s
}

func (r1 ParameterRule) Equals(r2 ParameterRule) bool {
	return r1.Path == r2.Path && r1.Level == r2.Level
}

func (r1 ParameterRule) IsCovers(r2 ParameterRule) bool {
	if r1.Equals(r2) {
		return true
	}

	switch r1.Level {
	case ParameterLevelStrict:
		return false
	case ParameterLevelUnder:
		if r2.Level == ParameterLevelAll {
			return false
		}

		if strings.HasPrefix(r2.Path, r1.Path) {
			// Is r2 just unedr r1?
			s := strings.Replace(r2.Path, r1.Path, "", 1)
			if strings.Contains(s, "/") {
				return false
			}

			return true
		}
	case ParameterLevelAll:
		if strings.HasPrefix(r2.Path, r1.Path) {
			return true
		}
	}

	return false
}
