package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aarondl/null/v8"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
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
		response.BadRequest(w, fmt.Errorf("authentication ceremony has not been started or has expired"))
		return
	}
	app.sessionManager.Remove(r.Context(), authenticationSessionDataKey)

	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(r.Body)
	if err != nil {
		response.BadRequest(w, fmt.Errorf("invalid authentication response"))
		return
	}

	var lookupErr error
	lookupUser := func(rawID, userHandle []byte) (webauthn.User, error) {
		var credential *models.Credential
		credential, lookupErr = models.Credentials(
			models.CredentialWhere.CredID.EQ(rawID),
			models.CredentialWhere.WebauthnUserID.EQ(userHandle),
		).One(r.Context(), tx)
		if lookupErr != nil {
			return nil, lookupErr
		}
		return toWebAuthnUserWithCredentials(credential), nil
	}

	user, credential, err := app.webAuthn.ValidatePasskeyLogin(lookupUser, sessionData, parsedResponse)
	if err != nil {
		if lookupErr != nil && !errors.Is(lookupErr, sql.ErrNoRows) {
			response.InternalServerError(w, lookupErr)
		} else {
			response.Unauthorized(w)
		}
		return
	}
	validatedUser, ok := user.(*WebAuthnUser)
	if !ok {
		response.InternalServerError(w, fmt.Errorf("unexpected WebAuthn user type %T", user))
		return
	}

	if credential.Authenticator.CloneWarning {
		response.Unauthorized(w)
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

	if err := app.sessionManager.RenewToken(r.Context()); err != nil {
		response.InternalServerError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), authenticatedUserIDKey, validatedUser.userID)
	w.WriteHeader(http.StatusOK)
}
