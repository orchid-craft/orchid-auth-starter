package common

import (
	"os"
	"strconv"
)

func GetBoolEnv(env string, fallback ...bool) bool {
	valueEnv := os.Getenv(env)
	fallbackValue := false
	if len(fallback) > 0 {
		fallbackValue = fallback[0]
	}

	if valueEnv == "" {
		return fallbackValue
	}

	result, err := strconv.ParseBool(valueEnv)
	if err != nil {
		return fallbackValue
	}
	return result
}

func GetEnvWithDefault(env string, fallback ...string) string {
	valueEnv := os.Getenv(env)
	if valueEnv == "" {
		if len(fallback) > 0 {
			valueEnv = fallback[0]
		}
	}
	return valueEnv
}
