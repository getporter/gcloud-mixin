package main

import (
	"github.com/deislabs/porter-gcloud/pkg/gcloud"
	"github.com/spf13/cobra"
)

func buildVersionCommand(m *gcloud.Mixin) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the mixin version",
		Run: func(cmd *cobra.Command, args []string) {
			m.PrintVersion()
		},
	}
}
