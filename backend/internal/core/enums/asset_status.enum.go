package enums

type AssetStatus string

const (
	AssetStatusCompleted AssetStatus = "COMPLETED"
	AssetStatusDraft     AssetStatus = "DRAFT"
	AssetStatusFailed    AssetStatus = "FAILED"
	AssetStatusUploading AssetStatus = "UPLOADING"
	AssetStatusImporting AssetStatus = "IMPORTING"
	AssetStatusPartial   AssetStatus = "PARTIAL"
)
