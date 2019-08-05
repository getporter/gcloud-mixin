package gcloud

import (
	"io/ioutil"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func TestFlags_Sort(t *testing.T) {
	flags := Flags{
		NewFlag("b", "1"),
		NewFlag("a", "2"),
		NewFlag("c", "3"),
	}

	sort.Sort(flags)

	assert.Equal(t, "a", flags[0].Name)
	assert.Equal(t, "b", flags[1].Name)
	assert.Equal(t, "c", flags[2].Name)
}

func TestMixin_UnmarshalStep(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/step-input.yaml")
	require.NoError(t, err)

	var step Steps
	err = yaml.Unmarshal(b, &step)
	require.NoError(t, err)

	assert.Equal(t, "Create VM", step.Description)
	assert.Equal(t, Groups{"compute", "instances"}, step.Groups)
	assert.Equal(t, "create", step.Command)

	assert.Equal(t, []string{"myinst"}, step.Arguments)

	sort.Sort(step.Flags)
	assert.Equal(t, Flags{
		NewFlag("env", "CLIENT_VERSION=1.0.0", "SERVER_VERSION=1.1.0"),
		NewFlag("hostname", "example.com"),
		NewFlag("labels", "FOO=BAR,STUFF=THINGS"),
		NewFlag("quiet", "true")}, step.Flags)
}

func TestMixin_UnmarshalInvalidStep(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/step-input-invalid.yaml")
	require.NoError(t, err)

	var step Steps
	err = yaml.Unmarshal(b, &step)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid yaml type for flag env")
}
