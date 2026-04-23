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
		Sk:     string(core.SkPrefix(entityType)) + core.KeySeparator + id,
		Gsi1Pk: string(gsi1PkPrefix),
		Gsi1Sk: id,
	}
}

var Now = func() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func BuildUserBaseItem(userID string) models.BaseItem {
	return models.BaseItem{
		Pk:     string(core.PkPrefixUser) + core.KeySeparator + userID,
		Sk:     string(core.SkPrefixUser) + core.KeySeparator + userID,
		Gsi1Pk: string(core.Gsi1PkEntityUser),
		Gsi1Sk: userID,
	}
}

func BuildUserAssetBaseItem(userID, assetID string) models.BaseItem {
	return models.BaseItem{
		Pk:     string(core.PkPrefixUser) + core.KeySeparator + userID,
		Sk:     string(core.SkPrefixAsset) + core.KeySeparator + assetID,
		Gsi1Pk: string(core.Gsi1PkEntityAsset),
		Gsi1Sk: assetID,
	}
}

func BuildJobBaseItem(jobID string) models.BaseItem {
	return BuildBaseItem(core.EntityTypeJob, jobID, core.Gsi1PkEntityJob)
}

func BuildCampaignBaseItem(campaignID string) models.BaseItem {
	return BuildBaseItem(core.EntityTypeCampaign, campaignID, core.Gsi1PkEntityCampaign)
}
