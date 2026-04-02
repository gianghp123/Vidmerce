package configs

import (
	"os"
)

type S3Config struct {
	BucketName string
}

func LoadS3Config() *S3Config {
	return &S3Config{
		BucketName: os.Getenv("S3_BUCKET_NAME"),
	}
}
