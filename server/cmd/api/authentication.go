package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"net/http"
	"time"
	"webauthn.rasc.ch/internal/models"
	"webauthn.rasc.ch/internal/response"
)

const authenticationSessionDataKey = "webAuthnAuthenticationSessionData"

func (app *application) authenticationStart(w http.ResponseWriter, r *http.Request) {
	options, sessionData, err := app.webAuthn.BeginDiscoverableLogin(webauthn.WithUserVerification(protocol.VerificationPreferred))
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), authenticationSessionDataKey, sessionData)
	response.JSON(w, http.StatusOK, options.Response)
}

func (app *application) authenticationFinish(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	sessionData, ok := app.sessionManager.Get(r.Context(), authenticationSessionDataKey).(webauthn.SessionData)
	if !ok {
		err := fmt.Errorf("webAuthn session data not found")
		response.InternalServerError(w, err)
		return
	}

	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(r.Body)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	credential, err := app.webAuthn.ValidateDiscoverableLogin(app.createDiscovarableUserHandler(r.Context(), tx), sessionData, parsedResponse)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	if credential.Authenticator.CloneWarning {
		response.InternalServerError(w, fmt.Errorf("authenticator may be cloned"))
		return
	}

	cols := models.M{
		models.CredentialColumns.SignCount:    credential.Authenticator.SignCount,
		models.CredentialColumns.CloneWarning: credential.Authenticator.CloneWarning,
		models.CredentialColumns.LastUsed: null.Time{
			Time:  time.Now(),
			Valid: true,
		},
	}
	if credential.Flags.BackupEligible {
		cols[models.CredentialColumns.BackupState] = credential.Flags.BackupState
	}

	err = models.Credentials(
		models.CredentialWhere.WebauthnUserID.EQ(parsedResponse.Response.UserHandle),
		models.CredentialWhere.CredID.EQ(credential.ID),
	).
		UpdateAll(r.Context(), tx, cols)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), authenticationSessionDataKey)

	user, err := models.Credentials(models.CredentialWhere.CredID.EQ(credential.ID), qm.Select(models.CredentialColumns.UserID)).One(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "userID", user.UserID)
	w.WriteHeader(http.StatusOK)
}

func (app *application) createDiscovarableUserHandler(ctx context.Context, tx *sql.Tx) webauthn.DiscoverableUserHandler {
	return func(rawID, userHandle []byte) (webauthn.User, error) {
		credential, err := models.Credentials(models.CredentialWhere.WebauthnUserID.EQ(userHandle)).One(ctx, tx)
		if err != nil {
			return nil, err
		}
		return toWebAuthnUserWithCredentials(credential)
	}
}
