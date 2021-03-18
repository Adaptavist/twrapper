package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequiredVarsFail(t *testing.T) {
	c := Config{
		RequiredVars: []string{"HOPEFULLY_RANDOM_ENV_VAR"},
	}
	err := c.checkRequirements()
	assert.NotNil(t, err, "error should be nil")
}

func TestRequiredVarsFailOnEmptyVar(t *testing.T) {
	if err := os.Setenv("TEST_KEY", ""); err != nil {
		panic(err)
	}
	c := Config{
		RequiredVars: []string{"TEST_KEY"},
	}
	err := c.checkRequirements()
	assert.NotNil(t, err, "error should be nil")
	os.Clearenv()
}

func TestRequiredVars(t *testing.T) {
	if err := os.Setenv("TEST_KEY", "VALUE"); err != nil {
		panic(err)
	}
	c := Config{
		RequiredVars: []string{"TEST_KEY"},
	}
	err := c.checkRequirements()
	assert.Nil(t, err, "error should not be nil")
	os.Clearenv()
}
