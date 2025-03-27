package util

import "strconv"

func ParseIntWithDefault(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return parsed
}
