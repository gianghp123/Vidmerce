package configs

import (
	"os"
	"strconv"
)

type AWSConfig struct {
	IsLocal bool
}

func LoadAWSConfig() *AWSConfig {

	isLocalStr := os.Getenv("IS_LOCAL")

	isLocal := false // default
	if val, err := strconv.ParseBool(isLocalStr); err == nil {
		isLocal = val
	}

	return &AWSConfig{
		IsLocal: isLocal,
	}
}
