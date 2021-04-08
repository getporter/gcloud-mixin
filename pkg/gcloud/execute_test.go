package gcloud

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"

	"get.porter.sh/porter/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	test.TestMainWithMockedCommandHandlers(m)
}

func TestMixin_Execute(t *testing.T) {
	testcases := []struct {
		name        string
		file        string
		wantOutput  string
		wantCommand string
	}{
		{"install", "testdata/install-input.yaml", "",
			"gcloud --quiet compute config-ssh --format json --ssh-config-file ./gce-ssh-config --ssh-key-file ./gce-ssh-key"},
		{"upgrade", "testdata/upgrade-input.yaml", "",
			"gcloud --quiet compute instances update myinst --format json --update-labels color=blue,ready=true"},
		{"invoke", "testdata/invoke-input.yaml", "vms",
			"gcloud --quiet compute instances list --format json"},
		{"uninstall", "testdata/uninstall-input.yaml", "",
			"gcloud --quiet compute instances delete myinst --delete-disks all --format json"},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			m := NewTestMixin(t)

			m.Setenv(test.ExpectedCommandEnv, tc.wantCommand)
			mixinInputB, err := ioutil.ReadFile(tc.file)
			require.NoError(t, err)

			m.In = bytes.NewBuffer(mixinInputB)

			err = m.Execute()
			require.NoError(t, err, "execute failed")

			if tc.wantOutput == "" {
				outputs, _ := m.FileSystem.ReadDir("/cnab/app/porter/outputs")
				assert.Empty(t, outputs, "expected no outputs to be created")
			} else {
				wantPath := path.Join("/cnab/app/porter/outputs", tc.wantOutput)
				exists, _ := m.FileSystem.Exists(wantPath)
				assert.True(t, exists, "output file was not created %s", wantPath)
			}
		})
	}
}
