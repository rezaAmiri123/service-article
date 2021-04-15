package utils

import (
	"strconv"
)

// Get config path for local or docker
func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	return "./cmd/config/config-local"
}

func UintToString(n uint64) string{
	return strconv.FormatUint(uint64(n), 10)
}

func StringToUint(str string) uint64{
	u, _ := strconv.ParseUint(str, 0, 64)
	return u
}