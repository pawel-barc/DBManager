package utils

import (
	"time"

	"safebase/config"

	"github.com/golang-jwt/jwt/v5"
)

// ----- CREATETOKEN -----  crée un JWT pour un utilisateur donné
func CreateToken(userID uint, username string, duration time.Duration) (string, error) {
	// claim -> body du JWT
	claims := jwt.MapClaims{
		"user_id": int(userID),
		"username": username,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}

	// Création du token avec la méthode de signature HS256 et les claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Retourne le token signé sous forme de string et une erreur si ça échoue
	return token.SignedString([]byte(config.JWTSecret))
}