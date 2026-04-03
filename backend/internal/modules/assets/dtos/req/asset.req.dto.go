package req

type CreateAssetReq struct {
	Name       string  `json:"name" binding:"required"`
	Price      float64 `json:"price" binding:"required,gt=0"`
	ProductURL string  `json:"productUrl" binding:"required,url"`
	ImageCount int     `json:"imageCount" binding:"min=1,max=10"`
}

type ListAssetsQuery struct {
	Limit  int    `form:"limit,default=20"`
	Cursor string `form:"cursor"`
}
