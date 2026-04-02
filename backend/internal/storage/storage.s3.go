package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Storage struct {
	client *s3.Client
	bucket string
}

func NewS3Storage(client *s3.Client, bucket string) Storage {
	return &s3Storage{
		client: client,
		bucket: bucket,
	}
}

func (s *s3Storage) GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, expire time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client, func(o *s3.PresignOptions) {
		o.Expires = expire
	})

	req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	return req.URL, nil
}

func (s *s3Storage) GeneratePresignedGetURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client, func(o *s3.PresignOptions) {
		o.Expires = expire
	})

	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", err
	}

	return req.URL, nil
}
