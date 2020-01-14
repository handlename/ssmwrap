package ssmwrap

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type DefaultSSMConnector struct{}

func (c DefaultSSMConnector) fetchParametersByPaths(client *ssm.SSM, paths []string, recursive bool) (map[string]string, error) {
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
				input.SetNextToken(nextToken)
			}

			output, err := client.GetParametersByPath(input)
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

func (c DefaultSSMConnector) fetchParametersByNames(client *ssm.SSM, names []string) (map[string]string, error) {
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
		params[name] = ""
		input.Names = append(input.Names, aws.String(name))
	}

	output, err := client.GetParameters(input)
	if err != nil {
		return params, fmt.Errorf("failed to GetParameters: %s", err)
	}

	for _, param := range output.Parameters {
		params[*param.Name] = *param.Value
	}

	return params, nil
}

func newSSMClient(retries int) (*ssm.SSM, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}

	// config.MaxRetries is defaults to -1
	if retries <= 0 {
		retries = -1
	}

	return ssm.New(sess, aws.NewConfig().WithMaxRetries(retries)), nil
}
