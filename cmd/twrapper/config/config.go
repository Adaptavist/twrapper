package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/adaptavist/terraform-wrapper/v1/pkg/aws"
	"github.com/adaptavist/terraform-wrapper/v1/pkg/terraform/model/hcl"
	"go.uber.org/zap"
)

type Config struct {
	Logger       *zap.Logger
	RequiredVars *[]string      `mapstructure:"REQUIRED_VARS" yaml:"required_vars" json:"required_vars"` // RequiredVars specifies env vars that must be set for terraform
	Backend      *Backend       `mapstructure:"BACKEND" yaml:"backend" json:"backend"`                   // Backend configuration for Terraform
	Cloud        *hcl.Cloud     `mapstructure:"CLOUD" yaml:"cloud" json:"cloud"`                         // Cloud configurtion for Terraform
	AWS          *aws.ConfigAWS `mapstructure:"AWS" yaml:"aws" json:"aws"`                               // AWS configuration for Terraform
}

func (c *Config) CloudConfigIsSet() bool {
	return c.Cloud != nil
}

func (c *Config) BackendConfigIsSet() bool {
	return c.Backend != nil
}

func (c Config) MakeBackend() (backend *hcl.Backend) {
	if c.Backend == nil {
		return nil
	}

	return hcl.NewBackend(c.Backend.Type, c.Backend.Props)
}

// FindBackendKey returns the backend key from the CloudConfig, or BackendConfig
func (c Config) FindBackendKey() string {
	if c.CloudConfigIsSet() {
		if c.Cloud.Workspace != "" {
			return c.Cloud.Workspace
		}
	}

	if backend := c.MakeBackend(); backend.GetWorkspace() != "" {
		return backend.GetWorkspace()
	}

	return ""
}

// MakeCloud returns a completed cloud configuration
func (c Config) MakeCloud() *hcl.Cloud {
	return c.Cloud
}

// checkRequiredVars checks if required variables are missing
func (c Config) CheckRequiredVars() (errs []error) {
	if c.RequiredVars != nil {
		for _, k := range *c.RequiredVars {
			v, ok := os.LookupEnv(k)
			if ok {
				if strings.TrimSpace(v) == "" {
					errs = append(errs, fmt.Errorf("%s is empty", k))
				}
			} else {
				errs = append(errs, fmt.Errorf("%s is missing", k))
			}
		}
	}
	return
}

func Must(conf *Config, err error) *Config {
	if err != nil {
		panic(err)
	}
	return conf
}
