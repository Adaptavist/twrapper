package cmd

import (
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/cmd/root"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/utils"
	"github.com/spf13/cobra"
)

func Execute() {
	log := utils.GetLogger()
	cnf := config.Must(config.New(log))
	cmd := root.Root(cnf)
	cobra.CheckErr(cmd.Execute())
}
