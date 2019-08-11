package gcloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/PaesslerAG/jsonpath"
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

	// Always request json formatted output
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

	// Split stdout in a buffer and also send it back to porter
	outputB := &bytes.Buffer{}
	cmd.Stdout = io.MultiWriter(m.Out, outputB)
	cmd.Stderr = m.Err

	prettyCmd := fmt.Sprintf("%s %s", cmd.Path, strings.Join(cmd.Args, " "))
	if m.Debug {
		fmt.Fprintln(m.Out, prettyCmd)
	}

	err = cmd.Start()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("couldn't run command %s", prettyCmd))
	}

	err = cmd.Wait()

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error running command %s", prettyCmd))
	}
	fmt.Fprintf(m.Out, "Finished operation: %s\n", step.Description)

	m.processOutputs(step.Step, outputB)

	return nil
}

func (m *Mixin) processOutputs(step Step, outputB *bytes.Buffer) error {
	if len(step.Outputs) == 0 {
		return nil
	}

	var outputJson interface{}
	err := json.Unmarshal(outputB.Bytes(), &outputJson)
	if err != nil {
		return errors.Wrapf(err, "error unmarshaling json %s", outputB.String())
	}

	for _, output := range step.Outputs {
		value, err := jsonpath.Get(output.JsonPath, outputJson)
		if err != nil {
			return errors.Wrapf(err, "error evaluating jsonpath %q for output %q against %s", output.JsonPath, output.Name, outputB.String())
		}
		valueB, err := json.Marshal(value)
		if err != nil {
			return errors.Wrapf(err, "error marshaling jsonpath result %v for output %q", valueB, output.Name)
		}
		err = m.WriteMixinOutputToFile(output.Name, valueB)
		if err != nil {
			return errors.Wrapf(err, "error writing mixin output for %q", output.Name)
		}
	}

	return nil
}
