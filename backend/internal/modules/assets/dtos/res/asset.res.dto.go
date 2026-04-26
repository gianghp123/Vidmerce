package res

type UploadInfo struct {
	FileKey   string `json:"fileKey"`
	UploadURL string `json:"uploadUrl"`
	Order     int    `json:"order"`
	ExpiresIn int    `json:"expiresIn"`
}

type ImageInfo struct {
	ID       string `json:"id"`
	ImageURL string `json:"imageUrl"`
}

type PresignedUrlsRes struct {
	AssetID string       `json:"assetId"`
	Uploads []UploadInfo `json:"uploads"`
}

type AssetRes struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Price      float64     `json:"price"`
	Images     []ImageInfo `json:"images"`
	ProductURL string      `json:"productUrl"`
	Status     string      `json:"status"`
	CreatedAt  string      `json:"createdAt"`
}
