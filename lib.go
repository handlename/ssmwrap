package ssmwrap

import (
	"fmt"
)

type ExportOptions struct {
	Paths         []string
	Names         []string
	Prefix        string
	UseEntirePath bool
	Recursive     bool
	Retries       int
}

// Export fetches paramters from SSM and export those to environment variables.
// This is for use ssmwrap as a library.
func Export(options ExportOptions) error {
	ssm := DefaultSSMConnector{}
	dest := DestinationEnv{
		Prefix:        options.Prefix,
		UseEntirePath: options.UseEntirePath,
	}

	parameters := map[string]string{}
	client, err := newSSMClient(options.Retries)
	if err != nil {
		return err
	}

	{
		p, err := ssm.fetchParametersByPaths(client, options.Paths, options.Recursive)
		if err != nil {
			return fmt.Errorf("failed to fetch parameters from SSM: %w", err)
		}
		for key, value := range p {
			parameters[key] = value
		}
	}

	{
		p, err := ssm.fetchParametersByNames(client, options.Names)
		if err != nil {
			return fmt.Errorf("failed to fetch parameters from SSM: %w", err)
		}
		for key, value := range p {
			parameters[key] = value
		}
	}

	envVars := dest.formatParametersAsEnvVars(parameters)
	return dest.export(envVars)
}
