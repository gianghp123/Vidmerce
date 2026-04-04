package configs

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWSConfig struct {
	Region   string
	Endpoint *string
	IsLocal  bool
}

func LoadAWSConfig() *AWSConfig {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-southeast-1" // Default region
	}

	lsHost := os.Getenv("LOCALSTACK_HOSTNAME")

	if lsHost == "" {
		lsHost = "4566"
	}

	var endpoint *string
	isLocal := false

	if lsHost != "" {
		isLocal = true
		addr := os.Getenv("AWS_ENDPOINT_URL")
		endpoint = aws.String(addr)
	}

	return &AWSConfig{
		Region:   region,
		Endpoint: endpoint,
		IsLocal:  isLocal,
	}
}

// Helper để áp dụng cho DynamoDB
func (c *AWSConfig) DynamoDBOptions(o *dynamodb.Options) {
	o.Region = c.Region
	if c.IsLocal {
		o.BaseEndpoint = c.Endpoint
	}
}

// Helper để áp dụng cho S3
func (c *AWSConfig) S3Options(o *s3.Options) {
	o.Region = c.Region
	if c.IsLocal {
		o.BaseEndpoint = c.Endpoint
		o.UsePathStyle = true // Bắt buộc cho LocalStack S3
	}
}
