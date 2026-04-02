package models

type VideoMetadataEntity struct {
	BaseItem

	Title        string `dynamodbav:"title"`
	Status       string `dynamodbav:"status"`
	VideoURL     string `dynamodbav:"videoUrl"`
	ThumbnailURL string `dynamodbav:"thumbnailUrl"`
	CreatedAt    string `dynamodbav:"createdAt"`
}
