package main

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"net/http"
	"webauthn.rasc.ch/internal/models"
	"webauthn.rasc.ch/internal/response"
)

func (app *application) signInStart(w http.ResponseWriter, r *http.Request) {
	options, sessionData, err := app.webAuthn.BeginDiscoverableLogin(webauthn.WithUserVerification(protocol.VerificationPreferred))
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "webAuthnSignInSessionData", sessionData)
	response.JSON(w, http.StatusOK, options)
}

func (app *application) createDiscovarableUserHandler(r *http.Request) webauthn.DiscoverableUserHandler {
	return func(rawID, userHandle []byte) (webauthn.User, error) {
		tx := r.Context().Value(transactionKey).(*sql.Tx)

		userID := bytesToInt64(userHandle)
		user, err := models.FindAppUser(r.Context(), tx, userID)
		if err != nil {
			return nil, err
		}
		allUserCredentials, err := models.AppCredentials(models.AppCredentialWhere.AppUserID.EQ(user.ID)).All(r.Context(), tx)
		if err != nil {
			return nil, err
		}
		return toWebAuthnUserWithCredentials(user, allUserCredentials)
	}
}

func (app *application) signInFinish(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	sessionData, ok := app.sessionManager.Get(r.Context(), "webAuthnSignInSessionData").(webauthn.SessionData)
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
	credential, err := app.webAuthn.ValidateDiscoverableLogin(app.createDiscovarableUserHandler(r), sessionData, parsedResponse)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	// update credential
	byteBuffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(byteBuffer)
	err = encoder.Encode(credential)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	userID := bytesToInt64(parsedResponse.Response.UserHandle)
	err = models.AppCredentials(models.AppCredentialWhere.AppUserID.EQ(userID), models.AppCredentialWhere.ID.EQ(credential.ID)).
		UpdateAll(r.Context(), tx, models.M{models.AppCredentialColumns.Credential: byteBuffer.Bytes()})
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "webAuthnSignInSessionData")
	app.sessionManager.Put(r.Context(), "userID", userID)
	w.WriteHeader(http.StatusOK)
}
