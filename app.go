package ssmwrap

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/pkg/errors"
)

type SSMConnector interface {
	FetchParameters(paths []string, retries int) (map[string]string, error)
}

type Destination interface {
	Name() string
	Output(parameters map[string]string) error
}

type RunOptions struct {
	// Paths is target paths on SSM Parameter Store.
	// If there are multiple paths, all of related values will be loaded.
	Paths []string

	// Retry limit to request to SSM.
	Retries int

	// Command and arguments to run.
	Command []string
}

func Run(options RunOptions, ssm SSMConnector, dests []Destination) error {
	parameters, err := ssm.FetchParameters(options.Paths, options.Retries)
	if err != nil {
		return errors.Wrap(err, "failed to fetch parameters from SSM")
	}

	for _, dest := range dests {
		err = dest.Output(parameters)
		if err != nil {
			return errors.Wrapf(err, "error occured on output to %s", dest.Name())
		}
	}

	if err := runCommand(options.Command, os.Environ()); err != nil {
		return errors.Wrapf(err, "failed to run command %+s", options.Command)
	}

	return nil
}

func runCommand(command, envVars []string) error {
	bin, err := exec.LookPath(command[0])
	if err != nil {
		return err
	}

	return syscall.Exec(bin, command, envVars)
}
