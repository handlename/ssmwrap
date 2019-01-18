package ssmwrap

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

type DefaultSSMConnector struct{}

func (c DefaultSSMConnector) FetchParameters(paths []string, recursive bool, retries int) (map[string]string, error) {
	params := map[string]string{}

	sess, err := session.NewSession()
	if err != nil {
		return params, errors.Wrap(err, "failed to start session")
	}

	// config.MaxRetries is defaults to -1
	if retries <= 0 {
		retries = -1
	}

	client := ssm.New(sess, aws.NewConfig().WithMaxRetries(retries))

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
