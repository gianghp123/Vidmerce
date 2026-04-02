package req

type CreateAssetReq struct {
	AssetID    string  `json:"assetId" binding:"required"`
	Name       string  `json:"name" binding:"required"`
	Price      float64 `json:"price" binding:"required,gt=0"`
	ProductURL string  `json:"productUrl" binding:"required,url"`
}

type ListAssetsQuery struct {
	Limit  int    `form:"limit,default=20"`
	Cursor string `form:"cursor"`
}
