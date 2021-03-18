package aws

import (
	"github.com/adaptavist/terraform-wrapper/v1/pkg/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFixedRoleIsConfigured(t *testing.T) {
	roleARN := "arn:aws:aim:::role/RoleName"
	awsConfig := ConfigAWS{
		RoleARN:   roleARN,
		RoleTFVar: "role_arn",
	}
	opts := terraform.NewOpts()
	err := awsConfig.Configure(&opts)
	assert.Nil(t, err, "configured should return nil error")
	assert.Equal(t, roleARN, opts.Variables["role_arn"], "RoleARN should have made it to terraform as a var")
}
