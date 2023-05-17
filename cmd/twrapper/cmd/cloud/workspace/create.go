package workspace

import (
	"log"

	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/flags"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/utils"
	cloud "github.com/adaptavist/terraform-wrapper/v1/pkg/terraform/tfe-helper"
	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

func Create(c *config.Config) *cobra.Command {
	create := false
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create terraform workspace",
		Run: func(cmd *cobra.Command, args []string) {
			// User input
			orgName := flags.GetOrganization(c, cmd)
			workspaceName := flags.GetWorkspace(c, cmd)
			projectName := flags.GetProject(c, cmd)
			// Make sure certain things are properly set
			utils.FatalIfEmpty(workspaceName, "workspace name is required")
			utils.FatalIfEmpty(orgName, "organization name is required")

			// get our tfe helper
			cloudHelper, err := cloud.New()
			utils.DieOnError(err)

			// create workspace within a project if a project is specified
			if projectName != "" {
				var project *tfe.Project

				// only auto create if create flag is true
				if create {
					project, err = cloudHelper.CreateProjectIfNotExists(orgName, projectName)
					utils.DieOnError(err)
				} else {
					project, err = cloudHelper.GetProjectByName(orgName, projectName)
					utils.DieOnError(err)
				}

				// create workspace linked to project
				workspace, err := cloudHelper.CreateWorkspaceIfNotExists(orgName, workspaceName, project.ID)
				utils.DieOnError(err)
				// feedback
				log.Printf("Created %s/%s/%s\n", orgName, project.Name, workspace.Name)
				// no need for the function to continue
				return
			}

			workspace, err := cloudHelper.CreateWorkspaceIfNotExists(orgName, workspaceName, "")
			utils.DieOnError(err)
			// feedback
			log.Printf("Created %s/%s\n", orgName, workspace.Name)
		},
	}
	cmd.Flags().String("workspace", "", "Terraform Cloud Workspace Name")
	cmd.Flags().BoolVarP(&create, "create", "c", false, "Create project if it doesn't exist")
	return cmd
}
