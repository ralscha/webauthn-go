package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"net/http"
	"time"
	"webauthn.rasc.ch/cmd/api/dto"
	"webauthn.rasc.ch/internal/models"
	"webauthn.rasc.ch/internal/request"
	"webauthn.rasc.ch/internal/response"
)

const registrationSessionDataKey = "webAuthnRegistrationSessionData"
const registrationSessionUserId = "webAuthnRegistrationSessionUserId"

func (app *application) registrationStart(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)

	var usernameInput dto.UsernameInput
	if ok := request.DecodeJSONValidate[*dto.UsernameInput](w, r, &usernameInput, dto.ValidateUsernameInput); !ok {
		return
	}

	user := models.User{
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

	rnd := make([]byte, 64)
	if _, err := rand.Read(rnd); err != nil {
		response.InternalServerError(w, err)
		return
	}
	webAuthnUser := &WebAuthnUser{
		username: user.Username,
		id:       rnd,
	}

	requireResidentKey := true
	options, sessionData, err := app.webAuthn.BeginRegistration(webAuthnUser, webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
		ResidentKey:        protocol.ResidentKeyRequirementRequired,
		RequireResidentKey: &requireResidentKey,
		UserVerification:   protocol.VerificationPreferred,
	}), webauthn.WithConveyancePreference(protocol.PreferNoAttestation),
		webauthn.WithExclusions([]protocol.CredentialDescriptor{}),
		webauthn.WithExtensions(protocol.AuthenticationExtensions{"credProps": true}),
	)

	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), registrationSessionDataKey, sessionData)
	app.sessionManager.Put(r.Context(), registrationSessionUserId, user.ID)

	response.JSON(w, http.StatusOK, options.Response)
}

func (app *application) registrationFinish(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)

	options, ok := app.sessionManager.Get(r.Context(), registrationSessionDataKey).(webauthn.SessionData)
	if !ok {
		err := fmt.Errorf("webAuthn session data not found")
		response.InternalServerError(w, err)
		return
	}
	userId, ok := app.sessionManager.Get(r.Context(), registrationSessionUserId).(int)
	if !ok {
		err := fmt.Errorf("webAuthn session user id not found")
		response.InternalServerError(w, err)
		return
	}

	user, err := models.FindUser(r.Context(), tx, userId)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	webAuthnUser := &WebAuthnUser{
		username: user.Username,
		id:       options.UserID,
	}

	credential, err := app.webAuthn.FinishRegistration(webAuthnUser, options, r)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	transports := ""
	for i, t := range credential.Transport {
		if i > 0 {
			transports += ","
		}
		transports += string(t)
	}

	appCredential := models.Credential{
		CredID:         credential.ID,
		UserID:         user.ID,
		WebauthnUserID: options.UserID,
		LastUsed: null.Time{
			Time:  time.Now(),
			Valid: true,
		},
		Aaguid: null.Bytes{
			Bytes: credential.Authenticator.AAGUID,
			Valid: len(credential.Authenticator.AAGUID) > 0,
		},
		AttestationType: null.String{
			String: credential.AttestationType,
			Valid:  credential.AttestationType != "",
		},
		Attachment:     string(credential.Authenticator.Attachment),
		Transport:      transports,
		SignCount:      int(credential.Authenticator.SignCount),
		CloneWarning:   credential.Authenticator.CloneWarning,
		Present:        credential.Flags.UserPresent,
		Verified:       credential.Flags.UserVerified,
		BackupEligible: credential.Flags.BackupEligible,
		BackupState:    credential.Flags.BackupState,
		PublicKey:      credential.PublicKey,
	}
	if err := appCredential.Insert(r.Context(), tx, boil.Infer()); err != nil {
		response.InternalServerError(w, err)
		return
	}

	err = models.Users(models.UserWhere.ID.EQ(user.ID)).
		UpdateAll(r.Context(), tx, models.M{models.UserColumns.RegistrationStart: null.Time{Valid: false}})
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), registrationSessionDataKey)
	app.sessionManager.Remove(r.Context(), registrationSessionUserId)
	w.WriteHeader(http.StatusOK)
}
