package cmd

import (
	"errors"
	"fmt"
	"github.com/adaptavist/terraform-wrapper/v1/pkg/aws"
	"github.com/adaptavist/terraform-wrapper/v1/pkg/terraform"
	"log"
	"os"
	"strings"
)

// Config model
type Config struct {
	RequiredVars       []string               `mapstructure:"required_vars" yaml:"required_vars" json:"required_vars"`                // RequiredVars specifies env vars that must be set for terraform
	BackendType        string                 `mapstructure:"backend_type" yaml:"backend_type" json:"backend_type"`                   // What is the Terraform backend type (s3)?
	BackendConfig      map[string]interface{} `mapstructure:"backend_config" yaml:"backend_config" json:"backend_config"`             // BackendConfig
	BackendKeyIDVar    string                 `mapstructure:"backend_key_id_var" yaml:"backend_key_id_var" json:"backend_key_id_var"` // Variable we'll find the key_id from
	BackendKeyTemplate string                 `mapstructure:"backend_key_template" yaml:"backend_key_template" json:"backend_key_template"`
	AWS                aws.ConfigAWS          `mapstructure:"aws" yaml:"aws" json:"aws"`
}

func (c Config) getBackendKey() (key string, err error) {
	if c.BackendKeyIDVar == "" {
		err = errors.New("backend_key_id_var is not set")
		return
	}

	keyID, ok := os.LookupEnv(c.BackendKeyIDVar)
	if !ok {
		err = fmt.Errorf("%s is not set or empty", c.BackendKeyIDVar)
		return
	}

	if !IsValidUUID(keyID) {
		err = fmt.Errorf("%s is not a valid UUID v4", keyID)
	}

	// enforce lower case
	key = strings.ToLower(keyID)

	if c.BackendKeyTemplate != "" {
		key = strings.Replace(c.BackendKeyTemplate, "{key_id}", key, -1)

		if strings.Count(key, keyID) != 1 {
			err = errors.New("key_id appears more than one in backend_key")
		}
	} else {
		key = fmt.Sprintf("%s.tfstate", key)
	}

	return
}

func (c Config) Backend() (backend terraform.Backend) {
	backend = terraform.NewBackend(c.BackendType, c.BackendConfig)
	backendKey, err := c.getBackendKey()
	fatalIfNotNil(err, "%s")
	for k := range backend {
		switch k {
		case "s3":
			backend[k]["key"] = backendKey
		case "local":
			backend[k]["path"] = backendKey
		default:
			log.Fatal("unknown backend")
		}
	}
	return
}

// checkRequirements checks if required variables are missing
func (c Config) checkRequirements() (err error) {
	hasErrors := false
	for _, k := range c.RequiredVars {
		v, ok := os.LookupEnv(k)
		if ok {
			if strings.TrimSpace(v) == "" {
				hasErrors = true
				fmt.Printf("%s is empty\n", k)
			}
		} else {
			hasErrors = true
			fmt.Printf("%s is missing\n", k)
		}
	}

	if hasErrors {
		err = errors.New("fix above errors and try again")
	}
	return
}

// ActionArgs applies some extra args based on the Terraform verb
func ActionArgs(args []string) []string {
	switch args[0] {
	case "apply":
		return append(args, "-auto-approve")
	case "destroy":
		return append(args, "-force")
	}
	return args
}
