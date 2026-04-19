package models

import "github.com/gianghp123/Vidmerce/backend/internal/core/enums"

type StoryboardEntity struct {
	BaseItem
	JobID    string  `dynamodbav:"jobId"`
	Headline string  `dynamodbav:"headline"`
	BodyCopy string  `dynamodbav:"bodyCopy"`
	CTA      string  `dynamodbav:"cta"`
	Slides   []Slide `dynamodbav:"slides"`
}

type Slide struct {
	SlideNumber int             `dynamodbav:"slideNumber"`
	ImageRole   enums.ImageRole `dynamodbav:"imageRole"`
	S3Key       string          `dynamodbav:"s3Key"`
	OverlayText string          `dynamodbav:"overlayText,omitempty"`
}
