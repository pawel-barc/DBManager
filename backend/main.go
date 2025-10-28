package main

import (
	"log"
	"net/http"

	"safebase/db"
	"safebase/routes"
)

func main() {
	// Connexiona à la base des données
	db.ConnectDB()
	// Configuration du routeur (définition des routes de l'application)
	router := routes.SetupRouter()

	// Message d'information dans le terminal
	log.Println("Demarrage du serveur sur :8080...")
	// Lancement du serveur HTTP et gestions d'erreurs éventuelles
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}