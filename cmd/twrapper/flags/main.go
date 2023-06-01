package flags

import (
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/spf13/cobra"
)

func getFlagString(name string, cmd *cobra.Command) string {
	o, err := cmd.Flags().GetString(name)
	if err == nil && o != "" {
		return o
	}
	return ""
}

func getFlagBool(name string, cmd *cobra.Command) bool {
	if o, err := cmd.Flags().GetBool(name); err == nil {
		return o
	}
	return false
}

func GetApprove(c *config.Config, cmd *cobra.Command) bool {
	return getFlagBool("approve", cmd)
}

func GetForce(c *config.Config, cmd *cobra.Command) bool {
	return getFlagBool("force", cmd)
}

func GetTeam(c *config.Config, cmd *cobra.Command) string {
	return getFlagString("team", cmd)
}

func GetName(c *config.Config, cmd *cobra.Command) string {
	return getFlagString("name", cmd)
}

func GetOrganization(c *config.Config, cmd *cobra.Command) string {
	if o := getFlagString("org", cmd); o != "" {
		return o
	}
	if c.Cloud.Organization != "" {
		return c.Cloud.Organization
	}
	return ""
}

func GetProject(c *config.Config, cmd *cobra.Command) string {
	if o := getFlagString("project", cmd); o != "" {
		return o
	}
	if c.Cloud.Project != "" {
		return c.Cloud.Project
	}
	return ""
}

func GetWorkspace(c *config.Config, cmd *cobra.Command) string {
	if o := getFlagString("workspace", cmd); o != "" {
		return o
	}
	if c.Cloud.Organization != "" {
		return c.Cloud.Workspace
	}
	return ""
}
