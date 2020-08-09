package config

import (
	"os"
	"strconv"
	"strings"
)

func Env(key string, defaultValue string) string {
	var value = os.Getenv(key)
	if len(strings.TrimSpace(value)) == 0 {
		return defaultValue
	}
	return value
}

func EnvInt(key string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err == nil {
		return value
	}
	return defaultValue
}

func EnvBool(key string, defaultValue bool) bool {
	var value = os.Getenv(key)
	if len(strings.TrimSpace(value)) == 0 {
		return defaultValue
	}
	return value == "true"
}
