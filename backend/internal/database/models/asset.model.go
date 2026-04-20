package models

import "github.com/gianghp123/Vidmerce/backend/internal/core/enums"

type AssetEntity struct {
	BaseItem
	// Creation timestamp
	CreatedAt string `dynamodbav:"createdAt"`
	// Number of images
	ImageCount int `dynamodbav:"imageCount"`
	// Product name
	Name string `dynamodbav:"name"`
	// Product price
	Price float64 `dynamodbav:"price"`
	// Product URL
	ProductURL string            `dynamodbav:"productUrl"`
	Status     enums.AssetStatus `dynamodbav:"status"`
}
