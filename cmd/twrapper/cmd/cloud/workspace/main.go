package workspace

import (
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/spf13/cobra"
)

func Root(c *config.Config) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:              "workspace",
		Short:            "Terraform workspace commands",
		TraverseChildren: true, // Needed to pass global flags to the command
	}
	cmd.AddCommand(
		Create(c),
		Delete(c),
		List(c),
		Read(c),
	)
	cmd.Flags().String("project", "", "Terraform workspace project name")
	return
}
