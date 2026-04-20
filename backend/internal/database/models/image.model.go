package models

import "github.com/gianghp123/Vidmerce/backend/internal/core/enums"

type ImageEntity struct {
	BaseItem
	// S3 file key
	FileKey string `dynamodbav:"fileKey"`
	// ImageEntity order
	Order  int               `dynamodbav:"order"`
	Status enums.ImageStatus `dynamodbav:"status"`
}
