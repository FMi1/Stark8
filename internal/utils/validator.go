package utils

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

func CustomPasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Controlla la lunghezza della password
	if len(password) < 5 || len(password) > 30 {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// Controlla che ci sia almeno una maiuscola, una minuscola, un numero e un carattere speciale
	return hasUpper && hasLower && hasNumber && hasSpecial
}
