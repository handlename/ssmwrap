package cli

import (
	"fmt"
	"io/fs"
	"strconv"
	"strings"

	"github.com/handlename/ssmwrap/internal/app"
)

type RuleFlags []app.Rule

func (f *RuleFlags) Set(value string) error {
	optLines := strings.Split(value, ",")
	opts := make(map[string]string, len(optLines))

	for _, opt := range optLines {
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			return f.Errorf(value, "invalid format")
		}

		opts[parts[0]] = parts[1]
	}

	rule := app.Rule{}

	if _, ok := opts["path"]; !ok {
		return f.Errorf(value, "`path` is required")
	}

	path, level, err := f.parsePath(opts["path"])
	if err != nil {
		return f.Errorf(value, err.Error())
	}

	rule.ParameterRule = app.ParameterRule{
		Path:  path,
		Level: level,
	}

	switch opts["type"] {
	case string(app.DestinationTypeEnv):
		rule.DestinationRule = app.DestinationRule{
			Type:           app.DestinationTypeEnv,
			To:             opts["to"],
			TypeEnvOptions: &app.DestinationTypeEnvOptions{},
		}

		if v, ok := opts["prefix"]; ok {
			rule.DestinationRule.TypeEnvOptions.Prefix = v
		}

		if v, ok := opts["entirepath"]; ok {
			entirePath, err := strconv.ParseBool(v)
			if err != nil {
				return f.Errorf(value, "invalid `entirepath`")
			}

			rule.DestinationRule.TypeEnvOptions.EntirePath = entirePath
		}
	case string(app.DestinationTypeFile):
		if _, ok := opts["to"]; !ok {
			return f.Errorf(value, "`to` is required for `type=file`")
		}

		// TODO: check if `to` is valid as file path

		rule.DestinationRule = app.DestinationRule{
			Type:            app.DestinationTypeFile,
			To:              opts["to"],
			TypeFileOptions: &app.DestinationTypeFileOptions{},
		}

		if modeStr, ok := opts["mode"]; ok {
			mode, err := strconv.ParseUint(modeStr, 8, 32)
			if err != nil {
				return f.Errorf(value, "invalid `mode`")
			}

			rule.DestinationRule.TypeFileOptions.Mode = fs.FileMode(mode)
		}

		if uidStr, ok := opts["uid"]; ok {
			uid, err := strconv.Atoi(uidStr)
			if err != nil {
				return f.Errorf(value, "invalid `uid`")
			}

			rule.DestinationRule.TypeFileOptions.Uid = uid
		}

		if gidStr, ok := opts["gid"]; ok {
			gid, err := strconv.Atoi(gidStr)
			if err != nil {
				return f.Errorf(value, "invalid `gid`")
			}

			rule.DestinationRule.TypeFileOptions.Gid = gid
		}
	default:
		return f.Errorf(value, "invalid `type`")
	}

	*f = append(*f, rule)

	return nil
}

func (f RuleFlags) parsePath(value string) (string, app.ParameterLevel, error) {
	if strings.HasSuffix(value, "/**/*") {
		return value[:len(value)-4], app.ParameterLevelAll, nil
	}

	if strings.HasSuffix(value, "/*") {
		return value[:len(value)-1], app.ParameterLevelUnder, nil
	}

	if strings.HasSuffix(value, "/") {
		return "", app.ParameterLevelStrict, fmt.Errorf("path must not end with `/`")
	}

	if !strings.HasPrefix(value, "/") {
		return "", app.ParameterLevelStrict, fmt.Errorf("path must start with `/`")
	}

	// TODO: validate if path not contains invalid characters

	return value, app.ParameterLevelStrict, nil
}

func (f RuleFlags) Errorf(value, format string, args ...interface{}) error {
	return fmt.Errorf("-rule "+value+": "+format, args...)
}
