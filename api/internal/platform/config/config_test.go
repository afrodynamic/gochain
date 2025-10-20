package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("ADDR", "")
	t.Setenv("CHAIN", "")
	c := Load()
	if c.Addr == "" || c.Chain == "" {
		t.Fatalf("got=%+v", c)
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("ADDR", "127.0.0.1:9999")
	t.Setenv("CHAIN", "gochain")
	c := Load()
	if c.Addr != "127.0.0.1:9999" || c.Chain != "gochain" {
		t.Fatalf("got=%+v", c)
	}
}
