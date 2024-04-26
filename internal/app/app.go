package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/samber/lo"
)

type SSMWrap struct {
	// Retry limit to request to SSM.
	Retries int

	// Command and arguments to run.
	Command []string
}

func NewSSMWrap() *SSMWrap {
	return &SSMWrap{
		Retries: 3,
	}
}

func (s *SSMWrap) Run(ctx context.Context, rules []Rule, command []string) error {
	slog.DebugContext(ctx, fmt.Sprintf("start to process %d rules", len(rules)))

	ssmClient, err := s.ssmClient(ctx)
	if err != nil {
		return err
	}

	// store related ssm params

	slog.DebugContext(ctx, "start to store parameters")

	store := NewParameterStore(ssmClient, DefaultSSMConnector{})
	if err := store.Store(ctx, lo.Map(rules, func(r Rule, _ int) ParameterRule {
		return r.ParameterRule
	})); err != nil {
		return fmt.Errorf("failed to refresh parameters: %w", err)
	}

	slog.DebugContext(ctx, fmt.Sprintf("%d parameters stored successfully", len(store.Parameters)))

	// execute rules

	for _, r := range rules {
		slog.DebugContext(ctx, "executing rule", slog.String("rule", r.String()))

		if err := r.Execute(*store); err != nil {
			return fmt.Errorf("failed to execute rule %s: %w", r, err)
		}
	}

	slog.DebugContext(ctx, fmt.Sprintf("%d rules processed successfully", len(rules)))

	// execute command

	if len(command) == 0 {
		return fmt.Errorf("command required")
	}

	bin, err := exec.LookPath(command[0])
	if err != nil {
		return fmt.Errorf("command is not executable %s: %w", command[0], err)
	}

	return syscall.Exec(bin, command, os.Environ())
}

func (s SSMWrap) ssmClient(ctx context.Context) (*ssm.Client, error) {
	opts := []func(*config.LoadOptions) error{}

	if 0 < s.Retries {
		opts = append(opts, config.WithRetryMaxAttempts(s.Retries))
	}

	conf, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load default aws config: %w", err)
	}

	return ssm.NewFromConfig(conf), nil
}
