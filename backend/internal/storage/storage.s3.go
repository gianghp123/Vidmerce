package storage

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
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

func (s *s3Storage) ObjectExists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			code := apiErr.ErrorCode()

			if code == "NotFound" || code == "NoSuchKey" || code == "AccessDenied" {
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}

func (s *s3Storage) DeleteObject(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}
