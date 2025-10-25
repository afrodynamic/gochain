package config

import "os"

type Config struct {
	Address string
	Chain   string
}

func Load() Config {
	return Config{
		Address: getEnvironmentVariable("ADDR", "127.0.0.1:8080"),
		Chain:   getEnvironmentVariable("CHAIN", "gochain"),
	}
}

func getEnvironmentVariable(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
