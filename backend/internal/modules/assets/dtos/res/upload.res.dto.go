package res

type UploadUrlRes struct {
	AssetID   string `json:"assetId"`
	FileKey   string `json:"fileKey"`
	UploadURL string `json:"uploadUrl"`
	ExpiresIn int    `json:"expiresIn"`
}

type GenerateUploadUrlsRes struct {
	AssetID    string         `json:"assetId"`
	Uploads    []UploadUrlRes `json:"uploads"`
	MaxAllowed int            `json:"maxAllowed"`
	ExpiresIn  int            `json:"expiresIn"`
}
