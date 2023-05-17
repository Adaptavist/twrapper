package migrator

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.uber.org/zap"
)

func (m *Migrator) s3ObjectInfo() (bucket, key string, err error) {
	if m.config.Backend == nil {
		return bucket, key, fmt.Errorf("backend config is nil")
	}

	bucket, err = m.config.Backend.Props.GetString("bucket")
	if err != nil {
		return bucket, key, err
	}

	key, err = m.config.Backend.Props.GetString("key")
	if err != nil {
		return bucket, key, err
	}

	return bucket, key, err
}

func (m *Migrator) initBackendS3(ctx context.Context) error {
	if m.s3 == nil {
		m.s3 = utils.GetS3Client(m.config)
	}

	// Get bucket details
	bucket, key, err := m.s3ObjectInfo()
	if err != nil {
		return err
	}

	object, err := m.s3.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		// If error is "NoSuchKey", we can skip decoding the file
		var bne *types.NoSuchKey
		if errors.As(err, &bne) {
			return nil
		}
		// if the error is anything else, return the error
		return err
	}

	m.backendStateBytes, err = io.ReadAll(object.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(m.backendStateBytes, &m.backendState); err != nil {
		return err
	}
	return nil
}

func (m *Migrator) cleanupS3() error {
	if m.s3 == nil {
		return errors.New("s3 client is nil, so not migration took place")
	}

	// Get bucket details
	bucket, key, err := m.s3ObjectInfo()
	if err != nil {
		return err
	}
	backupKey := fmt.Sprintf("%s.backup", key)

	_, err = m.s3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(backupKey),
		Body:                 bytes.NewReader(m.backendStateBytes),
		ServerSideEncryption: types.ServerSideEncryptionAes256,
	})

	if err != nil {
		return err
	}

	m.config.Logger.Info("backup state uploaded to backend", zap.String("bucket", bucket), zap.String("key", backupKey))

	_, err = m.s3.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return err
	}

	m.config.Logger.Info("deleted state file from backend", zap.String("bucket", bucket), zap.String("key", key))

	return err
}
