package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, pass, name, host, port,
	)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erreur lors de l'ouverture de la base des données:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Impossible de se connecter à la base des données", err)
	}

	log.Println("Connexion PostgreSQL OK →", name)
}
