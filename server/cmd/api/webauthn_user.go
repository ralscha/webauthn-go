package main

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"webauthn.rasc.ch/internal/models"
)

type WebAuthnUser struct {
	username    string
	id          []byte
	credentials []webauthn.Credential
}

func (u *WebAuthnUser) WebAuthnID() []byte {
	return u.id
}

func (u *WebAuthnUser) WebAuthnName() string {
	return u.username
}

func (u *WebAuthnUser) WebAuthnDisplayName() string {
	return u.username
}

func (u *WebAuthnUser) WebAuthnIcon() string {
	return ""
}

func (u *WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

func toWebAuthnUserWithCredentials(credential *models.Credential) (*WebAuthnUser, error) {
	webAuthnCredential := webauthn.Credential{
		ID:        credential.CredID,
		PublicKey: credential.CredPublicKey,
		Authenticator: webauthn.Authenticator{
			SignCount: uint32(credential.Counter),
		},
	}

	return &WebAuthnUser{
		username:    "",
		id:          credential.WebauthnUserID,
		credentials: []webauthn.Credential{webAuthnCredential},
	}, nil
}
