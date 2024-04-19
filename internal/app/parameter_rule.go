package app

type ParameterLevel int

const (
	// ParameterLevelStrict means the path will be searched strictly.
	ParameterLevelStrict ParameterLevel = 0

	// ParameterLevelUnder means the path will be searched just under the path.
	ParameterLevelUnder ParameterLevel = 1

	// ParameterLevelAll means the path will be searched under the path recursively.
	ParameterLevelAll ParameterLevel = 2
)

type ParameterRule struct {
	// Path is the target path on SSM Parameter Store.
	Path string

	// Level means how deep the path should be searched.
	Level ParameterLevel
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
