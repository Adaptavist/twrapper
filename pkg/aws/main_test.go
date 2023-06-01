package aws

import (
	"testing"

	"github.com/adaptavist/terraform-wrapper/v1/pkg/terraform/runner"
	"github.com/stretchr/testify/assert"
)

func TestFixedRoleIsConfigured(t *testing.T) {
	roleARN := "arn:aws:aim:::role/RoleName"
	awsConfig := ConfigAWS{
		RoleARN:   roleARN,
		RoleTFVar: "role_arn",
	}
	terraform := runner.New()
	err := awsConfig.Configure(&terraform)
	assert.Nil(t, err, "configured should return nil error")
	assert.Equal(t, roleARN, terraform.Variables["role_arn"], "RoleARN should have made it to terraform as a var")
}
