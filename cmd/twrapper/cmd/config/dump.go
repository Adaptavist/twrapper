package config

import (
	"fmt"
	"os"

	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func Check(c *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:                "dump",
		Short:              "Dump wrapper config to stdout",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			d, err := yaml.Marshal(c)
			utils.DieOnError(err)
			fmt.Println(string(d))
			os.Exit(0)
		},
	}
}
