//go:generate packr2
package gcloud

import (
	"github.com/deislabs/porter/pkg/context"
)

type Mixin struct {
	*context.Context
}

// New gcloud mixin client, initialized with useful defaults.
func New() (*Mixin, error) {
	return &Mixin{
		Context: context.New(),
	}, nil

}
