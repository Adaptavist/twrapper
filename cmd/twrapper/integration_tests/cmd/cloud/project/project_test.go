//go:build integration_test
// +build integration_test

package project

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/cmd/cloud"
	"github.com/adaptavist/terraform-wrapper/v1/pkg/integration/test/helpers"
	"github.com/stretchr/testify/assert"
)

func ProjectIntegrationTest(t *testing.T) {
	conf := helpers.Config()
	cmd := cloud.Root(conf)
	run := func(args ...string) {
		cmd.SetArgs(args)
		err := cmd.Execute()
		assert.Nil(t, err)
	}

	name := fmt.Sprintf("twrapper-integration-test-%x", rand.New(rand.NewSource(time.Now().UnixNano())))

	// Run some commands
	run("project", "list")
	run("project", "create", "--name", name)
	run("project", "delete", "--name", name, "--approve")
}
