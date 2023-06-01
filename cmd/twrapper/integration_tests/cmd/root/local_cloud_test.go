//go:build integration_test
// +build integration_test

package root

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/cmd/cloud/workspace"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/cmd/root"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/adaptavist/terraform-wrapper/v1/pkg/integration/test/helpers"
	"github.com/adaptavist/terraform-wrapper/v1/pkg/integration/test/helpers/fixture"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTerraformLocalCloudMiration(t *testing.T) {
	dir := t.TempDir()
	defer helpers.CleanupDIR(dir)

	// Setup some vars
	key := uuid.NewString()

	os.Setenv("TW_BACKEND_KEY", key)                         // Not actually in this test, but the config contains it
	os.Setenv("TW_WORKSPACE_ID", strings.Split(key, "-")[0]) // Used to create the workspace

	workspaceKey := path.Base(path.Dir(dir))

	// bootstrap our testing environment
	conf := helpers.Config()
	conf.Cloud.Workspace = workspaceKey
	conf.Backend.Type = "local"
	conf.Backend.Props = config.ConfigMap{
		"path": "terraform.tfstate",
	}

	// Setup a working directory
	fixture.Write("random/main.tf", fmt.Sprintf("%s/%s", dir, "main.tf"))
	fixture.Write("random/terraform.tfstate", fmt.Sprintf("%s/%s", dir, conf.Backend.Props["path"]))

	// we're going to run a few commands, so lets wrap up nicely.
	f := func(args ...string) {
		cmd := root.Root(conf)
		cmd.SetArgs(args)
		err := cmd.Execute()
		assert.Nil(t, err)
	}

	// run some commands
	f("init", "--chdir", dir)
	f("plan", "--chdir", dir)
	f("apply", "--chdir", dir)
	f("destroy", "--chdir", dir)

	// cleanup
	cmd := workspace.Root(conf)
	cmd.SetArgs([]string{"delete", "--workspace", workspaceKey, "--approve"})
	err := cmd.Execute()
	assert.Nil(t, err)
}
