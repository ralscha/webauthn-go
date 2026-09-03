package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Environment string

const (
	Production  Environment = "production"
	Development Environment = "development"
)

type Config struct {
	Environment Environment
	Session     struct {
		SecureCookie bool
		CookieDomain string
		Lifetime     time.Duration
	}
	DB struct {
		User         string
		Password     string
		Host         string
		Database     string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
		MaxLifetime  string
	}
	HTTP struct {
		Addr                  string
		ReadTimeoutInSeconds  int64
		WriteTimeoutInSeconds int64
		IdleTimeoutInSeconds  int64
	}
	WebAuthn struct {
		RPID          string
		RPDisplayName string
		RPOrigins     string
	}
}

func applyDefaults(v *viper.Viper) {
	v.SetDefault("environment", Production)
	v.SetDefault("http.readTimeoutInSeconds", 30)
	v.SetDefault("http.writeTimeoutInSeconds", 30)
	v.SetDefault("http.idleTimeoutInSeconds", 120)
	v.SetDefault("db.maxOpenConns", 4)
	v.SetDefault("db.maxIdleConns", 2)
	v.SetDefault("db.maxIdleTime", "15m")
	v.SetDefault("db.maxLifetime", "2h")
	v.SetDefault("session.secureCookie", true)
	v.SetDefault("session.lifetime", "24h")
}

func LoadConfig() (Config, error) {
	v := newConfigViper()

	err := v.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	return decodeConfig(v)
}

func newConfigViper() *viper.Viper {
	v := viper.New()
	applyDefaults(v)
	v.SetConfigName("app")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetEnvPrefix("WEBAUTHN")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	return v
}

func decodeConfig(v *viper.Viper) (Config, error) {
	var cfg Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		return cfg, err
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg Config) validate() error {
	if cfg.Environment != Development && cfg.Environment != Production {
		return fmt.Errorf("environment must be %q or %q", Development, Production)
	}

	required := map[string]string{
		"db.user":                cfg.DB.User,
		"db.host":                cfg.DB.Host,
		"db.database":            cfg.DB.Database,
		"http.addr":              cfg.HTTP.Addr,
		"webauthn.rpId":          cfg.WebAuthn.RPID,
		"webauthn.rpDisplayName": cfg.WebAuthn.RPDisplayName,
		"webauthn.rpOrigins":     cfg.WebAuthn.RPOrigins,
	}
	for key, value := range required {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s must not be empty", key)
		}
	}

	if cfg.Session.Lifetime <= 0 {
		return fmt.Errorf("session.lifetime must be greater than zero")
	}
	if cfg.HTTP.ReadTimeoutInSeconds <= 0 || cfg.HTTP.WriteTimeoutInSeconds <= 0 || cfg.HTTP.IdleTimeoutInSeconds <= 0 {
		return fmt.Errorf("HTTP timeouts must be greater than zero")
	}
	if cfg.DB.MaxOpenConns <= 0 || cfg.DB.MaxIdleConns < 0 || cfg.DB.MaxIdleConns > cfg.DB.MaxOpenConns {
		return fmt.Errorf("database connection-pool settings are invalid")
	}
	if _, err := time.ParseDuration(cfg.DB.MaxIdleTime); err != nil {
		return fmt.Errorf("invalid db.maxIdleTime: %w", err)
	}
	if _, err := time.ParseDuration(cfg.DB.MaxLifetime); err != nil {
		return fmt.Errorf("invalid db.maxLifetime: %w", err)
	}

	return nil
}
