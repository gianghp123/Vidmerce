package utils

import (
	"strings"

	"github.com/gianghp123/Vidmerce/backend/internal/core"
)

func GetOwnerIDFromPK(pk string) string {
	parts := strings.SplitN(pk, core.KeySeparator, 2)
	if len(parts) == 2 && parts[0] == string(core.PkPrefixUser) {
		return parts[1]
	}
	return ""
}
