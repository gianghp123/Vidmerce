package configs

import (
	"os"
	"strconv"
)

type AWSConfig struct {
	IsLocal  bool
	Endpoint *string
	Region   string
}

func LoadAWSConfig() *AWSConfig {
	isLocalStr := os.Getenv("IS_LOCAL")

	isLocal := false
	if val, err := strconv.ParseBool(isLocalStr); err == nil {
		isLocal = val
	}

	// NOTE: env name (you wrote AWS_ENTPOINT, assuming typo -> AWS_ENDPOINT)
	endpoint := os.Getenv("AWS_ENDPOINT")

	var endpointPtr *string
	if endpoint != "" {
		endpointPtr = &endpoint
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-southeast-1"
	}

	return &AWSConfig{
		IsLocal:  isLocal,
		Endpoint: endpointPtr,
		Region:   region,
	}
}
