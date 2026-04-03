package models

import "github.com/gianghp123/Vidmerce/backend/services/internal/core/enums"

type AssetEntity struct {
	BaseItem

	Name       string            `dynamodbav:"name"`
	Price      float64           `dynamodbav:"price"`
	ProductURL string            `dynamodbav:"productUrl"`
	Status     enums.AssetStatus `dynamodbav:"status"`
	ImageCount int               `dynamodbav:"imageCount"`
	CreatedAt  string            `dynamodbav:"createdAt"`
}

type ImageEntity struct {
	BaseItem

	FileKey string            `dynamodbav:"fileKey"`
	Status  enums.ImageStatus `dynamodbav:"status"`
	Order   int               `dynamodbav:"order"`
}
