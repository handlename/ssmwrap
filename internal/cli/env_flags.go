package cli

type EnvFlags struct {
	RuleFlags
}

func (f *EnvFlags) Set(value string) error {
	opts, err := f.parseValue(value)
	if err != nil {
		return f.Errorf(value, err.Error())
	}

	opts["type"] = "env"

	rule, err := f.buildRule(opts)
	if err != nil {
		return f.Errorf(value, err.Error())
	}

	f.Rules = append(f.Rules, *rule)

	return nil
}
