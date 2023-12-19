package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"net/http"
	"webauthn.rasc.ch/internal/models"
	"webauthn.rasc.ch/internal/response"
)

const loginSessionDataKey = "webAuthnLoginSessionData"

func (app *application) loginStart(w http.ResponseWriter, r *http.Request) {
	options, sessionData, err := app.webAuthn.BeginDiscoverableLogin(webauthn.WithUserVerification(protocol.VerificationPreferred))
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), loginSessionDataKey, sessionData)
	response.JSON(w, http.StatusOK, options)
}

func (app *application) loginFinish(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	sessionData, ok := app.sessionManager.Get(r.Context(), loginSessionDataKey).(webauthn.SessionData)
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
	userID := bytesToInt64(parsedResponse.Response.UserHandle)

	err = models.AppCredentials(
		models.AppCredentialWhere.AppUserID.EQ(userID),
		models.AppCredentialWhere.ID.EQ(credential.ID),
	).
		UpdateAll(r.Context(), tx,
			models.M{models.AppCredentialColumns.SignCount: credential.Authenticator.SignCount},
		)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), loginSessionDataKey)
	app.sessionManager.Put(r.Context(), "userID", userID)
	w.WriteHeader(http.StatusOK)
}

func (app *application) createDiscovarableUserHandler(ctx context.Context, tx *sql.Tx) webauthn.DiscoverableUserHandler {
	return func(rawID, userHandle []byte) (webauthn.User, error) {
		userID := bytesToInt64(userHandle)
		user, err := models.FindAppUser(ctx, tx, userID)
		if err != nil {
			return nil, err
		}
		allUserCredentials, err := models.AppCredentials(models.AppCredentialWhere.AppUserID.EQ(user.ID)).All(ctx, tx)
		if err != nil {
			return nil, err
		}
		return toWebAuthnUserWithCredentials(user, allUserCredentials)
	}
}
