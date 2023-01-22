package main

import (
	"context"
	"database/sql"
	"encoding/gob"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/procyon-projects/chrono"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/exp/slog"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
	"webauthn.rasc.ch/internal/config"
	"webauthn.rasc.ch/internal/database"
)

type application struct {
	config         *config.Config
	database       *sql.DB
	sessionManager *scs.SessionManager
	wg             sync.WaitGroup
	taskScheduler  chrono.TaskScheduler
	webAuthn       *webauthn.WebAuthn
}

func main() {
	rand.Seed(time.Now().UnixNano())
	gob.Register(webauthn.SessionData{})
	gob.Register(webauthn.Credential{})

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("reading config failed %v\n", err)
	}

	var logger *slog.Logger

	switch cfg.Environment {
	case config.Development:
		boil.DebugMode = true
		logger = slog.New(slog.NewTextHandler(os.Stdout))
	case config.Production:
		logger = slog.New(slog.NewJSONHandler(os.Stdout))
	}

	slog.SetDefault(logger)

	db, err := database.New(cfg)
	if err != nil {
		slog.Error("opening database connection failed", err)
		os.Exit(1)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	slog.Info("database connection pool established")

	sm := scs.New()
	sm.Store = postgresstore.NewWithCleanupInterval(db, 30*time.Minute)
	sm.Lifetime = cfg.SessionLifetime
	sm.Cookie.SameSite = http.SameSiteStrictMode
	if cfg.CookieDomain != "" {
		sm.Cookie.Domain = cfg.CookieDomain
	}
	sm.Cookie.Secure = cfg.SecureCookie
	slog.Info("secure cookie", "secure", sm.Cookie.Secure)

	wa, err := webauthn.New(&webauthn.Config{
		RPDisplayName: cfg.WebAuthn.RPDisplayName,
		RPID:          cfg.WebAuthn.RPID,
		RPOrigins:     []string{cfg.WebAuthn.RPOrigins},
	})
	if err != nil {
		slog.Error("initializing webauthn failed", err)
		os.Exit(1)
	}

	app := &application{
		config:         &cfg,
		database:       db,
		sessionManager: sm,
		taskScheduler:  chrono.NewDefaultTaskScheduler(),
		webAuthn:       wa,
	}

	_, err = app.taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		app.cleanup()
	}, 20*time.Minute)

	if err != nil {
		slog.Error("scheduling cleanup task failed", err)
		os.Exit(1)
	}

	err = app.serve()
	if err != nil {
		slog.Error("http serve failed", err)
		os.Exit(1)
	}
}
