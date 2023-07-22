package utils

import "strconv"

func GetInt(str string, fallbackValue int) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		return fallbackValue
	}

	return value
}
