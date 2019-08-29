package gcloud

import (
	"io/ioutil"
	"sort"
	"testing"

	"github.com/deislabs/porter/pkg/exec/builder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func TestMixin_UnmarshalInstallAction(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/install-input.yaml")
	require.NoError(t, err)

	var action Action
	err = yaml.Unmarshal(b, &action)
	require.NoError(t, err)

	require.Equal(t, 1, len(action.Steps))
	step := action.Steps[0]

	assert.Equal(t, "Configure SSH", step.Description)
	assert.Equal(t, Groups{"compute"}, step.Groups)
	assert.Equal(t, "config-ssh", step.Command)

	sort.Sort(step.Flags)
	assert.Equal(t, builder.Flags{
		builder.NewFlag("ssh-config-file", "./gce-ssh-config"),
		builder.NewFlag("ssh-key-file", "./gce-ssh-key")}, step.Flags)
}

func TestMixin_UnmarshalUpgradeAction(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/upgrade-input.yaml")
	require.NoError(t, err)

	var action Action
	err = yaml.Unmarshal(b, &action)
	require.NoError(t, err)

	require.Equal(t, 1, len(action.Steps))
	step := action.Steps[0]

	assert.Equal(t, "Tag VM", step.Description)
	require.Empty(t, step.Outputs)

	assert.Equal(t, Groups{"compute", "instances"}, step.Groups)
	assert.Equal(t, "update", step.Command)

	assert.Equal(t, []string{"myinst"}, step.Arguments)

	sort.Sort(step.Flags)
	assert.Equal(t, builder.Flags{
		builder.NewFlag("update-labels", "color=blue,ready=true")}, step.Flags)
}

func TestMixin_UnmarshalInvokeAction(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/invoke-input.yaml")
	require.NoError(t, err)

	var action Action
	err = yaml.Unmarshal(b, &action)
	require.NoError(t, err)

	require.Equal(t, 1, len(action.Steps))
	step := action.Steps[0]

	assert.Equal(t, "List VMs", step.Description)
	assert.Equal(t, Groups{"compute", "instances"}, step.Groups)
	assert.Equal(t, "list", step.Command)
	assert.Empty(t, step.Arguments)
	assert.Equal(t, []Output{{Name: "vms", JsonPath: "$[*].id"}}, step.Outputs)

	assert.Empty(t, step.Flags)
}

func TestMixin_UnmarshalUninstallAction(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/uninstall-input.yaml")
	require.NoError(t, err)

	var action Action
	err = yaml.Unmarshal(b, &action)
	require.NoError(t, err)

	require.Equal(t, 1, len(action.Steps))
	step := action.Steps[0]

	assert.Equal(t, "Deprovision VM", step.Description)
	require.Empty(t, step.Outputs)

	assert.Equal(t, Groups{"compute", "instances"}, step.Groups)
	assert.Equal(t, "delete", step.Command)

	assert.Equal(t, []string{"myinst"}, step.Arguments)

	sort.Sort(step.Flags)
	assert.Equal(t, builder.Flags{
		builder.NewFlag("delete-disks", "all")}, step.Flags)
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
	assert.Equal(t, builder.Flags{
		builder.NewFlag("env", "CLIENT_VERSION=1.0.0", "SERVER_VERSION=1.1.0"),
		builder.NewFlag("hostname", "example.com"),
		builder.NewFlag("labels", "FOO=BAR,STUFF=THINGS"),
		builder.NewFlag("quiet", "true")}, step.Flags)
}

func TestMixin_UnmarshalInvalidStep(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/step-input-invalid.yaml")
	require.NoError(t, err)

	var step Steps
	err = yaml.Unmarshal(b, &step)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid yaml type for flag env")
}
