package backend

import (
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/spf13/cobra"
)

func Root(c *config.Config) *cobra.Command {
	root := &cobra.Command{
		Use:              "backend",
		Short:            "Terraform backend commands",
		TraverseChildren: true,
	}
	root.AddCommand(
		Show(c),
		List(c),
	)
	return root
}
