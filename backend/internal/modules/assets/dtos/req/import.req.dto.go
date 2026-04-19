package req

type ImportAssetReq struct {
	ProductURL string `json:"productUrl" binding:"required,url" example:"https://..."`
}
