package core

const TableName = "MediaProjectTable"
const MaxImagesPerAsset = 5

type EntityType string

const (
	EntityTypeAsset    EntityType = "ASSET"
	EntityTypeImage    EntityType = "IMAGE"
	EntityTypeJob      EntityType = "JOB"
	EntityTypeCampaign EntityType = "CAMPAIGN"
	EntityTypeUser     EntityType = "USER"
)

type Gsi1PkPrefix string

const (
	Gsi1PkEntityAsset    Gsi1PkPrefix = "ENTITY#ASSET"
	Gsi1PkEntityJob      Gsi1PkPrefix = "ENTITY#JOB"
	Gsi1PkEntityCampaign Gsi1PkPrefix = "ENTITY#CAMPAIGN"
	Gsi1PkEntityUser     Gsi1PkPrefix = "ENTITY#USER"
)

type SortKey string

const (
	SortKeyMetadata SortKey = "METADATA"
	SortKeyImage    SortKey = "IMAGE"
)

const KeySeparator string = "#"
