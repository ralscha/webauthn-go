package main

import (
	"encoding/binary"
	"encoding/json"
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

func toWebAuthnUser(user *models.AppUser) *WebAuthnUser {
	return &WebAuthnUser{
		username: user.Username,
		id:       int64ToBytes(user.ID),
	}
}

func toWebAuthnUserWithCredentials(user *models.AppUser, credentials []*models.AppCredential) (*WebAuthnUser, error) {
	webAuthnCredentials := make([]webauthn.Credential, len(credentials))
	for i, c := range credentials {
		err := json.Unmarshal([]byte(c.Credential), &webAuthnCredentials[i])
		if err != nil {
			return nil, err
		}
	}

	return &WebAuthnUser{
		username:    user.Username,
		id:          int64ToBytes(user.ID),
		credentials: webAuthnCredentials,
	}, nil
}

func int64ToBytes(i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

func bytesToInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}
