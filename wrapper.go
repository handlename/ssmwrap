package ssmwrap

import (
	"os/exec"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

type Options struct {
	// Paths is target paths on SSM Parameter Store.
	// If there are multiple paths, all of related values will be loaded.
	Paths []string

	// Retry limit to request to SSM.
	Retries int

	// Output values to environment variables.
	EnvOutput bool

	// Prefix for name of exported environment variables.
	EnvPrefix string

	// Command and arguments to run.
	Command []string
}

func Run(options Options) error {
	parameters, err := fetchParameters(options.Paths, options.Retries)
	if err != nil {
		return errors.Wrap(err, "failed to fetch parameters from SSM")
	}

	envVars := []string{}
	if options.EnvOutput {
		envVars = prepareEnvVars(parameters, options.EnvPrefix)
	}

	if err := runCommand(options.Command, envVars); err != nil {
		return errors.Wrapf(err, "failed to run command %+s", options.Command)
	}

	return nil
}

func fetchParameters(paths []string, retries int) (map[string]string, error) {
	params := map[string]string{}

	sess, err := session.NewSession()
	if err != nil {
		return params, errors.Wrap(err, "failed to start session")
	}

	// config.MaxRetries is defaults to -1
	if retries <= 0 {
		retries = -1
	}

	client := ssm.New(sess, aws.NewConfig().WithMaxRetries(retries))

	for _, path := range paths {
		nextToken := ""

		for {
			input := &ssm.GetParametersByPathInput{
				Path:           &path,
				Recursive:      aws.Bool(true),
				WithDecryption: aws.Bool(true),
			}

			if nextToken != "" {
				input.SetNextToken(nextToken)
			}

			output, err := client.GetParametersByPath(input)
			if err != nil {
				return params, errors.Wrap(err, "failed to GetParametersByPath")
			}

			for _, param := range output.Parameters {
				params[*param.Name] = *param.Value
			}

			if output.NextToken == nil {
				break
			}

			nextToken = *output.NextToken
		}
	}

	return params, nil
}

func runCommand(command, envVars []string) error {
	bin, err := exec.LookPath(command[0])
	if err != nil {
		return err
	}

	return syscall.Exec(bin, command, envVars)
}
