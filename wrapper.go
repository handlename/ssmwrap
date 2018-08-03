package ssmwrap

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

	// ssm parameters takes precedence over the current environment variables.
	// In otehr words, ssm parameters overwrite the current environment variables.
	envVars = append(os.Environ(), envVars...)

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

// prepareEnvironmentVariables transform SSM parameters to environment variables like `FOO=bar`
// Tha last parts of parameter name separated by `/` will be used.
// `prefix` will append to head of name of environment variables.
func prepareEnvVars(parameters map[string]string, prefix string) []string {
	envVars := []string{}

	for name, value := range parameters {
		parts := strings.Split(name, "/")
		key := strings.ToUpper(prefix + parts[len(parts)-1])
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	return envVars
}

func runCommand(command, envVars []string) error {
	bin, err := exec.LookPath(command[0])
	if err != nil {
		return err
	}

	return syscall.Exec(bin, command, envVars)
}
