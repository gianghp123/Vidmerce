package enums

type AssetStatus string

const (
	StatusAssetUploading AssetStatus = "UPLOADING"
	StatusAssetDraft     AssetStatus = "DRAFT"
	StatusAssetCompleted AssetStatus = "COMPLETED"
)

type ImageStatus string

const (
	StatusImageUploading ImageStatus = "UPLOADING"
	StatusImageCompleted ImageStatus = "COMPLETED"
)
