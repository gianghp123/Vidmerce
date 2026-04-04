package configs

import (
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWSConfig struct {
	Endpoint *string
	IsLocal  bool
}

func LoadAWSConfig() *AWSConfig {
	endpointStr := os.Getenv("AWS_ENDPOINT_URL")

	isLocalStr := os.Getenv("IS_LOCAL")

	isLocal := false // default
	if val, err := strconv.ParseBool(isLocalStr); err == nil {
		isLocal = val
	}

	var endpoint *string
	if isLocal {
		endpoint = aws.String(endpointStr)
	}

	return &AWSConfig{
		Endpoint: endpoint,
		IsLocal:  isLocal,
	}
}

// Helper để áp dụng cho DynamoDB
func (c *AWSConfig) DynamoDBOptions(o *dynamodb.Options) {
	if c.IsLocal {
		o.BaseEndpoint = c.Endpoint
	}
}

// Helper để áp dụng cho S3
func (c *AWSConfig) S3Options(o *s3.Options) {
	if c.IsLocal {
		o.BaseEndpoint = c.Endpoint
		o.UsePathStyle = true // Bắt buộc cho LocalStack S3
	}
}
