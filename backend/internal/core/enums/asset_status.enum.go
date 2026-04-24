package enums

type AssetStatus string

const (
	AssetStatusCompleted AssetStatus = "COMPLETED"
	AssetStatusFailed    AssetStatus = "FAILED"
	AssetStatusImporting AssetStatus = "IMPORTING"
)
