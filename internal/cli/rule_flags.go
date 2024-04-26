package cli

import (
	"fmt"
	"io/fs"
	"regexp"
	"strconv"
	"strings"

	"github.com/handlename/ssmwrap/internal/app"
	"github.com/samber/lo"
)

var validPathRegexp = regexp.MustCompile(`^/[-_/a-zA-Z0-9]+((/\**)?/\*)?$`)

type RuleFlags struct {
	Rules []app.Rule
}

func (f RuleFlags) String() string {
	ss := make([]string, len(f.Rules)+2)
	ss = append(ss, "[")
	for _, r := range f.Rules {
		ss = append(ss, r.String())
	}
	ss = append(ss, "]")

	return strings.Join(ss, " ")
}

func (f *RuleFlags) Set(value string) error {
	opts, err := f.parseValue(value)
	if err != nil {
		return f.Errorf(value, err.Error())
	}

	rule, err := f.buildRule(opts)
	if err != nil {
		return f.Errorf(value, err.Error())
	}

	f.Rules = append(f.Rules, *rule)

	return nil
}

func (f RuleFlags) parseValue(value string) (map[string]string, error) {
	optLines := strings.Split(value, ",")
	opts := make(map[string]string, len(optLines))

	for _, opt := range optLines {
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid format")
		}

		opts[parts[0]] = parts[1]
	}

	return opts, nil
}

func (f RuleFlags) buildRule(opts map[string]string) (*app.Rule, error) {

	rule := &app.Rule{}

	if _, ok := opts["path"]; !ok {
		return nil, fmt.Errorf("`path` is required")
	}

	path, level, err := f.parsePath(opts["path"])
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	rule.ParameterRule = app.ParameterRule{
		Path:  path,
		Level: level,
	}

	switch opts["type"] {
	case string(app.DestinationTypeEnv):
		if err := f.checkOptionsCombinations(app.DestinationTypeEnv, opts); err != nil {
			return nil, err
		}

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
				return nil, fmt.Errorf("invalid `entirepath`")
			}

			rule.DestinationRule.TypeEnvOptions.EntirePath = entirePath
		}
	case string(app.DestinationTypeFile):
		if err := f.checkOptionsCombinations(app.DestinationTypeFile, opts); err != nil {
			return nil, err
		}

		if _, ok := opts["to"]; !ok {
			return nil, fmt.Errorf("`to` is required for `type=file`")
		}

		if level != app.ParameterLevelStrict {
			return nil, fmt.Errorf("`path` end with `/*` or `/**/*` is not allowed for `type=file`")
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
				return nil, fmt.Errorf("invalid `mode`")
			}

			rule.DestinationRule.TypeFileOptions.Mode = fs.FileMode(mode)
		}

		if uidStr, ok := opts["uid"]; ok {
			uid, err := strconv.Atoi(uidStr)
			if err != nil {
				return nil, fmt.Errorf("invalid `uid`")
			}

			rule.DestinationRule.TypeFileOptions.Uid = uid
		}

		if gidStr, ok := opts["gid"]; ok {
			gid, err := strconv.Atoi(gidStr)
			if err != nil {
				return nil, fmt.Errorf("invalid `gid`")
			}

			rule.DestinationRule.TypeFileOptions.Gid = gid
		}
	default:
		return nil, fmt.Errorf("invalid `type`")
	}

	return rule, nil
}

func (f RuleFlags) parsePath(value string) (string, app.ParameterLevel, error) {
	if !validPathRegexp.MatchString(value) {
		return "", app.ParameterLevelStrict, fmt.Errorf("invalid `path` format")
	}

	if strings.HasSuffix(value, "/**/*") {
		return value[:len(value)-4], app.ParameterLevelAll, nil
	}

	if strings.HasSuffix(value, "/*") {
		return value[:len(value)-1], app.ParameterLevelUnder, nil
	}

	return value, app.ParameterLevelStrict, nil
}

func (f RuleFlags) checkOptionsCombinations(t app.DestinationType, opts map[string]string) error {
	for _, key := range lo.Keys(opts) {
		switch key {
		case "prefix":
			if t != app.DestinationTypeEnv {
				return f.Errorf(key, "`prefix` is only allowed for `type=env`")
			}
		case "entirepath":
			if t != app.DestinationTypeEnv {
				return f.Errorf(key, "`entirepath` is only allowed for `type=env`")
			}

			if _, ok := opts["to"]; ok {
				return f.Errorf(key, "can't use `to` with `entirepath` in same time")
			}
		case "mode":
			if t != app.DestinationTypeFile {
				return f.Errorf(key, "`mode` is only allowed for `type=file`")
			}
		case "uid":
			if t != app.DestinationTypeFile {
				return f.Errorf(key, "`uid` is only allowed for `type=file`")
			}
		case "gid":
			if t != app.DestinationTypeFile {
				return f.Errorf(key, "`gid` is only allowed for `type=file`")
			}
		}
	}

	return nil
}

func (f RuleFlags) Errorf(value, format string, args ...interface{}) error {
	return fmt.Errorf("-rule "+value+": "+format, args...)
}
