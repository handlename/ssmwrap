//go:build integrated

package ssmwrap_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrated(t *testing.T) {
	RootDir := "."
	// PathPrefixPlaceholder := ""+envPrefix
	// EnvPrefixPlaceholder := "__ENV_PREFIX__"

	ctx := t.Context()
	pathPrefix := "/test/ssmwrap"
	envPrefix := fmt.Sprintf("%d_", time.Now().UnixNano())

	// Prepare SSM Params

	params := []struct {
		name  string
		value string
	}{
		{
			name:  pathPrefix + "/foo",
			value: "foo_value",
		},
		{
			name:  pathPrefix + "/foo/bar/bar1",
			value: "foo_bar_bar1_value",
		},
		{
			name:  pathPrefix + "/foo/bar/bar2",
			value: "foo_bar_bar2_value",
		},
		{
			name:  pathPrefix + "/foo/bar/buzz/buzz1",
			value: "foo_bar_buzz_buzz1_value",
		},
		{
			name:  pathPrefix + "/hoge",
			value: "hoge_value",
		},
	}

	cfg := lo.Must(config.LoadDefaultConfig(ctx))
	client := ssm.NewFromConfig(cfg)

	for _, p := range params {
		input := &ssm.PutParameterInput{
			Name:      aws.String(p.name),
			Value:     aws.String(p.value),
			Type:      types.ParameterTypeString,
			Overwrite: aws.Bool(true),
		}

		lo.Must(client.PutParameter(ctx, input))
		t.Logf("SSM Param %s is created or updated", p.name)
	}

	// Define test cases

	tests := []struct {
		flags        []string
		expectedEnvs []string // Key=Value
	}{
		{
			flags: []string{"-env", "path=" + pathPrefix + "/foo,prefix=" + envPrefix},
			expectedEnvs: []string{
				envPrefix + "FOO=foo_value",
			},
		},
		{
			flags: []string{"-env", "path=" + pathPrefix + "/foo/bar/*,prefix=" + envPrefix},
			expectedEnvs: []string{
				envPrefix + "BAR1=foo_bar_bar1_value",
				envPrefix + "BAR2=foo_bar_bar2_value",
			},
		},
		{
			flags: []string{"-env", "path=" + pathPrefix + "/foo/bar/**/*,prefix=" + envPrefix},
			expectedEnvs: []string{
				envPrefix + "BAR1=foo_bar_bar1_value",
				envPrefix + "BAR2=foo_bar_bar2_value",
				envPrefix + "BUZZ1=foo_bar_buzz_buzz1_value",
			},
		},
	}

	// Run test cases

	for _, tt := range tests {
		t.Run(strings.Join(tt.flags, " "), func(t *testing.T) {
			// Build command args
			args := []string{"run", "cmd/ssmwrap/main.go"}
			args = append(args, tt.flags...)
			args = append(args, []string{"--", "env"}...)

			// Prepare command
			cmd := exec.Command("go", args...)
			cmd.Dir = RootDir
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			// Run!
			require.NoError(t, cmd.Run())
			if stderr.String() != "" {
				t.Logf("logs:\n%s", stderr.String())
			}
			gotEnvs := strings.Split(stdout.String(), "\n")

			// Check
			for _, expected := range tt.expectedEnvs {
				assert.Contains(t, gotEnvs, expected)
			}
			assert.Equal(t,
				len(tt.expectedEnvs),
				len(lo.Filter(gotEnvs, func(line string, _ int) bool {
					return strings.HasPrefix(line, envPrefix)
				})),
			)
		})
	}
}
