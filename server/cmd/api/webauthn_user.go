package main

import (
	"encoding/binary"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"strings"
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
		var transports []protocol.AuthenticatorTransport
		if c.Transports.Valid {
			for _, t := range strings.Split(c.Transports.String, ",") {
				transports = append(transports, protocol.AuthenticatorTransport(t))
			}
		}
		webAuthnCredentials[i] = webauthn.Credential{
			ID:              c.ID,
			PublicKey:       c.PublicKey,
			AttestationType: c.AttestationType.String,
			Transport:       transports,
			Authenticator: webauthn.Authenticator{
				AAGUID:    c.AaGUID,
				SignCount: uint32(c.SignCount),
			},
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
