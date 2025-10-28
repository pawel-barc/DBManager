package models

// Struct User représente un utilisateur dans l'application
type User struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password,omitempty"`
	RefreshToken string `json:"refresh_token"`
}
