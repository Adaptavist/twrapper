package project

import (
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/spf13/cobra"
)

func Root(c *config.Config) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:              "project",
		Short:            "Terraform project commands",
		TraverseChildren: true, // Needed to pass global flags to the command
	}
	cmd.AddCommand(
		List(c),
		Create(c),
		Delete(c),
	)
	return

}
