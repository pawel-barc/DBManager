package validation

import "testing"

//////////////========PASSWORD=========//////////////

func TestValidatePassword_TooShort(t *testing.T) {
	password := "ABc1!"
	err := ValidatePassword(password)
	if err == nil {
		t.Errorf("Erreur attendu pour un mot de passe trop court")
	}
}

func TestValidatePassword_MissingDigit(t *testing.T) {
	password := "Mypassword!"
	err := ValidatePassword(password)
	if err == nil {
		t.Errorf("Erreur attendu pour un mot de passe sans un chiffre")
	}
}

func TestValidatePassword_MissingSpecialCharacter(t *testing.T) {
	password := "Mypassword1"
	err := ValidatePassword(password)
	if err == nil {
		t.Errorf("Erreur attendu pour un mot de passe sans un caractère spécial")
	}
}

func TestValidatePassword_MissingUpperCase(t *testing.T) {
	password := "mypassword1!"
	err := ValidatePassword(password)
	if err == nil {
		t.Errorf("Erreur attendu pour un mot de passe sans une lettre majuscule")
	}
}

func TestValidatePassword_MissingLowerCase(t *testing.T) {
	password := "MYPASSWORD1!"
	err := ValidatePassword(password)
	if err == nil {
		t.Errorf("Erreur attendu pour un mot de passe sans une lettre minuscule")
	}
}

func TestValidate_ValidPassword(t *testing.T) {
	password := "Mypassword1?"
	err := ValidatePassword(password)
	if err != nil {
		t.Errorf("Pas d'erreur attendu, %v", err)
	}
}

//////////////========EMAIL=========//////////////

func TestValidate_ValideEmail(t *testing.T) {
	email := "myemail@yahoo.com"
	err := ValidateEmail(email)
	if err != nil {
		t.Errorf("Pas d'erreur attendu, %v", err)
	}
}

func TestValidate_InvalideEmail(t *testing.T) {
	email := "myemailyahoo.com"
	err := ValidateEmail(email)
	if err == nil {
		t.Errorf("Email invalide")
	}
}

