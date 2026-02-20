package main

import (
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

func toWebAuthnUserWithCredentials(credential *models.Credential) (*WebAuthnUser, error) {
	attestationType := ""
	if credential.AttestationType.Valid {
		attestationType = credential.AttestationType.String
	}

	var transports []protocol.AuthenticatorTransport
	splitted := strings.SplitSeq(credential.Transport, ",")
	for s := range splitted {
		transports = append(transports, protocol.AuthenticatorTransport(s))
	}

	webAuthnCredential := webauthn.Credential{
		ID:              credential.CredID,
		PublicKey:       credential.PublicKey,
		AttestationType: attestationType,
		Transport:       transports,
		Flags: webauthn.CredentialFlags{
			UserPresent:    credential.Present,
			UserVerified:   credential.Verified,
			BackupEligible: credential.BackupEligible,
			BackupState:    credential.BackupState,
		},
		Authenticator: webauthn.Authenticator{
			AAGUID:       credential.Aaguid.Bytes,
			SignCount:    uint32(credential.SignCount),
			CloneWarning: credential.CloneWarning,
			Attachment:   protocol.AuthenticatorAttachment(credential.Attachment),
		},
		Attestation: webauthn.CredentialAttestation{},
	}

	return &WebAuthnUser{
		username:    "",
		id:          credential.WebauthnUserID,
		credentials: []webauthn.Credential{webAuthnCredential},
	}, nil
}
