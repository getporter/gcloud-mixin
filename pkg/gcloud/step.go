package gcloud

import (
	"github.com/pkg/errors"
)

type Step struct {
	Description string   `yaml:"description"`
	Groups      Groups   `yaml:"groups"`
	Command     string   `yaml:"command"`
	Arguments   []string `yaml:"arguments,omitempty"`
	Flags       Flags    `yaml:"flags,omitempty"`
	Outputs     []Output `yaml:"outputs,omitempty"`
}

type Groups []string

func (groups *Groups) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var groupMap interface{}
	err := unmarshal(&groupMap)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal yaml into Step.Groups")
	}

	switch t := groupMap.(type) {
	case string:
		*groups = append(*groups, t)
	case []interface{}:
		for i := range t {
			group, ok := t[i].(string)
			if !ok {
				return errors.Errorf("invalid yaml type for group item: %T", t[i])
			}
			*groups = append(*groups, group)
		}
	default:
		return errors.Errorf("invalid yaml type for group item: %T", t)
	}

	return nil
}

type Flags []Flag

func (flags Flags) Len() int {
	return len(flags)
}

func (flags Flags) Swap(i, j int) {
	flags[i], flags[j] = flags[j], flags[i]
}

func (flags Flags) Less(i, j int) bool {
	return flags[i].Name < flags[j].Name
}

type Flag struct {
	Name   string
	Values []string
}

func NewFlag(name string, values ...string) Flag {
	f := Flag{
		Name:   name,
		Values: make([]string, len(values)),
	}
	copy(f.Values, values)
	return f
}

func (flags *Flags) UnmarshalYAML(unmarshal func(interface{}) error) error {
	flagMap := map[interface{}]interface{}{}
	err := unmarshal(&flagMap)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal yaml into Step.Flags")
	}

	*flags = make(Flags, 0, len(flagMap))
	for k, v := range flagMap {
		f := Flag{}
		f.Name = k.(string)

		switch t := v.(type) {
		case string:
			f.Values = make([]string, 1)
			f.Values[0] = t
		case []interface{}:
			f.Values = make([]string, len(t))
			for i := range t {
				iv, ok := t[i].(string)
				if !ok {
					return errors.Errorf("invalid yaml type for flag %s: %T", f.Name, t[i])
				}
				f.Values[i] = iv
			}
		default:
			return errors.Errorf("invalid yaml type for flag %s: %T", f.Name, t)
		}

		*flags = append(*flags, f)
	}

	return nil
}

func (flags Flags) MarshalYAML() (interface{}, error) {
	result := make(map[string]interface{}, len(flags))

	for _, flag := range flags {
		if len(flag.Values) == 1 {
			result[flag.Name] = flag.Values[0]
		} else {
			result[flag.Name] = flag.Values
		}
	}

	return result, nil
}

type Output struct {
	Name     string `yaml:"name"`
	JsonPath string `yaml:"jsonPath"`
}
