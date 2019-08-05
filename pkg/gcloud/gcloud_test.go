package gcloud

import (
	"bytes"
	"io/ioutil"
	"os"
	"sort"
	"testing"

	"github.com/deislabs/porter/pkg/test"
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
	assert.Equal(t, Flags{
		NewFlag("ssh-config-file", "./gce-ssh-config"),
		NewFlag("ssh-key-file", "./gce-ssh-key")}, step.Flags)
}

func TestMixin_UnmarshalUpgradelAction(t *testing.T) {
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
	assert.Equal(t, Flags{
		NewFlag("update-labels", "color=blue,ready=true")}, step.Flags)
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
	assert.Equal(t, Flags{
		NewFlag("delete-disks", "all")}, step.Flags)
}

func TestMain(m *testing.M) {
	test.TestMainWithMockedCommandHandlers(m)
}

func TestMixin_Execute(t *testing.T) {
	m := NewTestMixin(t)

	testcases := []struct {
		name        string
		file        string
		wantCommand string
	}{
		{"install", "testdata/install-input.yaml",
			"gcloud --quiet compute config-ssh --format json --ssh-config-file ./gce-ssh-config --ssh-key-file ./gce-ssh-key"},
		{"upgrade", "testdata/upgrade-input.yaml",
			"gcloud --quiet compute instances update myinst --format json --update-labels color=blue,ready=true"},
		{"invoke", "testdata/invoke-input.yaml",
			"gcloud --quiet compute instances list --format json"},
		{"uninstall", "testdata/uninstall-input.yaml",
			"gcloud --quiet compute instances delete myinst --delete-disks all --format json"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv(test.ExpectedCommandEnv, tc.wantCommand)
			mixinInputB, err := ioutil.ReadFile(tc.file)
			require.NoError(t, err)

			m.In = bytes.NewBuffer(mixinInputB)

			err = m.Execute()
			require.NoError(t, err, "execute failed")
		})
	}
}
