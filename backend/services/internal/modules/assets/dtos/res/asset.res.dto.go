package res

type CreateAssetRes struct {
	AssetID  string `json:"assetId"`
	Status   string `json:"status"`
	ImageURL string `json:"imageUrl"`
}

type AssetRes struct {
	AssetID    string  `json:"assetId"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	ImageURL   string  `json:"imageUrl"`
	ProductURL string  `json:"productUrl"`
	CreatedAt  string  `json:"createdAt"`
}
