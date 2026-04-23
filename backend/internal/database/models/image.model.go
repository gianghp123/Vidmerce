package models

type ImageEntity struct {
	BaseItem
	FileKey string `dynamodbav:"fileKey"`
	Order   int    `dynamodbav:"order"`
}
