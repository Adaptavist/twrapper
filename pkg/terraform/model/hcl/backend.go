package hcl

import (
	"fmt"
)

// Backend is the main container for backend configuration of Terraform, it supports only local and s3 at present.
type Backend struct {
	Type   string                 `yaml:"type" json:"type"`
	Config map[string]interface{} `yaml:"config" json:"config"`
}

func (b *Backend) IsType(t string) bool {
	return b.Type == t
}

// Get backend configuration property
func (b *Backend) Get(prop string) interface{} {
	if p, ok := b.Config[prop]; ok {
		return p
	}
	return nil
}

func (b *Backend) GetString(prop string) string {
	if p, ok := b.Config[prop]; ok {
		return p.(string)
	}
	return ""
}

// Set backend configuration property
func (b *Backend) Set(prop string, value interface{}) {
	b.Config[prop] = value
}

// GetWorkspace gets the workspace for the given backend
func (b *Backend) GetWorkspace() string {
	switch b.Type {
	case "s3":
		if v, ok := b.Config["key"]; ok {
			return v.(string)
		}
	case "local":
		if v, ok := b.Config["path"]; ok {
			return v.(string)
		}
	}
	return ""
}

// SetWorkspace sets the workspace for the given backend
func (b *Backend) SetWorkspace(workspace string) {
	switch b.Type {
	case "s3":
		b.Config["key"] = workspace
	case "local":
		b.Config["path"] = workspace
	}
}

func (b *Backend) Validate() error {
	switch b.Type {
	case "local":
		if b.Get("path").(string) == "" {
			return fmt.Errorf("required field 'path' not found")
		}
	case "s3":
		if b.Get("Bucket").(string) == "" {
			return fmt.Errorf("required field 'bucket' not found")
		}

		if b.Get("Region").(string) == "" {
			return fmt.Errorf("required field 'region' not found")
		}

		if b.Get("DynamoDBTable").(string) == "" {
			return fmt.Errorf("required field 'table' not found")
		}

		if b.Get("Key").(string) == "" {
			return fmt.Errorf("required field 'key' not found")
		}

		if v := b.Get("encrpyted"); v != nil {
			return fmt.Errorf("required field 'encrpyted' not found")
		}
	default:
		return fmt.Errorf("unknown backend type: %s", b.Type)
	}
	return nil
}

func (b *Backend) ToHCL() (interface{}, error) {
	return map[string]interface{}{
		b.Type: b.Config,
	}, nil
}

func NewBackend(backendType string, backendConfig map[string]interface{}) *Backend {
	return &Backend{
		Type:   backendType,
		Config: backendConfig,
	}
}
