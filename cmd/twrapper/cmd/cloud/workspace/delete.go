package workspace

import (
	"log"

	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/flags"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/utils"
	cloud "github.com/adaptavist/terraform-wrapper/v1/pkg/terraform/tfe-helper"
	"github.com/spf13/cobra"
)

func Delete(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete terraform workspace",
		Run: func(cmd *cobra.Command, args []string) {
			// Get input
			org := flags.GetOrganization(c, cmd)
			workspace := flags.GetWorkspace(c, cmd)
			approve := flags.GetApprove(c, cmd)

			// Validate input
			utils.FatalIfEmpty(org, "organization is required")
			utils.FatalIfEmpty(workspace, "workspace is required")
			utils.FatalIfFalse(approve, "approval required when deleting a workspace")

			// Get client
			cloudHelper, err := cloud.New()
			utils.DieOnError(err)

			// Do some work
			err = cloudHelper.DeleteWorkspace(org, workspace)
			utils.DieOnError(err)

			// Give some feedback
			log.Printf("Deleted %s/%s\n", org, workspace)
		},
	}
	cmd.Flags().String("workspace", "", "Terraform Cloud Workspace Name")
	cmd.Flags().Bool("approve", false, "Explicitly approve deletion")
	return cmd
}
