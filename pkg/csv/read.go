package csv

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/pkg/text"
)

type S3Reader interface {
	ReadS3File(ctx context.Context, bucketName string, fileName string) ([]byte, error)
}

type s3Reader struct {
	config config.Config
	logs   logs.Logger
}

func NewS3Reader(cfg config.Config, logger logs.Logger) S3Reader {
	return &s3Reader{
		config: cfg,
		logs:   logger,
	}
}

func (s3Reader *s3Reader) ReadS3File(ctx context.Context, bucketName, fileName string) ([]byte, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s3Reader.config.MerchantScore.Region),
	})
	if err != nil {
		err = getSpecificError(err, bucketName, fileName)
		s3Reader.logs.Error(ctx, err.Error(), text.LogTagMethod, "ReadS3File")
		return nil, err
	}

	client := s3.New(sess)

	rawObject, err := client.GetObject(
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(fileName),
		})
	if err != nil {
		err = getSpecificError(err, bucketName, fileName)
		s3Reader.logs.Error(ctx, err.Error(), text.LogTagMethod, "ReadS3File")
		return nil, err
	}

	fileBytes, err := io.ReadAll(rawObject.Body)
	if err != nil {
		s3Reader.logs.Error(ctx, err.Error(), text.LogTagMethod, "ReadS3File")
		return nil, err
	}

	return fileBytes, nil
}

func getSpecificError(err error, bucket, fileName string) error {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case s3.ErrCodeNoSuchBucket:
			return fmt.Errorf("bucket '%s' does not exist", bucket)
		case s3.ErrCodeNoSuchKey:
			return fmt.Errorf("object with key '%s' does not exist in bucket '%s'", fileName, bucket)
		}
	}
	return err
}
