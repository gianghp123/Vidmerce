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

type PkPrefix string

const (
	PkPrefixUser     PkPrefix = "USER"
	PkPrefixAsset    PkPrefix = "ASSET"
	PkPrefixCampaign PkPrefix = "CAMPAIGN"
	PkPrefixJob      PkPrefix = "JOB"
)

type SkPrefix string

const (
	SkPrefixUser     SkPrefix = "USER"
	SkPrefixAsset    SkPrefix = "ASSET"
	SkPrefixImage    SkPrefix = "IMAGE"
	SkPrefixCampaign SkPrefix = "CAMPAIGN"
	SkPrefixJob      SkPrefix = "JOB"
)

type Gsi1PkPrefix string

const (
	Gsi1PkEntityAsset    Gsi1PkPrefix = "ENTITY#ASSET"
	Gsi1PkEntityJob      Gsi1PkPrefix = "ENTITY#JOB"
	Gsi1PkEntityCampaign Gsi1PkPrefix = "ENTITY#CAMPAIGN"
	Gsi1PkEntityUser     Gsi1PkPrefix = "ENTITY#USER"
)

const KeySeparator string = "#"

type SortKey = SkPrefix

const (
	SortKeyMetadata = SkPrefix("METADATA")
	SortKeyImage    = SkPrefix("IMAGE")
)

const (
	UserIDKey = "user_id"
	RoleKey   = "role"
)
