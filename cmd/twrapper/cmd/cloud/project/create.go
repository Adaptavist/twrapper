package project

import (
	"log"

	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/flags"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/utils"
	"github.com/adaptavist/terraform-wrapper/v1/pkg/terraform/tfe-helper"
	"github.com/spf13/cobra"
)

func Create(c *config.Config) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "create",
		Short: "create terraform project",
		Run: func(cmd *cobra.Command, args []string) {
			// Input
			org := flags.GetOrganization(c, cmd)
			projectName := flags.GetName(c, cmd)
			teamName := flags.GetTeam(c, cmd)

			// Validation
			utils.FatalIfEmpty(org, "organanisation is required")
			utils.FatalIfEmpty(projectName, "project name is required")

			// Get client
			helper, err := tfe.New()
			utils.DieOnError(err)

			// Create project or fetch existing one
			project, err := helper.CreateProjectIfNotExists(org, projectName)
			utils.DieOnError(err)

			// Provide some feedback
			log.Printf("Project %s (%s) has been created", project.Name, project.ID)

			// Attach a team if provided by the user
			if teamName != "" {
				team, err := helper.GetTeamByName(org, teamName)
				utils.DieOnError(err)
				_, err = helper.AddTeamProjectAccessIfNotExist(project.ID, team.ID)
				utils.DieOnError(err)
			}

			// Provide some feedback
			projectAccess, err := helper.GetProjectTeams(project.ID)
			utils.DieOnError(err)

			log.Println("Teams that have access to the project:")
			for _, projectAccess := range projectAccess {
				log.Printf("%s - %s\n", projectAccess.ID, projectAccess.Name)
			}
		},
	}
	cmd.Flags().String("name", "", "Terraform Cloud Project Name")
	cmd.Flags().String("team", "", "Terraform Cloud Team that owns the project")
	return
}
