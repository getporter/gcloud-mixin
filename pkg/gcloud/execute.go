package gcloud

import (
	"bytes"
	"encoding/json"

	"github.com/PaesslerAG/jsonpath"
	"github.com/deislabs/porter/pkg/exec/builder"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func (m *Mixin) loadAction() (*Action, error) {
	var action Action
	err := builder.LoadAction(m.Context, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &action)
		return &action, err
	})
	return &action, err
}

func (m *Mixin) Execute() error {
	action, err := m.loadAction()
	if err != nil {
		return err
	}

	output, err := builder.ExecuteSingleStepAction(m.Context, action)
	if err != nil {
		return err
	}
	step := action.Steps[0]
	return m.processOutputs(step.Step, output)
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
