package validation

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

// La validation du mot de passe
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("Le mot de passe doit contenir au moins huit caractères")
	}
	var hasLower, hasUpper, hasDigit, hasSpecial bool

	for _, c := range password {
		switch {
		case unicode.IsLower(c):
			hasLower = true

		case unicode.IsUpper(c):
			hasUpper = true
			
		case unicode.IsDigit(c) || unicode.IsSymbol(c):
			hasDigit = true
			
		case unicode.IsPunct(c):
			hasSpecial = true	
		}
	}

	if !hasLower || !hasUpper || !hasDigit || !hasSpecial {
		return errors.New("le mot de passe doit contenir au moins huit caractères, une lettre majuscule, une lettre minuscule, un chiffre et un carctère spécial")
	}
	return nil
}
// La validation d'email
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

	match, _ := regexp.MatchString(regex, email)
	if !match {
		return errors.New("email invalide")
	}
	return nil
}

// La validation du nom
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if len(username) < 2 {
		return errors.New("le nom doit contenir au moins deux caractères")
	}
	return nil
}