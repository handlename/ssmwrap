package app

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/samber/lo"
)

type ParameterStore struct {
	client *ssm.Client
	conn   SSMConnector

	Parameters []Parameter
}

func NewParameterStore(client *ssm.Client, conn SSMConnector) *ParameterStore {
	return &ParameterStore{
		client: client,
		conn:   conn,
	}
}

func (c *ParameterStore) Store(ctx context.Context, rules []ParameterRule) error {
	c.Parameters = []Parameter{}

	paths := map[ParameterLevel][]string{
		ParameterLevelStrict: {},
		ParameterLevelUnder:  {},
		ParameterLevelAll:    {},
	}

	for _, rule := range rules {
		switch rule.Level {
		case ParameterLevelStrict:
			paths[ParameterLevelStrict] = append(paths[ParameterLevelStrict], rule.Path)
		case ParameterLevelUnder:
			paths[ParameterLevelUnder] = append(paths[ParameterLevelUnder], rule.Path)
		case ParameterLevelAll:
			paths[ParameterLevelAll] = append(paths[ParameterLevelAll], rule.Path)
		default:
			slog.Warn("invalid ParameterRule path level", slog.Int("level", int(rule.Level)))
		}
	}

	add := func(params map[string]string) {
		for key, value := range params {
			c.Parameters = append(c.Parameters, Parameter{
				Path:  key,
				Value: value,
			})
		}
	}

	if p, err := c.conn.fetchParametersByNames(ctx, c.client, paths[ParameterLevelStrict]); err != nil {
		return fmt.Errorf("failed to fetch parameters from SSM by strict paths %v: %w", paths[ParameterLevelStrict], err)
	} else {
		add(p)
	}

	if p, err := c.conn.fetchParametersByPaths(ctx, c.client, paths[ParameterLevelUnder], false); err != nil {
		return fmt.Errorf("failed to fetch parameters from SSM by just under paths %v: %w", paths[ParameterLevelUnder], err)
	} else {
		add(p)
	}

	if p, err := c.conn.fetchParametersByPaths(ctx, c.client, paths[ParameterLevelAll], true); err != nil {
		return fmt.Errorf("failed to fetch parameters from SSM by under paths recursively %v: %w", paths[ParameterLevelAll], err)
	} else {
		add(p)
	}

	return nil
}

func (c ParameterStore) Retrieve(path string, level ParameterLevel) ([]Parameter, error) {
	switch level {
	case ParameterLevelStrict:
		if param := c.FindByName(path); param == nil {
			return []Parameter{}, nil
		} else {
			return []Parameter{*param}, nil
		}
	case ParameterLevelUnder:
		return c.SearchByPath(path, false), nil
	case ParameterLevelAll:
		return c.SearchByPath(path, true), nil
	default:
		return nil, fmt.Errorf("invalid ParameterLevel: %d", level)
	}
}

func (c ParameterStore) FindByName(name string) *Parameter {
	params := lo.Filter(c.Parameters, func(p Parameter, _ int) bool {
		return p.Path == name
	})
	if len(params) == 0 {
		return nil
	}
	if 2 <= len(params) {
		slog.Warn("found multiple parameters with the same name", slog.String("name", name))
	}

	return &params[0]
}

func (c ParameterStore) SearchByPath(path string, recursive bool) []Parameter {
	return lo.Filter(c.Parameters, func(p Parameter, _ int) bool {
		if !strings.HasPrefix(p.Path, path) {
			return false
		}

		if !recursive {
			rest := strings.Replace(p.Path, path, "", 1)
			if strings.Contains(rest, "/") {
				return false
			}
		}

		return true
	})
}
