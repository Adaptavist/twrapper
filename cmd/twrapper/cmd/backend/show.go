package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/utils"
	"github.com/adaptavist/terraform-wrapper/v1/pkg/terraform/model/state"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

func printStateAsTable(data []byte) {
	// Prep object for unmarshaling
	tfState := state.State{}

	// Unmarshal JSON
	err := json.Unmarshal(data, &tfState)
	utils.FatalIfNotNil(err, "failed to unmarshal state")

	resources := tfState.FilterByMode("managed")
	fmt.Println("")
	fmt.Printf("Resources: %d\n", len(resources))
	fmt.Println("----------------------------------------------------------------------------------------------------")
	fmt.Printf("| %-30s | %-63s |\n", "module", "key")
	fmt.Println("----------------------------------------------------------------------------------------------------")
	for _, resource := range resources {
		fmt.Printf("| %-30s | %-63s |\n", resource.Module, resource.Key())
	}
	fmt.Println("----------------------------------------------------------------------------------------------------")
	fmt.Println("")
	fmt.Printf("Outputs: %d\n", len(tfState.Outputs))
}

func Show(config *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show objects in the configure backend",
		Long:  "Show objects by key E.G twrapper backend show key1 key2",
		Run: func(cmd *cobra.Command, args []string) {
			// Get bucket name
			bucket := config.MakeBackend().GetString("bucket")
			utils.FatalIfEmpty(bucket, "Unable to get bucket from config")

			// Get s3 client
			client := utils.GetS3Client(config)

			for _, object := range args {
				// Get Item
				item, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
					Bucket: aws.String(bucket),
					Key:    aws.String(object),
				})
				utils.DieOnError(err)

				// Get content of the item
				body, err := io.ReadAll(item.Body)
				utils.DieOnError(err)

				switch cmd.Flag("format").Value.String() {
				case "json":
					fmt.Println(string(body))
				case "table":
					printStateAsTable(body)
				default:
					fmt.Println(string(body))
				}
			}
		},
	}
	cmd.PersistentFlags().String("output", "json", "Output format")
	return cmd
}
