package utils

import (
	"time"

	"github.com/gianghp123/Vidmerce/backend/internal/core"
	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
)

func BuildPk(entityType core.EntityType, id string) string {
	return string(entityType) + core.KeySeparator + id
}

func BuildBaseItem(entityType core.EntityType, id string, gsi1PkPrefix core.Gsi1PkPrefix) models.BaseItem {
	return models.BaseItem{
		Pk:     BuildPk(entityType, id),
		Sk:     string(core.SortKeyMetadata),
		Gsi1Pk: string(gsi1PkPrefix),
		Gsi1Sk: id,
	}
}

func GetCurrentTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

var Now = GetCurrentTimestamp

func BuildAssetBaseItem(id string) models.BaseItem {
	return BuildBaseItem(core.EntityTypeAsset, id, core.Gsi1PkEntityAsset)
}

func BuildJobBaseItem(id string) models.BaseItem {
	return BuildBaseItem(core.EntityTypeJob, id, core.Gsi1PkEntityJob)
}

func BuildCampaignBaseItem(id string) models.BaseItem {
	return BuildBaseItem(core.EntityTypeCampaign, id, core.Gsi1PkEntityCampaign)
}
