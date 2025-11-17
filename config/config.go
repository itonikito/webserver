package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort    string
	BatFilePath   string
	TimeoutSec    int
	MaxConcurrent int
	LogLevel      string
}

func Load() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "5005"),
		//BatFilePath:   getEnv("BAT_FILE_PATH", "testdb"),
		//Для теста +
		BatFilePath: getEnv("BAT_FILE_PATH", "/Users/nikita/project/go/webserver/scripts/testdb.sh"),
		//Для теста-
		TimeoutSec:    getEnvAsInt("TIMEOUT_SEC", 30),
		MaxConcurrent: getEnvAsInt("MAX_CONCURRENT", 5),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
