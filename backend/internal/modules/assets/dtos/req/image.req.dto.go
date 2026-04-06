package req

type GetImageUploadUrlReq struct {
	FileName string `form:"fileName" binding:"required" example:"image.jpg"`
}

type UploadAssetImageReq struct {
	FileKey string `json:"fileKey" binding:"required" example:"assets/abc123/1.jpg"`
}
