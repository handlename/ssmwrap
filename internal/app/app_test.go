package app

import (
	"os"
	"strings"
)

type EnvCleaner struct {
	orig []string
}

func (e *EnvCleaner) Clean() {
	e.orig = os.Environ()
	os.Clearenv()
}

func (e EnvCleaner) Restore() {
	os.Clearenv()

	for _, env := range e.orig {
		p := strings.SplitN(env, "=", 2)
		os.Setenv(p[0], p[1])
	}
}
