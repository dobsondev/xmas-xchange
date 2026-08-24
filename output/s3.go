package output

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// s3PutObjectAPI narrows *s3.Client to just what uploadToS3 needs, for testability.
type s3PutObjectAPI interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

// s3GetObjectAPI narrows *s3.Client to just what downloadFromS3 needs, for testability.
type s3GetObjectAPI interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func uploadToS3(ctx context.Context, client s3PutObjectAPI, bucket, key string, data []byte) (string, error) {
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("s3://%s/%s", bucket, key), nil
}

// WriteS3 uploads data to bucket/key in region, returning the s3:// location.
func WriteS3(ctx context.Context, region, bucket, key string, data []byte) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return "", err
	}

	return uploadToS3(ctx, s3.NewFromConfig(cfg), bucket, key, data)
}

func downloadFromS3(ctx context.Context, client s3GetObjectAPI, bucket, key string) ([]byte, error) {
	out, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()

	return io.ReadAll(out.Body)
}

// ReadS3 downloads bucket/key in region and returns its bytes.
func ReadS3(ctx context.Context, region, bucket, key string) ([]byte, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	return downloadFromS3(ctx, s3.NewFromConfig(cfg), bucket, key)
}

// ParseS3URI parses an "s3://bucket/key" location. ok is false if location
// isn't an s3:// URI or is missing a bucket/key.
func ParseS3URI(location string) (bucket, key string, ok bool) {
	const prefix = "s3://"
	if !strings.HasPrefix(location, prefix) {
		return "", "", false
	}

	rest := strings.TrimPrefix(location, prefix)
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}

	return parts[0], parts[1], true
}
