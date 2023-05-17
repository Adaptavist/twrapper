package cloud

import (
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/cmd/cloud/project"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/cmd/cloud/workspace"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/spf13/cobra"
)

func Root(c *config.Config) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:              "cloud",
		Short:            "Terraform cloud commands",
		TraverseChildren: true, // Needed to pass global flags to the command
	}
	cmd.PersistentFlags().String("org", "", "Terraform Cloud Organization")
	cmd.AddCommand(
		project.Root(c),
		workspace.Root(c),
	)
	return
}
