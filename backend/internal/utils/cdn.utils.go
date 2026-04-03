package utils

import (
	"github.com/gianghp123/Vidmerce/backend/services/internal/configs"
)

func GetCDNURL(fileKey string) string {
	cfg := configs.LoadS3Config()
	if cfg.CDNURL == "" {
		return fileKey
	}
	return cfg.CDNURL + "/" + fileKey
}
