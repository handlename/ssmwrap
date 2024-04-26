package app

import (
	"fmt"
	"log/slog"
	"strings"
)

type Rule struct {
	ParameterRule   ParameterRule
	DestinationRule DestinationRule
}

func (r Rule) String() string {
	ss := []string{
		"path=" + r.ParameterRule.String(),
		"type=" + string(r.DestinationRule.Type),
	}

	if r.DestinationRule.To != "" {
		ss = append(ss, "to="+r.DestinationRule.To)
	}

	switch r.DestinationRule.Type {
	case DestinationTypeEnv:
		ss = append(ss, r.DestinationRule.TypeEnvOptions.String())
	case DestinationTypeFile:
		ss = append(ss, r.DestinationRule.TypeFileOptions.String())
	}

	return strings.Join(ss, ",")
}

func (r Rule) Execute(store ParameterStore) error {
	params, err := store.Retrieve(r.ParameterRule.Path, r.ParameterRule.Level)
	if err != nil {
		return fmt.Errorf("failed to retrieve parameters: %w", err)
	}

	var ex Exporter
	switch r.DestinationRule.Type {

	case DestinationTypeEnv:
		if r.DestinationRule.TypeEnvOptions == nil {
			return fmt.Errorf("TypeEnvOptions is required for DestinationTypeEnv")
		}

		var envName string
		if r.DestinationRule.TypeEnvOptions.Prefix != "" {
			envName = r.DestinationRule.TypeEnvOptions.Prefix + "_"
		}

		if r.DestinationRule.TypeEnvOptions.EntirePath {
			envName += strings.ReplaceAll(r.ParameterRule.Path, "/", "_")
		} else {
			parts := strings.Split(r.ParameterRule.Path, "/")
			envName += parts[len(parts)-1]
		}

		envName = strings.ToUpper(envName)
		ex = NewEnvExporter(envName)
	case DestinationTypeFile:
		if r.DestinationRule.TypeFileOptions == nil {
			return fmt.Errorf("TypeFileOption is required for DestinationTypeFile")
		}

		e := NewFileExporter(r.DestinationRule.To)

		if r.DestinationRule.TypeFileOptions.Mode != 0 {
			e.Mode = r.DestinationRule.TypeFileOptions.Mode
		}

		if r.DestinationRule.TypeFileOptions.Uid != 0 {
			e.Uid = r.DestinationRule.TypeFileOptions.Uid
		}

		if r.DestinationRule.TypeFileOptions.Gid != 0 {
			e.Gid = r.DestinationRule.TypeFileOptions.Gid
		}

		ex = e
	default:
		return fmt.Errorf("invalid destination type: %s", r.DestinationRule.Type)
	}

	for _, p := range params {
		slog.Debug(
			"exporting parameter",
			slog.String("type", string(r.DestinationRule.Type)),
			slog.String("address", ex.Address()),
		)

		if err := ex.Export(p.Value); err != nil {
			return fmt.Errorf("failed to export parameter to %s: %w", r.DestinationRule.To, err)
		}
	}

	return nil
}
