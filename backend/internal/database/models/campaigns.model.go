package models

import "github.com/gianghp123/Vidmerce/backend/internal/core/enums"

type CampaignEntity struct {
	BaseItem
	AssetID    string               `dynamodbav:"assetId"`
	Status     enums.CampaignStatus `dynamodbav:"status"`
	Storyboard Storyboard           `dynamodbav:"storyboard,omitempty"`
	CreatedAt  string               `dynamodbav:"createdAt"`
	UpdatedAt  string               `dynamodbav:"updatedAt,omitempty"`
}

type Storyboard struct {
	Headline string  `dynamodbav:"headline,omitempty"`
	BodyCopy string  `dynamodbav:"bodyCopy,omitempty"`
	CTA      string  `dynamodbav:"cta,omitempty"`
	Slides   []Slide `dynamodbav:"slides,omitempty"`
}
