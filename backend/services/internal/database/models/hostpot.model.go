package models

type Hotspot struct {
	X float64 `dynamodbav:"x"`
	Y float64 `dynamodbav:"y"`

	Asset AssetSnapshot `dynamodbav:"asset"`
}
