package gcloud

import (
	"get.porter.sh/porter/pkg/exec/builder"
	"github.com/pkg/errors"
)

var _ builder.ExecutableAction = Action{}
var _ builder.BuildableAction = Action{}

type Action struct {
	Name  string
	Steps []Step // using UnmarshalYAML so that we don't need a custom type per action
}

// MakeSteps builds a slice of Steps for data to be unmarshaled into.
func (a Action) MakeSteps() interface{} {
	return &[]Step{}
}

// UnmarshalYAML takes any yaml in this form
// ACTION:
// - gcloud: ...
// and puts the steps into the Action.Steps field
func (a *Action) UnmarshalYAML(unmarshal func(interface{}) error) error {
	results, err := builder.UnmarshalAction(unmarshal, a)
	if err != nil {
		return err
	}

	for actionName, action := range results {
		a.Name = actionName
		for _, result := range action {
			step := result.(*[]Step)
			a.Steps = append(a.Steps, *step...)
		}
		break // There is only 1 action
	}
	return nil
}

func (a Action) GetSteps() []builder.ExecutableStep {
	steps := make([]builder.ExecutableStep, len(a.Steps))
	for i := range a.Steps {
		steps[i] = a.Steps[i]
	}

	return steps
}

type Step struct {
	Instruction `yaml:"gcloud"`
}

var _ builder.ExecutableStep = Step{}
var _ builder.StepWithOutputs = Step{}
var _ builder.SuppressesOutput = Step{}

type Instruction struct {
	Description    string        `yaml:"description"`
	Groups         Groups        `yaml:"groups"`
	Command        string        `yaml:"command"`
	Arguments      []string      `yaml:"arguments,omitempty"`
	Flags          builder.Flags `yaml:"flags,omitempty"`
	Outputs        []Output      `yaml:"outputs,omitempty"`
	SuppressOutput bool          `yaml:"suppress-output,omitempty"`
}

func (s Step) GetCommand() string {
	return "gcloud"
}

func (s Step) GetArguments() []string {
	args := make([]string, 0, len(s.Arguments)+len(s.Groups)+2)

	// Always be in non-interactive mode, must be specified immediately after gcloud
	args = append(args, "--quiet")

	// Specify the gcloud group(s) and command
	args = append(args, s.Groups...)
	args = append(args, s.Command)

	// Append the positional arguments
	args = append(args, s.Arguments...)

	return args
}

func (s Step) GetFlags() builder.Flags {
	// Always request json formatted output
	return append(s.Flags, builder.NewFlag("format", "json"))
}

func (s Step) GetOutputs() []builder.Output {
	outputs := make([]builder.Output, len(s.Outputs))
	for i := range s.Outputs {
		outputs[i] = s.Outputs[i]
	}
	return outputs
}

func (s Step) SuppressesOutput() bool {
	return s.SuppressOutput
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

var _ builder.OutputJsonPath = Output{}

type Output struct {
	Name     string `yaml:"name"`
	JsonPath string `yaml:"jsonPath"`
}

func (o Output) GetName() string {
	return o.Name
}

func (o Output) GetJsonPath() string {
	return o.JsonPath
}
