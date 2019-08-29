package gcloud

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/deislabs/porter/pkg/test"
	"github.com/stretchr/testify/require"
)

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

	defer os.Unsetenv(test.ExpectedCommandEnv)
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
