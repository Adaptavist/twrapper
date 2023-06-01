package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/cmd/root"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/adaptavist/terraform-wrapper/v1/pkg/integration/test/helpers"
	"github.com/adaptavist/terraform-wrapper/v1/pkg/integration/test/helpers/fixture"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// cleanup temporary files after tests
func cleanup(dir string) {
	fmt.Printf("cleaning up temporary files in %s\n", dir)
	if err := os.RemoveAll(dir); err != nil {
		panic(err)
	}
}

// TestTerraformCommands runs through a complete terraform lifecycle
func TestTerraformCommands(t *testing.T) {
	dir := t.TempDir()
	defer cleanup(dir)

	// Twrapper configuration
	key := uuid.NewString()
	os.Setenv("TW_BACKEND_KEY", key)
	os.Setenv("TW_WORKSPACE_ID", strings.Split(key, "-")[0])

	// Get config
	conf := helpers.Config()
	conf.Backend = &config.Backend{
		Type: "local",
		Props: config.ConfigMap{
			"path": "terraform.tfstate",
		},
	}

	// Create terraform module
	fixture.Write("random/terraform.tfstate", fmt.Sprintf("%s/%s", dir, conf.Backend.Props["path"]))
	fixture.Write("random/main.tf", fmt.Sprintf("%s/main.tf", dir))

	// We're going to run a few commands, so lets wrap up nicely.
	f := func(args ...string) {
		cmd := root.Root(conf)
		cmd.SetArgs(args)
		err := cmd.Execute()
		assert.Nil(t, err)
	}

	// Run some commands
	f("init", "--chdir", dir)
	f("plan", "--chdir", dir)
	f("apply", "--chdir", dir)
	f("destroy", "--chdir", dir)
}
