package req

type GenerateUploadUrlsReq struct {
	Count int `json:"count" binding:"required,min=1"`
}
