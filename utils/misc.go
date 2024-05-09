package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GenKey(strs []string) string {
	// Concatenate JSON string and user ID with a separator
	combined := ""
	for _, v := range strs {
		combined += v
	}

	hash := md5.Sum([]byte(combined))
	return hex.EncodeToString(hash[:])
}
