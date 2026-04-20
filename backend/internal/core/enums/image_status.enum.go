package enums

type ImageStatus string

const (
	ImageStatusCompleted ImageStatus = "COMPLETED"
	ImageStatusFailed    ImageStatus = "FAILED"
	ImageStatusUploading ImageStatus = "UPLOADING"
)
