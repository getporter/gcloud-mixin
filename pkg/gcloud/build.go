package gcloud

import "fmt"

// Build will generate the necessary Dockerfile lines
// for an invocation image using this mixin
func (m *Mixin) Build() error {
	fmt.Fprintln(m.Out, `RUN apt-get update && apt-get install -y apt-transport-https ca-certificates curl gnupg`)
	fmt.Fprintln(m.Out, `RUN echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" > /etc/apt/sources.list.d/google-cloud-sdk.list`)
	fmt.Fprintln(m.Out, `RUN curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add -`)
	fmt.Fprintln(m.Out, `RUN apt-get update && apt-get install -y google-cloud-sdk`)

	return nil
}
