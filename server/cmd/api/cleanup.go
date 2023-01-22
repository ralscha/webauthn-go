package main

import (
	"context"
	"github.com/volatiletech/null/v8"
	"golang.org/x/exp/slog"
	"time"
	"webauthn.rasc.ch/internal/models"
)

func (app *application) cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete all users with a pending registration older than 10 minutes
	tenMinutesAgo := time.Now().Add(-10 * time.Minute)
	err := models.AppUsers(models.AppUserWhere.SignUpStart.LT(null.Time{
		Time:  tenMinutesAgo,
		Valid: true,
	})).DeleteAll(ctx, app.database)
	if err != nil {
		slog.Error("error deleting old pending sign ups", err)
	}
}
