package ssmwrap

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type DefaultSSMConnector struct{}

func (c DefaultSSMConnector) fetchParametersByPaths(ctx context.Context, client *ssm.Client, paths []string, recursive bool) (map[string]string, error) {
	params := map[string]string{}
	if len(paths) == 0 {
		return params, nil
	}

	for _, path := range paths {
		nextToken := ""

		for {
			input := &ssm.GetParametersByPathInput{
				Path:           &path,
				Recursive:      aws.Bool(recursive),
				WithDecryption: aws.Bool(true),
			}

			if nextToken != "" {
				input.NextToken = aws.String(nextToken)
			}

			output, err := client.GetParametersByPath(ctx, input)
			if err != nil {
				return params, fmt.Errorf("failed to GetParametersByPath: %w", err)
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

func (c DefaultSSMConnector) fetchParametersByNames(ctx context.Context, client *ssm.Client, names []string) (map[string]string, error) {
	params := make(map[string]string, len(names))
	if len(names) == 0 {
		return params, nil
	}

	input := &ssm.GetParametersInput{
		WithDecryption: aws.Bool(true),
	}
	for _, name := range names {
		if _, exists := params[name]; exists { // discard duplication
			continue
		}
		input.Names = append(input.Names, name)
	}

	output, err := client.GetParameters(ctx, input)
	if err != nil {
		return params, fmt.Errorf("failed to GetParameters: %s", err)
	}

	for _, param := range output.Parameters {
		params[*param.Name] = *param.Value
	}

	return params, nil
}

func newSSMClient(ctx context.Context, retries int) (*ssm.Client, error) {
	opts := []func(*config.LoadOptions) error{}

	if 0 < retries {
		opts = append(opts, config.WithRetryMaxAttempts(retries))
	}

	conf, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load default aws config: %w", err)
	}

	return ssm.NewFromConfig(conf), nil
}
