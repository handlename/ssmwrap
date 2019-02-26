package ssmwrap

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

type DefaultSSMConnector struct{}

func (c DefaultSSMConnector) FetchParameters(paths []string, recursive bool, retries int) (map[string]string, error) {
	client, err := newSSMClient(retries)
	if err != nil {
		return nil, err
	}
	return c.fetchParameters(client, paths, recursive)
}

func (c DefaultSSMConnector) fetchParameters(client *ssm.SSM, paths []string, recursive bool) (map[string]string, error) {
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

func (c DefaultSSMConnector) FetchParametersByNames(paths []string, retries int) (map[string]string, error) {
	client, err := newSSMClient(retries)
	if err != nil {
		return nil, err
	}
	return c.fetchParametersByNames(client, paths)
}

func (c DefaultSSMConnector) fetchParametersByNames(client *ssm.SSM, names []string) (map[string]string, error) {
	params := map[string]string{}
	if len(names) == 0 {
		return params, nil
	}

	input := &ssm.GetParametersInput{
		WithDecryption: aws.Bool(true),
	}
	for _, name := range names {
		input.Names = append(input.Names, aws.String(name))
	}

	output, err := client.GetParameters(input)
	if err != nil {
		return params, errors.Wrap(err, "failed to GetParameters")
	}

	for _, param := range output.Parameters {
		params[*param.Name] = *param.Value
	}

	return params, nil
}

func newSSMClient(retries int) (*ssm.SSM, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, "failed to start session")
	}

	// config.MaxRetries is defaults to -1
	if retries <= 0 {
		retries = -1
	}

	return ssm.New(sess, aws.NewConfig().WithMaxRetries(retries)), nil
}
