package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Variable globale pour le secret du JWT
var JWTSecret []byte

// La func init charge la valeur du JWT depuis l'environnement / .env
func init() {
	// On tente de charger un fichier .env (s'il existe)
	// mais on ne considère plus ça comme bloquant.
	_ = godotenv.Load()

	// On regarde dans l'environnement si on est en test ou pas
	env := os.Getenv("APP_ENV") // "test", "dev", "prod", etc.

	// On récupère le secret
	secret := os.Getenv("JWT_SECRET")

	// Si aucun secret n'est défini
	if secret == "" {
		if env == "test" {
			// En mode test : on met un secret par défaut
			secret = "default_test_secret"
			log.Println("⚠️ JWT_SECRET non défini, utilisation d'un secret de test par défaut")
		} else {
			// En dev / prod : on bloque, c'est une erreur critique
			log.Fatal("JWT_SECRET n'est pas défini dans l'environnement (.env)")
		}
	}

	// On stocke le secret dans la variable globale
	JWTSecret = []byte(secret)
}
