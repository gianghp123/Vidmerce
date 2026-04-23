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

func GetCurrentTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

var Now = GetCurrentTimestamp

func BuildUserBaseItem(userID string) models.BaseItem {
	return models.BaseItem{
		Pk:     string(core.PkPrefixUser) + core.KeySeparator + userID,
		Sk:     string(core.SkPrefixUser) + core.KeySeparator + userID,
		Gsi1Pk: string(core.Gsi1PkEntityUser),
		Gsi1Sk: userID,
	}
}

func BuildAssetBaseItem(assetID string) models.BaseItem {
	return BuildBaseItem(core.EntityTypeAsset, assetID, core.Gsi1PkEntityAsset)
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

func BuildUserJobBaseItem(userID, jobID string) models.BaseItem {
	return models.BaseItem{
		Pk:     string(core.PkPrefixUser) + core.KeySeparator + userID,
		Sk:     string(core.SkPrefixJob) + core.KeySeparator + jobID,
		Gsi1Pk: string(core.Gsi1PkEntityJob),
		Gsi1Sk: jobID,
	}
}

func BuildCampaignBaseItem(campaignID string) models.BaseItem {
	return BuildBaseItem(core.EntityTypeCampaign, campaignID, core.Gsi1PkEntityCampaign)
}

func BuildUserCampaignBaseItem(userID, campaignID string) models.BaseItem {
	return models.BaseItem{
		Pk:     string(core.PkPrefixUser) + core.KeySeparator + userID,
		Sk:     string(core.SkPrefixCampaign) + core.KeySeparator + campaignID,
		Gsi1Pk: string(core.Gsi1PkEntityCampaign),
		Gsi1Sk: campaignID,
	}
}

func BuildImageBaseItem(assetID string, order int) models.BaseItem {
	return models.BaseItem{
		Pk: string(core.PkPrefixAsset) + core.KeySeparator + assetID,
		Sk: string(core.SkPrefixImage) + core.KeySeparator + assetID + core.KeySeparator + string(rune('0'+order)),
	}
}

func BuildUserImageBaseItem(userID, assetID string, order int) models.BaseItem {
	return models.BaseItem{
		Pk: string(core.PkPrefixUser) + core.KeySeparator + userID,
		Sk: string(core.SkPrefixImage) + core.KeySeparator + assetID + core.KeySeparator + string(rune('0'+order)),
	}
}
