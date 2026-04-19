package enums

type EntityType string

const (
	EntityTypeAsset    EntityType = "ASSET"
	EntityTypeImage    EntityType = "IMAGE"
	EntityTypeJob      EntityType = "JOB"
	EntityTypeCampaign EntityType = "CAMPAIGN"
)

type GSI1PKPrefix string

const (
	GSI1PKEntityAsset    GSI1PKPrefix = "ENTITY#ASSET"
	GSI1PKEntityJob      GSI1PKPrefix = "ENTITY#JOB"
	GSI1PKEntityCampaign GSI1PKPrefix = "ENTITY#CAMPAIGN"
)

type SortKey string

const (
	SortKeyMetadata SortKey = "METADATA"
	SortKeyImage    SortKey = "IMAGE"
)

const KeySeparator string = "#"
