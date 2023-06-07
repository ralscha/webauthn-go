package main

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gobuffalo/validate"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"net/http"
	"time"
	"webauthn.rasc.ch/cmd/api/dto"
	"webauthn.rasc.ch/internal/models"
	"webauthn.rasc.ch/internal/request"
	"webauthn.rasc.ch/internal/response"
)

func (app *application) signUpStart(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)

	var usernameInput dto.UsernameInput
	if ok := request.DecodeJSONValidate[*dto.UsernameInput](w, r, &usernameInput, dto.ValidateUsernameInput); !ok {
		return
	}

	// check if username is already taken
	usernameExists, err := models.AppUsers(models.AppUserWhere.Username.EQ(usernameInput.Username)).Exists(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	if usernameExists {
		usernameExistsError := validate.NewErrors()
		usernameExistsError.Add("username", "exists")
		response.FailedValidation(w, usernameExistsError)
		return
	}

	user := models.AppUser{
		Username: usernameInput.Username,
		SignUpStart: null.Time{
			Time:  time.Now(),
			Valid: true,
		},
	}
	if err := user.Insert(r.Context(), tx, boil.Infer()); err != nil {
		response.InternalServerError(w, err)
		return
	}

	webAuthnUser := toWebAuthnUser(&user)

	options, sessionData, err := app.webAuthn.BeginRegistration(webAuthnUser, webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
		ResidentKey:      protocol.ResidentKeyRequirementRequired,
		UserVerification: protocol.VerificationPreferred,
	}))
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "webAuthnSignUpSessionData", sessionData)
	response.JSON(w, http.StatusOK, options)
}

func (app *application) signUpFinish(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	sessionData, ok := app.sessionManager.Get(r.Context(), "webAuthnSignUpSessionData").(webauthn.SessionData)
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
	webAuthnUser := toWebAuthnUser(user)

	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(r.Body)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	credential, err := app.webAuthn.CreateCredential(webAuthnUser, sessionData, parsedResponse)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	byteBuffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(byteBuffer)
	err = encoder.Encode(credential)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	appCredential := models.AppCredential{
		ID:         credential.ID,
		AppUserID:  user.ID,
		Credential: byteBuffer.Bytes(),
	}
	if err := appCredential.Insert(r.Context(), tx, boil.Infer()); err != nil {
		response.InternalServerError(w, err)
		return
	}

	err = models.AppUsers(models.AppUserWhere.ID.EQ(user.ID)).
		UpdateAll(r.Context(), tx, models.M{models.AppUserColumns.SignUpStart: null.Time{Valid: false}})
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "webAuthnSignUpSessionData")
	w.WriteHeader(http.StatusOK)
}
