package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"

	"github.com/adaptavist/terraform-wrapper/v1/pkg/terraform"
	"github.com/spf13/viper"
)

var config Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "atw",
	Short: "Wrapper for initialising Terraform",
	Long:  `Sets some things up before running Terraform in a CI environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := config.checkRequirements()
		fatalIfNotNil(err, "%s")

		opts := terraform.NewOpts()
		opts.WithArguments(ActionArgs(args))
		opts.WithBackend(config.Backend())

		if !config.AWS.IsEmpty() {
			fmt.Println("setting up for aws")
			err = config.AWS.Configure(&opts)
			fatalIfNotNil(err, "failed to configure Terraform for aws: %s")
		}

		err = terraform.Configure(opts)
		fatalIfNotNil(err, "%s")

		err = terraform.Init()
		fatalIfNotNil(err, "%s")

		err = terraform.Execute(opts)
		fatalIfNotNil(err, "%s")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cwd, err := os.Getwd()
	fatalIfNotNil(err, "%s")
	viper.AddConfigPath(cwd)
	viper.SetConfigName("terraform.atw.yml")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()

	if err != nil {
		log.Fatalf("failed to find %s with error: %s", viper.ConfigFileUsed(), err)
	}

	if err = viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable load %s", viper.ConfigFileUsed())
	}
}
