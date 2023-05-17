package project

import (
	"log"

	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/flags"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/utils"
	"github.com/adaptavist/terraform-wrapper/v1/pkg/terraform/tfe-helper"
	"github.com/spf13/cobra"
)

func Delete(c *config.Config) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "delete",
		Short: "delete terraform project",
		Run: func(cmd *cobra.Command, args []string) {
			// Input
			org := flags.GetOrganization(c, cmd)
			projectName := flags.GetName(c, cmd)
			approve := flags.GetApprove(c, cmd)

			// Validation
			utils.FatalIfEmpty(org, "organanisation is required")
			utils.FatalIfEmpty(projectName, "project name is required")
			utils.FatalIfFalse(approve, "Please approve the deletion")

			// Get client
			helper, err := tfe.New()
			utils.DieOnError(err)

			// Create project or fetch existing one
			err = helper.DeleteProject(org, projectName)
			utils.DieOnError(err)

			// Provide some feedback
			log.Printf("Project %s has been deleted", projectName)
		},
	}
	cmd.Flags().String("name", "", "Terraform Cloud Project Name")
	cmd.Flags().Bool("approve", false, "Approve the deletion")
	return
}
