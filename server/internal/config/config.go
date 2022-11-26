package config

import (
	"github.com/spf13/viper"
	"time"
)

type Environment string

const (
	Production  Environment = "production"
	Development Environment = "development"
)

type Config struct {
	Environment     Environment
	SecureCookie    bool
	CookieDomain    string
	SessionLifetime time.Duration
	DB              struct {
		User         string
		Password     string
		Connection   string
		Database     string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
		MaxLifetime  string
	}
	HTTP struct {
		Port                  string
		ReadTimeoutInSeconds  int64
		WriteTimeoutInSeconds int64
		IdleTimeoutInSeconds  int64
	}
}

func applyDefaults() {
	viper.SetDefault("environment", Production)
	viper.SetDefault("http.readTimeoutInSeconds", 30)
	viper.SetDefault("http.writeTimeoutInSeconds", 30)
	viper.SetDefault("http.idleTimeoutInSeconds", 120)
	viper.SetDefault("sessionLifetime", "24h")
	viper.SetDefault("db.maxOpenConns", 4)
	viper.SetDefault("db.maxIdleConns", 2)
	viper.SetDefault("db.maxIdleTime", "15m")
	viper.SetDefault("db.maxLifetime", "2h")
	viper.SetDefault("secureCookie", true)
}

func LoadConfig() (Config, error) {
	var cfg Config

	applyDefaults()
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return cfg, err
	}

	viper.SetEnvPrefix("WEBAUTHN")
	viper.AutomaticEnv()

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
