package enums

type JobType string

const (
	JobTypeGenerateCampaign JobType = "GENERATE_CAMPAIGN"
	JobTypeScrapeProduct    JobType = "SCRAPE_PRODUCT"
)
