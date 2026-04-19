package res

type CampaignRes struct {
	ID         string     `json:"id"`
	AssetID    string     `json:"assetId"`
	Status     string     `json:"status"`
	Storyboard Storyboard `json:"storyboard,omitempty"`
	CreatedAt  string     `json:"createdAt"`
	UpdatedAt  string     `json:"updatedAt,omitempty"`
}

type Storyboard struct {
	Headline string  `json:"headline,omitempty"`
	BodyCopy string  `json:"bodyCopy,omitempty"`
	CTA      string  `json:"cta,omitempty"`
	Slides   []Slide `json:"slides,omitempty"`
}

type Slide struct {
	SlideNumber int    `json:"slideNumber"`
	ImageRole   string `json:"imageRole"`
	S3Key       string `json:"s3Key"`
	OverlayText string `json:"overlayText,omitempty"`
}
