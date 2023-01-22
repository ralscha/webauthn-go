package dto

import (
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type UsernameInput struct {
	Username string `json:"username"`
}

func ValidateUsernameInput(u *UsernameInput) *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{
			Name:    "username",
			Field:   u.Username,
			Message: "required",
		},
	)
}
