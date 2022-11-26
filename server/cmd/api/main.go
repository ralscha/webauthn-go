package main

import (
	"context"
	"database/sql"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
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
	"webauthn.rasc.ch/internal/version"
)

type application struct {
	config         *config.Config
	database       *sql.DB
	sessionManager *scs.SessionManager
	wg             sync.WaitGroup
	taskScheduler  chrono.TaskScheduler
}

func main() {
	rand.Seed(time.Now().UnixNano())

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
		logger.Error("opening database connection failed", err)
		os.Exit(1)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	logger.Info("database connection pool established")

	sm := scs.New()
	sm.Store = postgresstore.NewWithCleanupInterval(db, 30*time.Minute)
	sm.Lifetime = cfg.SessionLifetime
	sm.Cookie.SameSite = http.SameSiteStrictMode
	if cfg.CookieDomain != "" {
		sm.Cookie.Domain = cfg.CookieDomain
	}
	sm.Cookie.Secure = cfg.SecureCookie
	logger.Info("secure cookie", "secure", sm.Cookie.Secure)

	err = initAuth(cfg)
	if err != nil {
		logger.Error("init auth failed", err)
		os.Exit(1)
	}

	app := &application{
		config:         &cfg,
		database:       db,
		sessionManager: sm,
		taskScheduler:  chrono.NewDefaultTaskScheduler(),
	}

	_, err = app.taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		app.cleanup()
	}, 20*time.Minute)

	if err != nil {
		logger.Error("scheduling cleanup task failed", err)
		os.Exit(1)
	}

	logger.Info("starting server", "addr", app.config.HTTP.Port, "version", version.Get())

	err = app.serve()
	if err != nil {
		logger.Error("http serve failed", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}
