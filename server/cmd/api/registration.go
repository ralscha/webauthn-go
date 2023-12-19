package main

import (
	"database/sql"
	"fmt"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gobuffalo/validate"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"net/http"
	"strings"
	"time"
	"webauthn.rasc.ch/cmd/api/dto"
	"webauthn.rasc.ch/internal/models"
	"webauthn.rasc.ch/internal/request"
	"webauthn.rasc.ch/internal/response"
)

const registrationSessionDataKey = "webAuthnRegistrationSessionData"

func (app *application) registrationStart(w http.ResponseWriter, r *http.Request) {
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
		RegistrationStart: null.Time{
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

	app.sessionManager.Put(r.Context(), registrationSessionDataKey, sessionData)
	response.JSON(w, http.StatusOK, options)
}

func (app *application) registrationFinish(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)

	sessionData, ok := app.sessionManager.Get(r.Context(), registrationSessionDataKey).(webauthn.SessionData)
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

	credential, err := app.webAuthn.FinishRegistration(webAuthnUser, sessionData, r)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	var transportStrings []string
	for _, transport := range credential.Transport {
		transportStrings = append(transportStrings, string(transport))
	}
	transports := strings.Join(transportStrings, ",")

	appCredential := models.AppCredential{
		ID:        credential.ID,
		AppUserID: user.ID,
		PublicKey: credential.PublicKey,
		AttestationType: null.String{
			String: credential.AttestationType,
			Valid:  credential.AttestationType != "",
		},
		AaGUID:    credential.Authenticator.AAGUID,
		SignCount: int(credential.Authenticator.SignCount),
		Transports: null.String{
			String: transports,
			Valid:  transports != "",
		},
	}
	if err := appCredential.Insert(r.Context(), tx, boil.Infer()); err != nil {
		response.InternalServerError(w, err)
		return
	}

	err = models.AppUsers(models.AppUserWhere.ID.EQ(user.ID)).
		UpdateAll(r.Context(), tx, models.M{models.AppUserColumns.RegistrationStart: null.Time{Valid: false}})
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), registrationSessionDataKey)
	w.WriteHeader(http.StatusOK)
}
