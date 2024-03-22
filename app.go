package ssmwrap

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type SSMConnector interface {
	fetchParametersByPaths(ctx context.Context, client *ssm.Client, paths []string, recursive bool) (map[string]string, error)
	fetchParametersByNames(ctx context.Context, client *ssm.Client, names []string) (map[string]string, error)
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

func Run(ctx context.Context, options RunOptions, ssm SSMConnector, dests []Destination) error {
	client, err := newSSMClient(ctx, options.Retries)
	if err != nil {
		return err
	}
	parameters := map[string]string{}

	{
		p, err := ssm.fetchParametersByPaths(ctx, client, options.Paths, options.Recursive)
		if err != nil {
			return fmt.Errorf("failed to fetch parameters from SSM: %w", err)
		}
		for key, value := range p {
			parameters[key] = value
		}
	}

	{
		p, err := ssm.fetchParametersByNames(ctx, client, options.Names)
		if err != nil {
			return fmt.Errorf("failed to fetch parameters from SSM: %w", err)
		}
		for key, value := range p {
			parameters[key] = value
		}
	}

	for _, dest := range dests {
		err = dest.Output(parameters)
		if err != nil {
			return fmt.Errorf("error occured on output to %s: %w", dest.Name(), err)
		}
	}

	if err := runCommand(options.Command, os.Environ()); err != nil {
		return fmt.Errorf("failed to run command %+s: %w", options.Command, err)
	}

	return nil
}

func runCommand(command, envVars []string) error {
	if len(command) == 0 {
		return errors.New("command required")
	}
	bin, err := exec.LookPath(command[0])
	if err != nil {
		return err
	}

	return syscall.Exec(bin, command, envVars)
}
