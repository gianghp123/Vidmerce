package models

type AssetSnapshot struct {
	AssetID    string  `dynamodbav:"assetId"`
	Name       string  `dynamodbav:"name"`
	Price      float64 `dynamodbav:"price"`
	ImageURL   string  `dynamodbav:"imageUrl"`
	ProductURL string  `dynamodbav:"productUrl"`
}
