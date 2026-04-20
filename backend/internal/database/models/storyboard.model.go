package models

type StoryboardEntity struct {
	BaseItem
	// CampaignEntity body copy
	BodyCopy string `dynamodbav:"bodyCopy"`
	// Call to action text
	Cta string `dynamodbav:"cta"`
	// CampaignEntity headline
	Headline string `dynamodbav:"headline"`
	// Associated job ID
	JobID string `dynamodbav:"jobId"`
	// StoryboardEntity slides
	Slides []SlideEntity `dynamodbav:"slides"`
}
