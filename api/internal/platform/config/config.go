package config

import "os"

type Config struct {
	Addr  string
	Chain string
}

func Load() Config {
	return Config{
		Addr:  getenv("ADDR", "127.0.0.1:8080"),
		Chain: getenv("CHAIN", "gochain"),
	}
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
