package main

import (
	"testing"

	"webauthn.rasc.ch/internal/models"
)

func TestToWebAuthnUserOmitsEmptyTransport(t *testing.T) {
	credential := &models.Credential{
		UserID:         42,
		CredID:         []byte("credential-id"),
		WebauthnUserID: []byte("user-handle"),
		PublicKey:      []byte("public-key"),
		Transport:      "",
	}

	user := toWebAuthnUserWithCredentials(credential)

	if user.userID != credential.UserID {
		t.Fatalf("internal user ID = %d, want %d", user.userID, credential.UserID)
	}
	credentials := user.WebAuthnCredentials()
	if len(credentials) != 1 {
		t.Fatalf("credentials = %d, want 1", len(credentials))
	}
	if len(credentials[0].Transport) != 0 {
		t.Fatalf("empty database transport produced %#v", credentials[0].Transport)
	}
}
