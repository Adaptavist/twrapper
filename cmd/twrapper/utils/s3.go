package utils

import (
	"context"
	"io"
	"log"

	cfg "github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetS3Client(c *cfg.Config) *s3.Client {
	// Get bucket name from config
	backend := c.MakeBackend()

	// Region
	region := ""
	if region = backend.GetString("region"); region == "" {
		log.Fatal("Unable to get region from config")
	}

	// Load the Shared AWS Configuration (~/.aws/config)
	conf, err := config.LoadDefaultConfig(context.TODO())
	conf.Region = region

	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	return s3.NewFromConfig(conf)
}

func GetS3ObjectContents(c *cfg.Config, backend map[string]interface{}) (body []byte, err error) {
	client := GetS3Client(c)
	object, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(backend["bucket"].(string)),
		Key:    aws.String(backend["key"].(string)),
	})

	if err != nil {
		return
	}

	body, err = io.ReadAll(object.Body)

	return
}
