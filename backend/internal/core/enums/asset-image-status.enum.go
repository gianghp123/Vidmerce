package enums

type AssetStatus string

const (
	StatusAssetUploading AssetStatus = "UPLOADING"
	StatusAssetDraft     AssetStatus = "DRAFT"
	StatusAssetCompleted AssetStatus = "COMPLETED"
	StatusAssetPartial   AssetStatus = "PARTIAL"
	StatusAssetFailed    AssetStatus = "FAILED"
)

type ImageStatus string

const (
	StatusImageUploading ImageStatus = "UPLOADING"
	StatusImageCompleted ImageStatus = "COMPLETED"
	StatusImageFailed    ImageStatus = "FAILED"
)
