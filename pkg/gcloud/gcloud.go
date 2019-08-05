package gcloud

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type Action struct {
	Steps []Steps // using UnmarshalYAML so that we don't need a custom type per action
}

// UnmarshalYAML takes any yaml in this form
// ACTION:
// - gcloud: ...
// and puts the steps into the Action.Steps field
func (a *Action) UnmarshalYAML(unmarshal func(interface{}) error) error {
	actionMap := map[interface{}][]interface{}{}
	err := unmarshal(&actionMap)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal yaml into an action map of gcloud steps")
	}

	for _, stepMaps := range actionMap {
		b, err := yaml.Marshal(stepMaps)
		if err != nil {
			return err
		}

		var steps []Steps
		err = yaml.Unmarshal(b, &steps)
		if err != nil {
			return err
		}

		a.Steps = append(a.Steps, steps...)
	}

	return nil
}

type Steps struct {
	Step `yaml:"gcloud"`
}

func (m *Mixin) Execute() error {
	payload, err := m.getPayloadData()
	if err != nil {
		return err
	}

	var action Action
	err = yaml.Unmarshal(payload, &action)
	if err != nil {
		return err
	}
	if len(action.Steps) != 1 {
		return errors.Errorf("expected a single step, but got %d", len(action.Steps))
	}
	step := action.Steps[0]

	// Always output json so that we can query it for outputs afterwards
	step.Flags = append(step.Flags, NewFlag("format", "json"))

	fmt.Fprintf(m.Out, "Starting operation: %s\n", step.Description)

	args := make([]string, 0, 2+len(step.Arguments)+len(step.Flags)*2)
	// Always be in non-interactive mode
	args = append(args, "--quiet")

	// Specify the gcloud groups and command to run
	args = append(args, step.Groups...)
	args = append(args, step.Command)

	// Append the positional arguments
	for _, arg := range step.Arguments {
		args = append(args, arg)
	}

	// Append the flags to the argument list
	sort.Sort(step.Flags)
	for _, flag := range step.Flags {
		for _, value := range flag.Values {
			args = append(args, fmt.Sprintf("--%s", flag.Name))
			args = append(args, value)
		}
	}

	cmd := m.NewCommand("gcloud", args...)
	cmd.Stdout = m.Out
	cmd.Stderr = m.Err

	err = cmd.Start()
	if err != nil {
		prettyCmd := fmt.Sprintf("%s %s", cmd.Path, strings.Join(cmd.Args, " "))
		return errors.Wrap(err, fmt.Sprintf("couldn't run command %s", prettyCmd))
	}

	err = cmd.Wait()

	if err != nil {
		prettyCmd := fmt.Sprintf("%s %s", cmd.Path, strings.Join(cmd.Args, " "))
		return errors.Wrap(err, fmt.Sprintf("error running command %s", prettyCmd))
	}
	fmt.Fprintf(m.Out, "Finished operation: %s\n", step.Description)

	return nil
}
