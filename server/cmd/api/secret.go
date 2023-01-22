package main

import (
	"golang.org/x/exp/slog"
	"net/http"
	"webauthn.rasc.ch/cmd/api/dto"
	"webauthn.rasc.ch/internal/response"
)

func (app *application) secret(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.GetInt64(r.Context(), "userID")
	slog.Info("User ID", userID)

	response.JSON(w, http.StatusOK, dto.SecretOutput{
		Message: "This is a secret message",
	})
}
