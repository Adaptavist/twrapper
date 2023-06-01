package tfe

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/adaptavist/terraform-wrapper/v1/pkg/terraform/model/credentials"
	"github.com/hashicorp/go-tfe"
)

func makeTerraformClient(token string) (*tfe.Client, error) {
	// Create config
	tfConfig := &tfe.Config{
		Token: token,
	}
	// Return client if it initializes
	if client, err := tfe.NewClient(tfConfig); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func getTerraformEnterpriseClient() (*tfe.Client, error) {
	// Environment variable takes precedence
	token := os.Getenv("TF_TOKEN_app_terraform_io")
	if token != "" {
		return makeTerraformClient(token)
	}

	// Try the local terraform credentials file
	creds := credentials.New()
	if cred, ok := creds.Get("app.terraform.io"); ok {
		return makeTerraformClient(cred.Token)
	}

	return nil, fmt.Errorf("terraform token not found")
}

type ClientHelper struct {
	Context context.Context
	Client  *tfe.Client
}

func NewClientHelper(ctx context.Context, client *tfe.Client) *ClientHelper {
	return &ClientHelper{
		Context: ctx,
		Client:  client,
	}
}

func New() (*ClientHelper, error) {
	client, err := getTerraformEnterpriseClient()
	if err != nil {
		return nil, err
	}

	return NewClientHelper(context.Background(), client), nil
}

func (h *ClientHelper) GetProjectByName(org, name string) (*tfe.Project, error) {
	options := tfe.ProjectListOptions{
		Name: name,
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
	}

	projects, err := h.Client.Projects.List(h.Context, org, &options)

	if err != nil {
		return nil, err
	}

	for options.PageNumber <= projects.TotalPages {
		for _, p := range projects.Items {
			if p.Name == name {
				return p, nil
			}
		}
		if options.PageNumber == projects.TotalPages {
			break
		}
		options.ListOptions.PageNumber = projects.NextPage
	}

	return nil, errors.New("project not found")
}

func (h *ClientHelper) CreateProject(org, name string) (*tfe.Project, error) {
	return h.Client.Projects.Create(h.Context, org, tfe.ProjectCreateOptions{
		Name: name,
	})
}

func (h *ClientHelper) CreateProjectIfNotExists(org, name string) (*tfe.Project, error) {
	project, err := h.GetProjectByName(org, name)

	if err != nil {
		if err.Error() == "project not found" || err.Error() == "Project not found, or you are not authorized to use it." {
			return h.CreateProject(org, name)
		}
		return nil, err
	}

	return project, nil
}

func (h *ClientHelper) GetProjectTeams(projectID string) ([]*tfe.Team, error) {
	teams := make([]*tfe.Team, 0)
	access, err := h.Client.TeamProjectAccess.List(h.Context, tfe.TeamProjectAccessListOptions{
		ProjectID: projectID,
	})

	if err != nil {
		return nil, err
	}

	for _, t := range access.Items {
		team, err := h.Client.Teams.Read(h.Context, t.Team.ID)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func (h *ClientHelper) AddTeamProjectAccess(projectID, teamID string) (*tfe.TeamProjectAccess, error) {
	return h.Client.TeamProjectAccess.Add(h.Context, tfe.TeamProjectAccessAddOptions{
		Team:    &tfe.Team{ID: teamID},
		Project: &tfe.Project{ID: projectID},
		Access:  tfe.TeamProjectAccessMaintain,
	})
}

func (h *ClientHelper) AddTeamProjectAccessIfNotExist(projectID, teamID string) (*tfe.TeamProjectAccess, error) {
	access, err := h.Client.TeamProjectAccess.List(h.Context, tfe.TeamProjectAccessListOptions{
		ProjectID: projectID,
	})

	if err != nil {
		return nil, err
	}

	for _, a := range access.Items {
		if a.Team.ID == teamID {
			return a, nil
		}
	}

	return h.AddTeamProjectAccess(projectID, teamID)
}

func (h *ClientHelper) GetTeamByName(org, team string) (*tfe.Team, error) {
	teams, err := h.Client.Teams.List(h.Context, org, &tfe.TeamListOptions{
		Names: []string{team},
	})

	if err != nil {
		return nil, err
	}

	for _, t := range teams.Items {
		if t.Name == team {
			return t, nil
		}
	}

	return nil, errors.New("No team found with name: " + team)
}

func (h *ClientHelper) GetWorkspace(org, name string) (*tfe.Workspace, error) {
	return h.Client.Workspaces.Read(h.Context, org, name)
}

func (h *ClientHelper) CreateWorkspaceIfNotExists(org, name, projectID string) (*tfe.Workspace, error) {
	workspaces, err := h.Client.Workspaces.Read(h.Context, org, name)

	if err != nil {
		// If the error received is not "not found" return it, otherwise, we'll continue
		if err.Error() != "resource not found" {
			return nil, err
		}

		execMode := "local"
		options := tfe.WorkspaceCreateOptions{
			Name:          &name,
			ExecutionMode: &execMode,
		}

		if projectID != "" {
			options.Project = &tfe.Project{ID: projectID}
		}

		workspace, err := h.Client.Workspaces.Create(h.Context, org, options)

		if err != nil {
			return nil, fmt.Errorf("failed to create workspace: %v", err)
		}

		return workspace, nil
	}

	return workspaces, nil
}

func (h *ClientHelper) CreateWorkspaceAndProjectIfNotExists(org, project, workspace *string) (w *tfe.Workspace, err error) {
	var p *tfe.Project
	// Create a project if it doesn't exist when a project is supplied
	if project != nil {
		if p, err = h.CreateProjectIfNotExists(*org, *project); err != nil {
			return nil, fmt.Errorf("failed to find or create project: %v", err)
		}
	}
	// Create the workspace
	w, err = h.CreateWorkspaceIfNotExists(*org, *workspace, p.ID)
	return w, err
}

func (h *ClientHelper) DeleteProject(org, name string) error {
	project, err := h.GetProjectByName(org, name)

	if err != nil {
		return err
	}

	return h.Client.Projects.Delete(h.Context, project.ID)
}

func (h *ClientHelper) DeleteWorkspace(org, workspace string) error {
	return h.Client.Workspaces.Delete(h.Context, org, workspace)
}
