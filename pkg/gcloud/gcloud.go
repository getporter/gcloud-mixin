package gcloud

import (
	"get.porter.sh/porter/pkg/runtime"
)

type Mixin struct {
	runtime.RuntimeConfig
}

// New gcloud mixin client, initialized with useful defaults.
func New() *Mixin {
	return &Mixin{
		RuntimeConfig: runtime.NewConfig(),
	}
}
