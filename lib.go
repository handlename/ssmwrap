package ssmwrap

import (
	"github.com/pkg/errors"
)

type ExportOptions struct {
	Paths     []string
	Prefix    string
	Recursive bool
	Retries   int
}

// Export fetches paramters from SSM and export those to environment variables.
// This is for use ssmwrap as a library.
func Export(options ExportOptions) error {
	ssm := DefaultSSMConnector{}
	dest := DestinationEnv{
		Prefix: options.Prefix,
	}

	parameters, err := ssm.FetchParameters(options.Paths, options.Recursive, options.Retries)
	if err != nil {
		return errors.Wrap(err, "failed to fetch parameters from SSM")
	}

	envVars := dest.formatParametersAsEnvVars(parameters)
	return dest.export(envVars)
}
