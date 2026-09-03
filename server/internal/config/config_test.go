package config

import (
	"strings"
	"testing"
)

const validConfig = `
environment: development
db:
  user: postgres
  password: password
  host: config-host:5432
  database: webauthn
http:
  addr: localhost:8080
session:
  secureCookie: false
webauthn:
  rpId: localhost
  rpDisplayName: Test Application
  rpOrigins: http://localhost:4200
`

func TestDecodeConfigReadsHostAndEnvironmentOverrides(t *testing.T) {
	t.Setenv("WEBAUTHN_DB_HOST", "environment-host:5432")

	v := newConfigViper()
	if err := v.ReadConfig(strings.NewReader(validConfig)); err != nil {
		t.Fatalf("read test config: %v", err)
	}

	cfg, err := decodeConfig(v)
	if err != nil {
		t.Fatalf("decode config: %v", err)
	}

	if cfg.DB.Host != "environment-host:5432" {
		t.Fatalf("DB host = %q, want environment override", cfg.DB.Host)
	}
	if cfg.Session.Lifetime.String() != "24h0m0s" {
		t.Fatalf("session lifetime = %s, want default of 24h", cfg.Session.Lifetime)
	}
}

func TestDecodeConfigRejectsUnknownEnvironment(t *testing.T) {
	v := newConfigViper()
	config := strings.Replace(validConfig, "environment: development", "environment: staging", 1)
	if err := v.ReadConfig(strings.NewReader(config)); err != nil {
		t.Fatalf("read test config: %v", err)
	}

	_, err := decodeConfig(v)
	if err == nil {
		t.Fatal("decodeConfig succeeded with an unknown environment")
	}
}
