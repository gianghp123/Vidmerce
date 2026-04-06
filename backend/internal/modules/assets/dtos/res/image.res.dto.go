package res

type AssetUpload struct {
	UploadURL string `json:"uploadUrl" example:"https://s3.amazonaws.com/bucket/asset-url"`
	FileKey   string `json:"fileKey" example:"assets/abc123/1.jpg"`
	ExpiresIn int    `json:"expiresIn" example:"300"`
}

type ImageIdRes struct {
	ImageID string `json:"imageId" example:"img-123"`
}
