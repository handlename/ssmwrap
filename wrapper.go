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

func Run(paths []string, prefix string, command []string) error {
	parameters, err := fetchParameters(paths)
	if err != nil {
		return errors.Wrap(err, "failed to fetch parameters from SSM")
	}

	envVars := prepareEnvVars(parameters, prefix)

	// ssm parameters takes precedence over the current environment variables.
	// In otehr words, ssm parameters overwrite the current environment variables.
	envVars = append(os.Environ(), envVars...)

	if err := runCommand(command, envVars); err != nil {
		return errors.Wrapf(err, "failed to run command %+s", command)
	}

	return nil
}

func fetchParameters(paths []string) (map[string]string, error) {
	params := map[string]string{}

	sess, err := session.NewSession()
	if err != nil {
		return params, errors.Wrap(err, "failed to start session")
	}

	client := ssm.New(sess)

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
