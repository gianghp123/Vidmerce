package models

import "github.com/gianghp123/Vidmerce/backend/internal/core/enums"

type SlideEntity struct {
	ImageRole enums.ImageRole `dynamodbav:"imageRole"`
	// Overlay text on image
	OverlayText *string `dynamodbav:"overlayText,omitempty"`
	// S3 key for image
	S3Key string `dynamodbav:"s3Key"`
	// SlideEntity number
	SlideNumber int `dynamodbav:"slideNumber"`
}
