package configs

import (
	"os"
)

// GetAllEnvConfigs is used get all enviroment var
func GetAllEnvConfigs(key map[string]string) map[string]string {
	for k := range key {
		key[k] = os.Getenv(k)
	}
	return key
}

// GetSingleEnvConfigs is used get single enviroment var
func GetSingleEnvConfigs(key string) string {
	value := os.Getenv(key)
	return value
}
