package models

import "github.com/gianghp123/Vidmerce/backend/internal/core/enums"

type CampaignEntity struct {
	BaseItem
	// Associated asset ID
	AssetID string `dynamodbav:"assetId"`
	// Creation timestamp
	CreatedAt string               `dynamodbav:"createdAt"`
	Status    enums.CampaignStatus `dynamodbav:"status"`
	// CampaignEntity storyboard
	StoryboardEntity *StoryboardEntity `dynamodbav:"storyboard,omitempty"`
	// Last update timestamp
	UpdatedAt *string `dynamodbav:"updatedAt,omitempty"`
}
