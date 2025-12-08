package utils

import (
	"safebase/db"
	"time"
)

// Fonction exécutée par le CRON, elle parcourt toutes les bases de données et vérifie la date du dernier backup
func CheckOldBackups() {

	rows, err := db.DB.Query(db.QueryGetDatabasesLastBackup)
	if err != nil {
		LogError("Imopossible de récupérer les dernières sauvegardes", err)
		return
	}

	for rows.Next() {
		var userID int
		var dbName string
		var lastBackup *time.Time
		if err := rows.Scan(&userID, &dbName, &lastBackup); err != nil {
			LogError("Erreur scan des backups", err)
			continue
		}
		// Si aucun backup n'a été réalisé depuis plus de 7 jours, une alerte est automatiquement envoyeé à l'utilisateur 
		if lastBackup == nil || time.Since(*lastBackup) > 7*24*time.Hour {
			NoBackupWarning(userID, dbName)
		}
	}
}