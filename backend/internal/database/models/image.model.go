package models

type ImageEntity struct {
	BaseItem
	// S3 file key
	FileKey string `dynamodbav:"fileKey"`
}
