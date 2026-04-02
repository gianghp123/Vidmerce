package models

import "github.com/gianghp123/Vidmerce/backend/services/internal/core/enums"

type VideoMetadataEntity struct {
	BaseItem

	Title        string            `dynamodbav:"title"`
	Status       enums.VideoStatus `dynamodbav:"status"`
	VideoURL     string            `dynamodbav:"videoUrl"`
	ThumbnailURL string            `dynamodbav:"thumbnailUrl"`
	CreatedAt    string            `dynamodbav:"createdAt"`
}
