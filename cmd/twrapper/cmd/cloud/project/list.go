package project

import (
	"context"
	"fmt"

	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/flags"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/utils"
	cloud "github.com/adaptavist/terraform-wrapper/v1/pkg/terraform/tfe-helper"
	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

func List(c *config.Config) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "list",
		Short: "list terraform projects",
		Run: func(cmd *cobra.Command, args []string) {
			// get our helper object
			cloudHelper, err := cloud.New()
			utils.DieOnError(err)

			// get tehe client
			client := cloudHelper.Client

			org := flags.GetOrganization(c, cmd)
			utils.FatalIfEmpty(org, "organisation required")

			// default options var
			listOptions := tfe.ListOptions{
				PageSize:   100,
				PageNumber: 1,
			}

			// do the initial query
			projectList, err := client.Projects.List(context.Background(), org, &tfe.ProjectListOptions{ListOptions: listOptions})
			utils.FatalIfNotNil(err, "Error listing projects")

			// paginate
			for projectList.CurrentPage <= projectList.TotalPages {
				for _, project := range projectList.Items {
					fmt.Println(project.ID, project.Name)
				}
				// stop looping once we reach the page limit
				if projectList.CurrentPage == projectList.TotalPages {
					break
				}
				// setup next page
				listOptions.PageNumber = projectList.NextPage
				projectList, err = client.Projects.List(context.Background(), org, &tfe.ProjectListOptions{ListOptions: listOptions})
				utils.FatalIfNotNil(err, "failed to list projects")
			}
		},
	}
	cmd.PersistentFlags().String("project", "", "return workspace matching project")
	return
}
