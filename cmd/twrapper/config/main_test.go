package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBackendConfigMap(t *testing.T) {
	c := ConfigMap{}
	c.Set("foo", "bar")
	val, err := c.GetString("foo")
	assert.NoError(t, err)
	assert.Equal(t, "bar", val)
}

func TestRequiredVarsFail(t *testing.T) {
	c := Config{
		RequiredVars: &[]string{"HOPEFULLY_RANDOM_ENV_VAR"},
	}
	err := c.CheckRequiredVars()
	assert.NotNil(t, err, "error should be nil")
}

func TestRequiredVarsFailOnEmptyVar(t *testing.T) {
	if err := os.Setenv("TEST_KEY", ""); err != nil {
		panic(err)
	}
	c := Config{
		RequiredVars: &[]string{"TEST_KEY"},
	}
	err := c.CheckRequiredVars()
	assert.NotNil(t, err, "error should be nil")
	os.Clearenv()
}

func TestRequiredVars(t *testing.T) {
	if err := os.Setenv("TEST_KEY", "VALUE"); err != nil {
		panic(err)
	}
	c := Config{
		RequiredVars: &[]string{"TEST_KEY"},
	}
	err := c.CheckRequiredVars()
	assert.Nil(t, err, "error should not be nil")
	os.Clearenv()
}
