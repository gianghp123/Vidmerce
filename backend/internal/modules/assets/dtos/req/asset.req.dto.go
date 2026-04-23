package req

type CreateAssetReq struct {
	ImageCount int `json:"imageCount" binding:"required,min=1,max=5" example:"3"`
}

type ConfirmUploadReq struct {
	Name       string   `json:"name" binding:"required"`
	Price      float64  `json:"price" binding:"required,gt=0"`
	ProductURL string   `json:"productUrl" binding:"required,url"`
	FileKeys   []string `json:"fileKeys" binding:"required,min=1"`
}

type ListAssetsQuery struct {
	Limit  int    `form:"limit,default=20" example:"20"`
	Cursor string `form:"cursor" example:"eyJpZCI6MTIzfQ=="`
}
