package core

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/gianghp123/Vidmerce/backend/internal/configs"
	"github.com/gianghp123/Vidmerce/backend/internal/storage"
)

type Clients struct {
	S3 storage.Storage
	DB *dynamodb.Client
}

func NewClients(ctx context.Context, cfg *configs.AWSConfig, s3Bucket string) (*Clients, error) {
	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return nil, err
	}

	// DynamoDB client
	db := dynamodb.NewFromConfig(sdkConfig, func(o *dynamodb.Options) {
		if cfg.Endpoint != nil {
			o.BaseEndpoint = cfg.Endpoint
		}
	})

	// S3 client
	s3Client := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		if cfg.Endpoint != nil {
			o.BaseEndpoint = cfg.Endpoint
			o.UsePathStyle = true
		}
	})

	return &Clients{
		DB: db,
		S3: storage.NewS3Storage(s3Client, s3Bucket),
	}, nil
}
