package configs

import (
	"os"
)

// GetSingleEnvConfigs is used get single enviroment var
func GetSingleEnvConfigs(key string) string {
	value := os.Getenv(key)
	return value
}
