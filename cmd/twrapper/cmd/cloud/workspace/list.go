package workspace

import (
	"context"
	"fmt"

	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/flags"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/utils"
	cloud "github.com/adaptavist/terraform-wrapper/v1/pkg/terraform/tfe-helper"
	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func List(c *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list terraform workspace",
		Run: func(cmd *cobra.Command, args []string) {
			// get our helper object
			cloudHelper, err := cloud.New()
			utils.DieOnError(err)
			// get tehe client
			tfeClient := cloudHelper.Client

			// Initial list options, will mutate as we paginate or set flags
			workspaceListOptions := tfe.WorkspaceListOptions{
				ListOptions: tfe.ListOptions{
					PageSize:   100,
					PageNumber: 1,
				},
			}

			org := flags.GetOrganization(c, cmd)
			c.Logger.Debug("filter", zap.String("org_name", org))

			// if project flag is set, lookup its id and add it to the filter.
			if p := flags.GetProject(c, cmd); p != "" {
				project, err := cloudHelper.GetProjectByName(org, p)
				utils.DieOnError(err)
				workspaceListOptions.ProjectID = project.ID
				c.Logger.Debug("filter", zap.String("project_name", p), zap.String("project_id", project.ID))
			}

			// do initial request
			workspaceList, err := tfeClient.Workspaces.List(context.Background(), org, &workspaceListOptions)
			utils.DieOnError(err)

			c.Logger.Debug("search",
				zap.Int("page_count", workspaceList.TotalPages),
				zap.Int("page_numb", workspaceList.CurrentPage),
			)

			// iterate through pages
			for workspaceList.CurrentPage <= (workspaceList.TotalPages) {
				for _, workspace := range workspaceList.Items {
					version := "No Version"
					if workspace.CurrentStateVersion != nil {
						version = workspace.CurrentStateVersion.ID
					}
					fmt.Println(workspace.ID, workspace.Name, workspace.Project.ID, workspace.Project.Name, version)
				}

				if workspaceList.Pagination.NextPage == 0 {
					break
				}

				workspaceListOptions.ListOptions.PageNumber = workspaceList.Pagination.NextPage
				workspaceList, err = tfeClient.Workspaces.List(context.Background(), org, &workspaceListOptions)
				utils.DieOnError(err)
			}
		},
	}
}
