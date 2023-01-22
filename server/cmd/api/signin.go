package main

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"net/http"
	"webauthn.rasc.ch/cmd/api/dto"
	"webauthn.rasc.ch/internal/models"
	"webauthn.rasc.ch/internal/request"
	"webauthn.rasc.ch/internal/response"
)

func (app *application) signInStart(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)

	var usernameInput dto.UsernameInput
	if ok := request.DecodeJSONValidate[*dto.UsernameInput](w, r, &usernameInput, dto.ValidateUsernameInput); !ok {
		return
	}

	// Find the user
	user, err := models.AppUsers(models.AppUserWhere.Username.EQ(usernameInput.Username)).One(r.Context(), tx)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Unauthorized(w)
			return
		}
		response.InternalServerError(w, err)
		return
	}

	credentials, err := models.AppCredentials(models.AppCredentialWhere.AppUserID.EQ(user.ID)).All(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	webAuthnUser, err := toWebAuthnUserWithCredentials(user, credentials)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	options, sessionData, err := app.webAuthn.BeginLogin(webAuthnUser)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "webAuthnSignInSessionData", sessionData)
	response.JSON(w, http.StatusOK, options)
}

func (app *application) signInFinish(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	sessionData, ok := app.sessionManager.Get(r.Context(), "webAuthnSignInSessionData").(webauthn.SessionData)
	if !ok {
		err := fmt.Errorf("webAuthn session data not found")
		response.InternalServerError(w, err)
		return
	}

	userID := bytesToInt64(sessionData.UserID)
	user, err := models.FindAppUser(r.Context(), tx, userID)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	allUserCredentials, err := models.AppCredentials(models.AppCredentialWhere.AppUserID.EQ(user.ID)).All(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	webAuthnUser, err := toWebAuthnUserWithCredentials(user, allUserCredentials)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(r.Body)
	credential, err := app.webAuthn.ValidateLogin(webAuthnUser, sessionData, parsedResponse)
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

	err = models.AppCredentials(models.AppCredentialWhere.AppUserID.EQ(user.ID), models.AppCredentialWhere.ID.EQ(credential.ID)).
		UpdateAll(r.Context(), tx, models.M{models.AppCredentialColumns.Credential: byteBuffer.Bytes()})
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "webAuthnSignInSessionData")
	app.sessionManager.Put(r.Context(), "userID", user.ID)
	w.WriteHeader(http.StatusOK)
}
