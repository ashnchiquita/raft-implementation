package utils

import (
	"strings"
)

func IsValidPair(key string, value string) bool {
	return IsValidKeyOrValue(key) && IsValidKeyOrValue(value)
}

func IsValidKeyOrValue(keyOrValue string) bool {
	return keyOrValue != "" && !strings.Contains(keyOrValue, " ") && !strings.Contains(keyOrValue, ",")
}
