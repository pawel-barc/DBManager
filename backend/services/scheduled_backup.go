package services

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"safebase/db"
	"safebase/utils"
	"time"
)

// Lance un backup planifié identifié par taskID
func ExecuteScheduledBackup(taskID int) {
	var DatabaseID, userID int
	err := db.DB.QueryRow(`SELECT database_id, user_id FROM scheduled_tasks WHERE id = $1`, taskID).Scan(&DatabaseID, &userID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogError(fmt.Sprintf("CRON: tâche introuvable id =%d", taskID), nil)
			return
		}
		utils.LogError("CRON erreur, lors de la récupération d'une tâche", err)
		return
	}

	// Vérification locale de la possession de la base
	var ownerCount int
	err = db.DB.QueryRow(`SELECT COUNT(1) FROM databases WHERE id = $1 AND user_id = $2`, DatabaseID, userID).Scan(&ownerCount)
	if err != nil {
		utils.LogError("CRON: erreur lors de la verfication de la priopriété de la base", err)
		return
	}
	if ownerCount == 0 {
		utils.LogError(fmt.Sprintf("CRON: L'utilisateur %d n'est pas le priopriétaire de la database %d (tâche %d)", userID, DatabaseID, taskID), nil)
		return
	}

	// Construire un nom de backup type CRON
	backupName := fmt.Sprintf("scheduled_db_%d_task_%d", DatabaseID, taskID)
	version := "v1"

	// Appel de la fonction non-HTTP qui effectue le backup
	if err := CreateBackupForDb(DatabaseID, backupName, version); err != nil {
		utils.LogError(fmt.Sprintf("CRON:  échec du backup planifié, tâche=%d db=%d", taskID, DatabaseID), err)
		return
	}

	// Mise à jour du last_run_at après succès
	if _, err := db.DB.Exec(db.QueryUpdateLastTaskRun, taskID); err != nil {
		utils.LogError(fmt.Sprintf("CRON: Impossible de mettre à jour last_run_at pour task: %d", taskID), err)
		return
	}
	utils.LogInfo(fmt.Sprintf("CRON: tâche %d exécuté avec succès", taskID))
}
// Fonction 
func CreateBackupForDb(databaseID int, backupName, version string) error {
	var dbType, host, username, password, dbName string
	var port string
	err := db.DB.QueryRow(db.QuerySelectDatabaseInfoNoUser, databaseID).Scan(&dbType, &host, &port, &username, &password, &dbName)
	if err != nil {
		utils.LogError("Impossible de récupérer les infos de la database", err)
		return err
	}

	// Création du dossier de backup s'il n'existe pas
	backupDir := "./backups"
	if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
		utils.LogError("Impossible de créer le dossier de backup", err)
		return err
	}

	// Création du nom de fichier unique avec date et heure
	filename := backupName + "_" + time.Now().Format("20060102_150405") + ".sql"
	filePath := filepath.Join(backupDir, filename)

	// Insertion du backup dans la base de récupération de l'ID
	var backupID int
	err = db.DB.QueryRow(db.QueryInsertBackup, databaseID, backupName, version, filePath).Scan(&backupID)
	if err != nil {
		utils.LogError("Impossible de créer l'entrée backup dans la base", err)
		return err
	}
	utils.LogInfo(fmt.Sprintf("Backup #%d crée (PENDING)", backupID))

	// Exécution du backup réel

	var backupErr error

	switch dbType {
	case "postgres":
		backupErr = RunBackupPostgres(backupID, dbName, host, port, username, password, filePath)
	case "mysql":
		backupErr = RunBackupMySQL(backupID, dbName, host, port, username, password, filePath)
	default:
		utils.LogError("Type de base non supporté: " +dbType, nil)
		return err
	}
	if backupErr != nil {
		return err
	}
return nil
}