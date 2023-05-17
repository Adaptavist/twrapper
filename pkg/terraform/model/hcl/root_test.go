package hcl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootWithCloud(t *testing.T) {
	expect := `{
	"terraform": {
		"cloud": {
			"organization": "my-org",
			"workspaces": {
				"name": "my-workspace"
			}
		}
	}
}`

	config := Root{
		Terraform: Terraform{
			Cloud: NewCloud("my-org", "my-workspace"),
		},
	}

	json, err := config.ToJSON()
	assert.NoError(t, err)

	assert.Equal(t, expect, string(json))
}

func TestRootWithBackend(t *testing.T) {
	expect := `{
	"terraform": {
		"backend": {
			"local": {
				"path": "foo"
			}
		}
	}
}`
	config := Root{
		Terraform: Terraform{
			Backend: NewBackend("local", map[string]interface{}{
				"path": "foo",
			}),
		},
	}

	json, err := config.ToJSON()
	assert.NoError(t, err)
	assert.Equal(t, expect, string(json))
}
