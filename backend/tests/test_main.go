package tests

import (
	"log"
	"os"
	"testing"

	"safebase/db"
)

// DDL complet pour initialiser le schéma de test
const schemaSQL = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    refresh_token TEXT
);

CREATE TABLE IF NOT EXISTS databases (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
    host VARCHAR(100) NOT NULL,
    port INT NOT NULL,
    db_username VARCHAR(50) NOT NULL,
    db_password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS backups (
    id SERIAL PRIMARY KEY,
    database_id INT NOT NULL REFERENCES databases(id),
    name VARCHAR(100),
    file_path TEXT NOT NULL,
    file_size FLOAT,
    backup_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    version VARCHAR(20),
    log TEXT
);

CREATE TABLE IF NOT EXISTS restores (
    id SERIAL PRIMARY KEY,
    backup_id INT NOT NULL REFERENCES backups(id),
    user_id INT NOT NULL REFERENCES users(id),
    restored_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    log TEXT 
);

CREATE TABLE IF NOT EXISTS scheduled_tasks (
    id SERIAL PRIMARY KEY,
    database_id INT NOT NULL REFERENCES databases(id),
    cron_expression VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    last_run_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS alerts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    alert_type VARCHAR(50),
    message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_read BOOLEAN DEFAULT FALSE
);
`

func TestMain(m *testing.M) {
	// 1️⃣ Connexion à la DB de test
	db.ConnectDB()
	if db.DB == nil {
		log.Fatal("❌ db.DB est nil après ConnectDB – vérifie la chaîne de connexion")
	}

	// 2️⃣ Création du schéma complet (users, databases, backups, restores, scheduled_tasks, alerts)
	if _, err := db.DB.Exec(schemaSQL); err != nil {
		log.Fatalf("❌ Impossible de créer le schéma de test: %v", err)
	}

	// 3️⃣ Nettoyage des données avant les tests
	if _, err := db.DB.Exec(`TRUNCATE alerts, restores, backups, scheduled_tasks, databases, users RESTART IDENTITY CASCADE;`); err != nil {
		log.Printf("⚠️ Impossible de nettoyer les tables avant les tests: %v", err)
	}

	// 4️⃣ Lancer tous les tests
	code := m.Run()

	// 5️⃣ Nettoyage après les tests
	if _, err := db.DB.Exec(`TRUNCATE alerts, restores, backups, scheduled_tasks, databases, users RESTART IDENTITY CASCADE;`); err != nil {
		log.Printf("⚠️ Impossible de nettoyer les tables après les tests: %v", err)
	}

	// 6️⃣ Fermeture de la connexion
	_ = db.DB.Close()

	os.Exit(code)
}
