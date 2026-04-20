package req

type CreateCampaignReq struct {
	AssetID string `json:"assetId" binding:"required" example:"uuid"`
}

type UpdateCampaignReq struct {
	Storyboard *UpdateStoryboard `json:"storyboard,omitempty"`
}

type UpdateStoryboard struct {
	Headline string  `json:"headline,omitempty"`
	BodyCopy string  `json:"bodyCopy,omitempty"`
	Cta      string  `json:"cta,omitempty"`
	Slides   []Slide `json:"slides,omitempty"`
}

type Slide struct {
	SlideNumber int    `json:"slideNumber"`
	ImageRole   string `json:"imageRole"`
	S3Key       string `json:"s3Key"`
	OverlayText string `json:"overlayText,omitempty"`
}
