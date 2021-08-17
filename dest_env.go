package ssmwrap

import (
	"fmt"
	"os"
	"strings"
)

// DestinationEnv is an implementation of Destination interface.
type DestinationEnv struct {
	Prefix string
}

func (d DestinationEnv) Name() string {
	return "Env"
}

func (d DestinationEnv) Output(parameters map[string]string) error {
	envVars := d.prepareEnvVars(parameters)
	return d.export(envVars)
}

// prepareEnvironmentVariables transform SSM parameters to environment variables like `FOO=bar`
// Tha last parts of parameter name separated by `/` will be used.
// `prefix` will append to head of name of environment variables.
func (d DestinationEnv) prepareEnvVars(parameters map[string]string) []string {
	envVars := d.formatParametersAsEnvVars(parameters)

	// ssm parameters takes precedence over the current environment variables.
	// In otehr words, ssm parameters overwrite the current environment variables.
	envVars = append(os.Environ(), envVars...)

	return envVars
}

func (d DestinationEnv) formatParametersAsEnvVars(parameters map[string]string) []string {
	envVars := []string{}

	for name, value := range parameters {
		parts := strings.Split(name, "/")
		key := strings.ToUpper(d.Prefix + strings.Join(parts[1:], "_"))
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	return envVars
}

func (d DestinationEnv) export(envVars []string) error {
	for _, v := range envVars {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("= is not contained in envvars")
		}
		if err := os.Setenv(parts[0], parts[1]); err != nil {
			return fmt.Errorf("setenv failed: %w", err)
		}
	}
	return nil
}
