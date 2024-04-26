package ssmwrap

import (
	"context"
	"fmt"

	"github.com/handlename/ssmwrap/internal/app"
)

type ExportOptions struct {
	Paths         []string
	Prefix        string
	UseEntirePath bool
	Retries       int
}

// Export fetches parameters from SSM and export those to environment variables.
// This is for use ssmwrap as a library.
func Export(ctx context.Context, options ExportOptions) error {
	rules := make([]app.Rule, 0, len(options.Paths))

	for _, path := range options.Paths {
		pr, err := app.NewParameterRule(path)
		if err != nil {
			return fmt.Errorf("failed to create ParameterRule: %w", err)
		}

		rules = append(rules, app.Rule{
			ParameterRule: *pr,
			DestinationRule: app.DestinationRule{
				Type: app.DestinationTypeEnv,
				TypeEnvOptions: &app.DestinationTypeEnvOptions{
					Prefix:     options.Prefix,
					EntirePath: options.UseEntirePath,
				},
			},
		})
	}

	sw := app.NewSSMWrap()
	if options.Retries != 0 {
		sw.Retries = options.Retries
	}

	if err := sw.Export(ctx, rules); err != nil {
		return fmt.Errorf("failed to export parameters: %w", err)
	}

	return nil
}
