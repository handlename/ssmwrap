package ssmwrap

import (
	"fmt"
	"os"
	"strings"
)

// prepareEnvironmentVariables transform SSM parameters to environment variables like `FOO=bar`
// Tha last parts of parameter name separated by `/` will be used.
// `prefix` will append to head of name of environment variables.
func prepareEnvVars(parameters map[string]string, prefix string) []string {
	envVars := formatParametersAsEnvVars(parameters, prefix)

	// ssm parameters takes precedence over the current environment variables.
	// In otehr words, ssm parameters overwrite the current environment variables.
	envVars = append(os.Environ(), envVars...)

	return envVars
}

func formatParametersAsEnvVars(parameters map[string]string, prefix string) []string {
	envVars := []string{}

	for name, value := range parameters {
		parts := strings.Split(name, "/")
		key := strings.ToUpper(prefix + parts[len(parts)-1])
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	return envVars
}
