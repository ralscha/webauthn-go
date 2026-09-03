package database

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/url"
	"time"
	"webauthn.rasc.ch/internal/config"
)

func DSN(cfg config.Config) string {
	return (&url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.DB.User, cfg.DB.Password),
		Host:   cfg.DB.Host,
		Path:   "/" + cfg.DB.Database,
	}).String()
}

func New(cfg config.Config) (*sql.DB, error) {
	connMaxIdleTime, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	connMaxLifetime, err := time.ParseDuration(cfg.DB.MaxLifetime)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", DSN(cfg))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime)
	db.SetConnMaxLifetime(connMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
