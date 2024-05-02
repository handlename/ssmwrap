package ssmwrap

import (
	"context"
	"fmt"

	"github.com/handlename/ssmwrap/internal/app"
)

type ExportOptions struct {
	Retries int
}

type ExportRule struct {
	// Path of parameter store.
	// If `path` ends with no-slash character, only the value of the path will be exported.
	// If `path` ends with `/**/*`, all values under the path will be exported.
	// If `path` ends with `/*`, only top level values under the path will be exported.
	Path string

	// Prefix for exported environment variable.
	Prefix string

	// UseEntirePath is flag if export entire path as environment variables name.
	// If true, all values under the path will be exported. (/path/to/param -> PATH_TO_PARAM)
	// If false, only top level values under the path will be exported. (/path/to/param -> PARAM)
	UseEntirePath bool
}

// Export fetches parameters from SSM and export those to environment variables.
// This is for use ssmwrap as a library.
func Export(ctx context.Context, ers []ExportRule, options ExportOptions) error {
	rules := make([]app.Rule, 0, len(ers))

	for _, er := range ers {
		pr, err := app.NewParameterRule(er.Path)
		if err != nil {
			return fmt.Errorf("failed to create ParameterRule: %w", err)
		}

		rules = append(rules, app.Rule{
			ParameterRule: *pr,
			DestinationRule: app.DestinationRule{
				Type: app.DestinationTypeEnv,
				TypeEnvOptions: &app.DestinationTypeEnvOptions{
					Prefix:     er.Prefix,
					EntirePath: er.UseEntirePath,
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
