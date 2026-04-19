package enums

type JobStatus string

const (
	StatusJobPending    JobStatus = "PENDING"
	StatusJobProcessing JobStatus = "PROCESSING"
	StatusJobCompleted  JobStatus = "COMPLETED"
	StatusJobFailed     JobStatus = "FAILED"
)

type JobType string

const (
	TypeJobScrapeProduct    JobType = "SCRAPE_PRODUCT"
	TypeJobGenerateCampaign JobType = "GENERATE_CAMPAIGN"
)
