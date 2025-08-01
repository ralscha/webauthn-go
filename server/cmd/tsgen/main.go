package main

import (
	"github.com/aarondl/null/v8"
	"github.com/gobuffalo/validate"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
	"webauthn.rasc.ch/cmd/api/dto"
)

func main() {
	converter := typescriptify.New()
	converter.Add(dto.SecretOutput{})
	converter.Add(dto.UsernameInput{})
	converter.Add(validate.Errors{})
	converter.CreateInterface = true
	converter.BackupDir = ""
	converter.ManageType(null.String{}, typescriptify.TypeOptions{TSType: "string"})

	err := converter.ConvertToFile("../client/src/app/api/types.ts")
	if err != nil {
		panic(err)
	}

}
