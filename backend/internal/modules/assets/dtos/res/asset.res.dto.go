package res

type UploadInfo struct {
	FileKey   string `json:"fileKey"`
	UploadURL string `json:"uploadUrl"`
	Order     int    `json:"order"`
	ExpiresIn int    `json:"expiresIn"`
}

type ImageInfo struct {
	ImageURL string `json:"imageUrl"`
	Order    int    `json:"order"`
}

type CreateAssetRes struct {
	AssetID string       `json:"assetId"`
	Status  string       `json:"status"`
	Uploads []UploadInfo `json:"uploads"`
}

type ConfirmAssetRes struct {
	AssetID string      `json:"assetId"`
	Status  string      `json:"status"`
	Images  []ImageInfo `json:"images"`
}

type AssetRes struct {
	AssetID    string      `json:"assetId"`
	Name       string      `json:"name"`
	Price      float64     `json:"price"`
	Images     []ImageInfo `json:"images,omitempty"`
	ProductURL string      `json:"productUrl"`
	Status     string      `json:"status"`
	CreatedAt  string      `json:"createdAt"`
}
