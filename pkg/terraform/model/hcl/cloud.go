package hcl

import (
	"errors"
)

// Flatterns terraform cloud configuration structure
type Cloud struct {
	Organization string `json:"organization"`
	Workspace    string `json:"workspace"`
	Project      string `json:"project"`
}

func (c *Cloud) Validate() error {
	if len(c.Organization) == 0 {
		return errors.New("organizations required")
	}

	if len(c.Workspace) == 0 {
		return errors.New("workspace required")
	}

	return nil
}

func NewCloud(organization string, workspace string) *Cloud {
	return &Cloud{
		Organization: organization,
		Workspace:    workspace,
	}
}

func NewCloudWithProject(organization string, workspace string, project string) *Cloud {
	return &Cloud{
		Organization: organization,
		Workspace:    workspace,
		Project:      project,
	}
}

func (c *Cloud) SetWorkspace(workspace string) {
	c.Workspace = workspace
}

func (c *Cloud) GetWorkspace() string {
	return c.Workspace
}

func (c *Cloud) IsEmpty() bool {
	return len(c.Organization) == 0 && len(c.Workspace) == 0
}

func (c *Cloud) ToHCL() (interface{}, error) {
	return map[string]interface{}{
		"organization": c.Organization,
		"workspaces": map[string]interface{}{
			"name": c.Workspace,
		},
	}, nil
}
