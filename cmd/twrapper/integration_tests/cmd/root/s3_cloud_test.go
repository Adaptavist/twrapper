//go:build integration_test
// +build integration_test

package root

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/cmd/root"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/utils"
	"github.com/adaptavist/terraform-wrapper/v1/pkg/integration/test/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var fixturePath string

// getFixturePath returns the path to the fixture directory
// because we change our cwd during testing, we set once and reuse to not lose the relative path
func getFixturePath() string {
	if fixturePath == "" {
		path, err := filepath.Abs("../../../../../test/fixtures/basic/main.tf")
		utils.DieOnError(err)
		fixturePath = path
	}
	return fixturePath
}

// copyFile copies source file to destination file
// stolen from https://github.com/mactsouk/opensource.com/blob/master/cp1.go
func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// localWorkspaceInit initialises a terraform workspace locally
func localWorkspaceInit(dir string) {
	_, err := copyFile(getFixturePath(), fmt.Sprintf("%s/main.tf", dir))
	utils.DieOnError(err)
}

// localWorkspaceClean deletes the local workspace directory
func localWorkspaceClean(dir string) {
	os.RemoveAll(dir)
}

// testBackendInit initialises the workspace with just the backend
// t *testing.T is used for assertions
// c *config.Config is the configuration to test with
// dir is the directory we're working in
func testBackendInit(t *testing.T, c *config.Config, dir string) {
	fmt.Println("PREPARE BACKEND FOR MIGRATION")
	fmt.Println("========================================")

	c.Cloud = nil

	f := func(args ...string) {
		cmd := root.Root(c)
		cmd.SetArgs(args)
		err := cmd.Execute()
		assert.Nil(t, err)
	}

	fmt.Println("INIT")
	f("init", "--chdir", dir)
	fmt.Println("APPLY")
	f("apply", "--chdir", dir)
}

// testBackendMigration runs a testBackendMigration on a given directory and workspace name.
// t *testing.T is used for assertions.
// workspace is the name of the workspace used in the backend
// dir is the directory to run the testBackendMigration in
func testBackendMigration(t *testing.T, workspace, dir string) {
	// Directories
	localWorkspaceInit(dir)
	defer localWorkspaceClean(dir)

	// Environment
	// os.Setenv("DEBUG", "true")
	os.Setenv("KEY_ID", workspace)
	fmt.Printf("test key id: %s\n", workspace)

	// Config
	conf := helpers.Config()

	// Deploy (old backend)
	testBackendInit(t, conf, dir)

	fmt.Println("MIGRATE TO CLOUD")
	fmt.Println("========================================")

	// Wrap repeated calls
	f := func(args ...string) {
		cmd := root.Root(conf)
		cmd.SetArgs(args)
		err := cmd.Execute()
		assert.Nil(t, err)
	}

	// Run through
	f("init", "--chdir", dir)
	f("plan", "--chdir", dir)
	f("apply", "--chdir", dir)
	f("destroy", "--chdir", dir)

	// Cleanup
	f("delete", "--workspace", workspace, "--approve")
}

func TestTerraformS3CloudMigrationWithUUID(t *testing.T) {
	dir := t.TempDir()
	workspace := uuid.New().String()
	testBackendMigration(t, workspace, dir)
}

func TestTerraformS3CloudMigration(t *testing.T) {
	dir := t.TempDir()
	workspace := path.Base(path.Dir(dir))
	testBackendMigration(t, workspace, dir)
}
