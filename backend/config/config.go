package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Variable globale pour le secret du JWT
var JWTSecret []byte

// La func init charge la valeur du JWT depuis le fichier .env
// func init() {
// 	if err := godotenv.Load(); err != nil {
// 		log.Println("Echec du chargement du fichier .env, les valeurs par défaut seront utilisées")
// 	}

// 	secret := os.Getenv("JWT_SECRET")
// 	if secret == "" {
// 		log.Fatal("JWT_SECRET n'est pas défini dans le fichier .env")
// 	}

// 	JWTSecret = []byte(secret)
// }

// Func Load est temporaire, remplace le fonction init pour le temps de tests  

func Load() error {
	_ = godotenv.Load()

	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		return errors.New("JWT_SECRET n'est pas défini")
	}
	JWTSecret = []byte(secret)
	return nil
}