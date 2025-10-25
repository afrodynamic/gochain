package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("ADDR", "")
	t.Setenv("CHAIN", "")

	config := Load()

	if config.Address == "" || config.Chain == "" {
		t.Fatalf("unexpected default configuration: %+v", config)
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("ADDR", "127.0.0.1:9999")
	t.Setenv("CHAIN", "gochain")

	config := Load()

	if config.Address != "127.0.0.1:9999" || config.Chain != "gochain" {
		t.Fatalf("expected overrides not applied, got: %+v", config)
	}
}
