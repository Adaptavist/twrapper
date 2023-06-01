package backend

import (
	"context"
	"fmt"

	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

func List(config *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List items in the configure backend",
		Run: func(cmd *cobra.Command, args []string) {
			// Get bucket name
			bucket := config.MakeBackend().GetString("bucket")
			utils.FatalIfEmpty(bucket, "Unable to get bucket from config")

			// Get Client
			client := utils.GetS3Client(config)

			// Use the built in pagination functions for AWS
			pager := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
				Bucket: aws.String(bucket),
			})

			// Start paginating
			for pager.HasMorePages() {
				page, err := pager.NextPage(context.TODO())
				utils.FatalIfNotNil(err, "failed to list objects")
				for _, obj := range page.Contents {
					fmt.Printf("%s\n", *obj.Key)
				}
			}
		},
	}
}
