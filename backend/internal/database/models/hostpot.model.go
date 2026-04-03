package models

type Hotspot struct {
	X float64 `dynamodbav:"x"`
	Y float64 `dynamodbav:"y"`

	Asset AssetSnapshot `dynamodbav:"asset"`
}

type AssetSnapshot struct {
	AssetID    string  `dynamodbav:"assetId"`
	Name       string  `dynamodbav:"name"`
	Price      float64 `dynamodbav:"price"`
	ImageURL   string  `dynamodbav:"imageUrl"`
	ProductURL string  `dynamodbav:"productUrl"`
}
