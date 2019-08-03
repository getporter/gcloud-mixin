package gcloud

import (
	"fmt"

	"github.com/deislabs/porter-gcloud/pkg"
)

func (m *Mixin) PrintVersion() {
	fmt.Fprintf(m.Out, "gcloud mixin %s (%s)\n", pkg.Version, pkg.Commit)
}
