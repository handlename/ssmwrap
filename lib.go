package ssmwrap

import (
	"github.com/pkg/errors"
)

type ExportOptions struct {
	Paths     []string
	Names     []string
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

	parameters := map[string]string{}
	client, err := newSSMClient(options.Retries)
	if err != nil {
		return err
	}

	{
		p, err := ssm.fetchParameters(client, options.Paths, options.Recursive)
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

	envVars := dest.formatParametersAsEnvVars(parameters)
	return dest.export(envVars)
}
