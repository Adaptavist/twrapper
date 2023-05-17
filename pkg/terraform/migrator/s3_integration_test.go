//go:build integration_test
// +build integration_test

package migrator

import (
	"os"
	"strings"
	"testing"

	"github.com/adaptavist/terraform-wrapper/v1/pkg/integration/test/helpers"
	"github.com/adaptavist/terraform-wrapper/v1/pkg/integration/test/helpers/fixture"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestS3Migration(t *testing.T) {
	key := uuid.NewString()

	// Twrapper configuration
	os.Setenv("TW_BACKEND_KEY", key)
	os.Setenv("TW_WORKSPACE_ID", strings.Split(key, "-")[0])

	conf := helpers.Config()

	t.Logf("backend key: %s", conf.Backend.Props["key"])

	// Writing the state file so we can migrate it.
	state := fixture.Get("random/terraform.tfstate")
	err := helpers.WriteS3BackendObject(conf, state)
	assert.NoError(t, err)

	// Create the migrator
	mig, err := New(conf)

	if err != nil {
		t.Errorf("failed to create migrator: %s", err)
	}

	if err = mig.Init(); err != nil {
		t.Errorf("Failed to init migrate: %v", err)
	}

	//Cleanup
	defer mig.Helper().DeleteWorkspace(
		mig.Config().Cloud.Organization,
		mig.Config().Cloud.Workspace)

	// Migrate
	if err = mig.Migrate(); err != nil {
		t.Errorf("Failed to migrate: %v", err)
	}

	if err = mig.Cleanup(); err != nil {
		t.Errorf("Failed to cleanup: %v", err)
	}
}
