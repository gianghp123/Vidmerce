package utils

import (
	"github.com/gianghp123/Vidmerce/backend/internal/core/enums"
	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
	"time"
)

func BuildPK(entityType enums.EntityType, id string) string {
	return string(entityType) + enums.KeySeparator + id
}

func BuildBaseItem(entityType enums.EntityType, id string, gsi1PKPrefix enums.GSI1PKPrefix) models.BaseItem {
	return models.BaseItem{
		PK:     BuildPK(entityType, id),
		SK:     string(enums.SortKeyMetadata),
		GSI1PK: string(gsi1PKPrefix),
		GSI1SK: id,
	}
}

func GetCurrentTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

var Now = GetCurrentTimestamp

func BuildAssetBaseItem(id string) models.BaseItem {
	return BuildBaseItem(enums.EntityTypeAsset, id, enums.GSI1PKEntityAsset)
}

func BuildJobBaseItem(id string) models.BaseItem {
	return BuildBaseItem(enums.EntityTypeJob, id, enums.GSI1PKEntityJob)
}

func BuildCampaignBaseItem(id string) models.BaseItem {
	return BuildBaseItem(enums.EntityTypeCampaign, id, enums.GSI1PKEntityCampaign)
}
