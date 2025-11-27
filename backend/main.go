package main

import (
	"log"
	"net/http"

	"safebase/config"
	"safebase/db"
	"safebase/routes"
)

func main() {
	// Connexiona à la base des données
	db.ConnectDB()
	// Configuration du routeur (définition des routes de l'application)

	// Migration: création de la base des données, Commentez après avoir utiliser
	// sqlBytes, err := os.ReadFile("db/migrations/001_init.sql")
	// if err != nil {
	// 	log.Fatal("Impossible de lire le fichier de migration:", err)
	// }

	// _, err = db.DB.Exec(string(sqlBytes))
	// if err != nil {
	// 	log.Fatal("Erreur lors de l'exécution de la migration", err)
	// }

	// log.Println("Migration exécutée avec succès")

	router := routes.SetupRouter()

	// Message d'information dans le terminal
	log.Println("Demarrage du serveur sur :8080...")
	// Lancement du serveur HTTP et gestions d'erreurs éventuelles
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
	config.InitLogger()
}