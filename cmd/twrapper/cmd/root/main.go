package root

import (
	"fmt"
	"os"

	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/cmd/backend"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/cmd/cloud"
	cfg "github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/cmd/config"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/terraform"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/utils"
	"github.com/spf13/cobra"
)

// make creates a new cobra.Command for wrapping terraform
func make(c *config.Config, action string, tfArgs ...string) *cobra.Command {
	return &cobra.Command{
		Use:   action,
		Short: fmt.Sprintf("Wraps `terraform %s`", action),
		Run: func(cmd *cobra.Command, args []string) {
			// If the chdir flag is provided, we need to change the working directory
			if flag := cmd.Flag("chdir"); flag.Value.String() != "" {
				c.Logger.Debug("getting current working directory")
				pwd, err := os.Getwd()
				utils.DieOnError(err)

				c.Logger.Debug("changing working dir")
				utils.ChangeDirOrDie(flag.Value.String())

				defer utils.ChangeDirOrDie(pwd)
			}
			terraform.Run(c, append(tfArgs, args...)...)
		},
	}
}

func Root(c *config.Config) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:              "twrapper",
		Short:            "Wrapper for initialising Terraform",
		Long:             "Sets some things up before running Terraform in a CI environment.",
		TraverseChildren: true, // Needed to pass global flags to the command
	}
	cmd.PersistentFlags().String("chdir", "", "Change working dir")
	cmd.PersistentFlags().Bool("force-delete", false, "Force deletion of backend objects")
	cmd.AddCommand(
		make(c, "init", "init"),
		make(c, "plan", "plan"),
		make(c, "apply", "apply", "-auto-approve"),
		make(c, "destroy", "apply", "-destroy", "-auto-approve"),
		cfg.Root(c),
		backend.Root(c),
		cloud.Root(c),
	)
	return
}
