package helpers

import (
	"bytes"
	"context"
	"errors"

	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/utils"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func WriteS3BackendObject(c *config.Config, content []byte) error {
	if c.Backend == nil {
		return errors.New("empty backend config")
	}

	bucket, err := c.Backend.Props.GetString("bucket")
	if err != nil {
		return err
	}

	key, err := c.Backend.Props.GetString("key")
	if err != nil {
		return err
	}

	return WriteS3Object(c, bucket, key, content)
}

// WriteS3Object writes an object to S3 for testing purposes
func WriteS3Object(c *config.Config, bucket, key string, content []byte) error {
	cli := utils.GetS3Client(c)

	req := &s3.PutObjectInput{
		Bucket:               &bucket,
		Key:                  &key,
		ServerSideEncryption: types.ServerSideEncryptionAes256,
		Body:                 bytes.NewReader(content),
	}

	_, err := cli.PutObject(context.TODO(), req)

	return err
}
