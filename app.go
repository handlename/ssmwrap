package ssmwrap

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

type SSMConnector interface {
	fetchParametersByPaths(client *ssm.SSM, paths []string, recursive bool) (map[string]string, error)
	fetchParametersByNames(client *ssm.SSM, names []string) (map[string]string, error)
}

type Destination interface {
	Name() string
	Output(parameters map[string]string) error
}

type RunOptions struct {
	// Paths is target paths on SSM Parameter Store.
	// If there are multiple paths, all of related values will be loaded.
	Paths []string

	// Names are names on SSM Parameter Store.
	Names []string

	// Recursive tell ssmwrap to retrieve values from SSM recursively.
	Recursive bool

	// Retry limit to request to SSM.
	Retries int

	// Command and arguments to run.
	Command []string
}

func Run(options RunOptions, ssm SSMConnector, dests []Destination) error {
	client, err := newSSMClient(options.Retries)
	if err != nil {
		return err
	}
	parameters := map[string]string{}

	{
		p, err := ssm.fetchParametersByPaths(client, options.Paths, options.Recursive)
		if err != nil {
			return errors.Wrap(err, "failed to fetch parameters from SSM")
		}
		for key, value := range p {
			parameters[key] = value
		}
	}

	{
		p, err := ssm.fetchParametersByNames(client, options.Names)
		if err != nil {
			return errors.Wrap(err, "failed to fetch parameters from SSM")
		}
		for key, value := range p {
			parameters[key] = value
		}
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
