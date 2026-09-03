package database

import (
	"net/url"
	"testing"
	"webauthn.rasc.ch/internal/config"
)

func TestDSNEscapesCredentialsAndDatabaseName(t *testing.T) {
	var cfg config.Config
	cfg.DB.User = "test@example.com"
	cfg.DB.Password = "p@ss/word"
	cfg.DB.Host = "localhost:5432"
	cfg.DB.Database = "passkey demo"

	parsed, err := url.Parse(DSN(cfg))
	if err != nil {
		t.Fatalf("parse DSN: %v", err)
	}

	if parsed.User.Username() != cfg.DB.User {
		t.Fatalf("username = %q, want %q", parsed.User.Username(), cfg.DB.User)
	}
	password, ok := parsed.User.Password()
	if !ok || password != cfg.DB.Password {
		t.Fatalf("password = %q, %v; want %q, true", password, ok, cfg.DB.Password)
	}
	if parsed.Host != cfg.DB.Host {
		t.Fatalf("host = %q, want %q", parsed.Host, cfg.DB.Host)
	}
	if parsed.Path != "/"+cfg.DB.Database {
		t.Fatalf("path = %q, want %q", parsed.Path, "/"+cfg.DB.Database)
	}
}
