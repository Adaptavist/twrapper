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

func Read(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "read terraform workspace",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.TODO()
			cloudHelper, err := cloud.New()
			utils.DieOnError(err)

			tfeClient := cloudHelper.Client
			orgName := flags.GetOrganization(c, cmd)

			c.Logger.Debug("search", zap.String("org_name", orgName), zap.Any("workspace_names", args))

			showVersions, err := cmd.Flags().GetBool("state-versions")

			if err != nil {
				c.Logger.Error("show versions", zap.Error(err))
			}

			for _, v := range args {
				workspace, err := tfeClient.Workspaces.Read(ctx, orgName, v)
				if err != nil {
					// Move onto the next arg if its simply a resource error
					if err.Error() == "resource not found" {
						fmt.Printf("%s: %s\n", err.Error(), v)
						continue
					}
					// Something else must be a miss, lets bail
					utils.DieOnError(err)
				}

				fmt.Println(workspace.ID, workspace.Name)

				if showVersions {
					opts := tfe.StateVersionListOptions{
						Organization: orgName,
						Workspace:    v,
					}

					for res, err := tfeClient.StateVersions.List(ctx, &opts); true; {
						utils.DieOnError(err)

						for _, version := range res.Items {
							fmt.Println(version.ID)
						}

						if res.NextPage == 0 {
							break
						}
					}
				}
			}
		},
	}
	cmd.Flags().Bool("state-versions", false, "show state versions")
	return cmd
}
