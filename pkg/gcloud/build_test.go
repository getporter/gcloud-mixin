package gcloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMixin_Build(t *testing.T) {
	m := NewTestMixin(t)

	err := m.Build()
	require.NoError(t, err, "Build failed")

	gotOutput := m.TestContext.GetOutput()

	wantOutput := `RUN apt-get update && apt-get install -y apt-transport-https ca-certificates curl gnupg
RUN echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" > /etc/apt/sources.list.d/google-cloud-sdk.list
RUN curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add -
RUN apt-get update && apt-get install -y google-cloud-sdk
`

	assert.Equal(t, wantOutput, gotOutput)
}
