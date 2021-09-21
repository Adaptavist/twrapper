package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdAllArgs(t *testing.T) {
	r := bytes.NewBufferString("")
	rootCmd.SetOut(r)

	err := os.Setenv("KEY_ID", "2F4A2B8D-9353-44EA-8248-386CB0FE7425")
	assert.Nil(t, err)

	rootCmd.SetArgs([]string{"terraform", "-version"})
	err = rootCmd.Execute()
	assert.Nil(t, err)

	_, err = ioutil.ReadAll(r)
	assert.Nil(t, err)
}
