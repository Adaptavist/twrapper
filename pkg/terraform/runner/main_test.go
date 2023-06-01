package runner

import (
	"testing"

	"github.com/adaptavist/terraform-wrapper/v1/pkg/terraform/model/hcl"
	"github.com/stretchr/testify/assert"
)

func TestBlankCloudConfig(t *testing.T) {
	terraform := New()
	terraform.WithBackend(hcl.NewBackend("s3", map[string]interface{}{
		"bucket":         "my-bucket",
		"dynamodb_table": "my-table",
		"key":            "my-key",
		"region":         "us-east-1",
		"encrypted":      true,
	}))

	json, err := terraform.Config.ToJSON()
	assert.NoError(t, err)
	assert.Equal(t, string(json), `{
	"terraform": {
		"backend": {
			"s3": {
				"bucket": "my-bucket",
				"dynamodb_table": "my-table",
				"encrypted": true,
				"key": "my-key",
				"region": "us-east-1"
			}
		}
	}
}`)
}
