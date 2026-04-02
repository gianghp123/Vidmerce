package models

type AssetEntity struct {
	BaseItem

	Name       string  `dynamodbav:"name"`
	Price      float64 `dynamodbav:"price"`
	ImageURL   string  `dynamodbav:"imageUrl"`
	ProductURL string  `dynamodbav:"productUrl"`
	CreatedAt  string  `dynamodbav:"createdAt"`
}
