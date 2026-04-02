package enums

type VideoStatus string

const (
	StatusPending    VideoStatus = "PENDING"
	StatusProcessing VideoStatus = "PROCESSING"
	StatusCompleted  VideoStatus = "COMPLETED"
	StatusFailed     VideoStatus = "FAILED"
)
