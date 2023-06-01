package terraform

import (
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/adaptavist/terraform-wrapper/v1/pkg/terraform/migrator"
	"github.com/adaptavist/terraform-wrapper/v1/pkg/terraform/runner"
	"github.com/adaptavist/terraform-wrapper/v1/pkg/terraform/tfe-helper"
	"go.uber.org/zap"
)

func Run(c *config.Config, args ...string) {
	c.Logger.Info("running terraform", zap.Any("args", args))

	// check if we are in a good state to be running terraform
	errs := c.CheckRequiredVars()

	if len(errs) > 0 {
		c.Logger.Error("missing fields", zap.Errors("errors", errs))
		return
	}

	// If we have a backend configuration and a cloud configuration, lets migrate to cloud.
	// In some cases (with global config) cloud may be set with an org but not workspace
	if c.Backend != nil && c.CloudConfigIsSet() {
		c.Logger.Info("mirating backend")
		err := migrator.Migrate(c)
		if err != nil {
			c.Logger.Fatal("failed to migrate backend to cloud", zap.Error(err))
		}
		c.Backend = nil
		c.Logger.Debug("new config", zap.Any("config", c))
	}

	// start configuration terraform workspace
	terraform := runner.New()
	terraform.WithArguments(args)

	if c.Backend != nil {
		terraform.WithBackend(c.MakeBackend())
	}

	if c.CloudConfigIsSet() {
		terraform.WithCloud(c.MakeCloud())
	}

	// At this point, we'll check if the terraform object has a cloud configuration
	if cloud := terraform.GetCloud(); cloud != nil {
		c.Logger.Info("checking Terraform Enterprise configuration")

		// Setup Terraform Cloud project and workspace
		client, err := tfe.New()
		if err != nil {
			c.Logger.Fatal("failed to create terraform enterprise client", zap.Error(err))
		}

		if err != nil {
			c.Logger.Fatal("failed to define workspace name", zap.Error(err))
		}

		// Create empty project id for when we create the workspace
		_, err = client.CreateWorkspaceAndProjectIfNotExists(&cloud.Organization, &cloud.Project, &cloud.Workspace)
		if err != nil {
			c.Logger.Fatal("failed to create create workspace and/or project", zap.Error(err))
		}
	}

	// Get terraform configured for AWS
	if c.AWS != nil && !c.AWS.IsEmpty() {
		c.Logger.Info("looking for AWS credentials")
		if err := c.AWS.Configure(&terraform); err != nil {
			c.Logger.Fatal("failed to configure aws session", zap.Error(err))
		}
	}

	c.Logger.Info("configuring terraform workspace")
	if err := terraform.Configure(); err != nil {
		c.Logger.Fatal("failed to configure terraform", zap.Error(err))
	}

	c.Logger.Info("running terraform")
	c.Logger.Info("-------------------------------------------------------------")

	if err := terraform.Go(); err != nil {
		c.Logger.Fatal("terraform operation failed", zap.Error(err))
	}
}
