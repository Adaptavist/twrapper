package config

import (
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/spf13/cobra"
)

func Root(c *config.Config) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:              "config",
		Short:            "Twrapper configuration commands",
		TraverseChildren: true, // Needed to pass global flags to the command
	}
	cmd.AddCommand(
		Check(c),
	)
	return
}
