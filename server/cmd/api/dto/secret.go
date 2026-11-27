package dto

import "time"

type PasskeyInfo struct {
	Label     string     `json:"label"`
	CreatedAt time.Time  `json:"createdAt"`
	LastUsed  *time.Time `json:"lastUsed"`
	Current   bool       `json:"current"`
}

type SecretOutput struct {
	Passkeys []PasskeyInfo `json:"passkeys"`
	Message  string        `json:"message"`
}
