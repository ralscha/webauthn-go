package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"net/http"
	"webauthn.rasc.ch/cmd/api/dto"
	"webauthn.rasc.ch/internal/models"
	"webauthn.rasc.ch/internal/response"
)

func (app *application) secret(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.GetInt(r.Context(), authenticatedUserIDKey)
	slog.Info("fetch secret", "User ID", userID)

	credentials, err := models.Credentials(models.CredentialWhere.UserID.EQ(userID)).All(r.Context(), app.database)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	current, _ := app.sessionManager.Get(r.Context(), currentPasskeyKey).([]byte)
	passkeys := make([]dto.PasskeyInfo, 0, len(credentials))
	for _, c := range credentials {
		label := "Passkey"
		switch c.Attachment {
		case "platform":
			label = "Device passkey"
		case "cross-platform":
			label = "External passkey"
		}
		id := sha256.Sum256(c.CredID)
		info := dto.PasskeyInfo{Label: fmt.Sprintf("%s · %x", label, id[:3]), CreatedAt: c.CreatedAt, Current: bytes.Equal(current, c.CredID)}
		if c.LastUsed.Valid {
			info.LastUsed = &c.LastUsed.Time
		}
		passkeys = append(passkeys, info)
	}
	response.JSON(w, http.StatusOK, dto.SecretOutput{Passkeys: passkeys,
		Message: "This is a secret message",
	})
}
