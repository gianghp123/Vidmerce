package req

type CreateAssetReq struct {
	Name       string  `json:"name" binding:"required" example:"iPhone 15"`
	Price      float64 `json:"price" binding:"required,gt=0" example:"999.99"`
	ProductURL string  `json:"productUrl" binding:"required,url" example:"https://example.com/products/iphone-15"`
	ImageCount int     `json:"imageCount" binding:"min=1,max=5" example:"3"`
}

type ListAssetsQuery struct {
	Limit  int    `form:"limit,default=20" example:"20"`
	Cursor string `form:"cursor" example:"eyJpZCI6MTIzfQ=="`
}
