package services

import (
	"fmt"
	"os"
	"os/exec"
	"safebase/db"
	"safebase/utils"
)

// RestoreBackup exécute de maniere synchrone une restauration à partir d'un fichier de backup
func RestoreBackup(backupID, userID int) (restoreID int, err error) {
	// Insérer une entrée dans la table restores avec le statut initial
	err = db.DB.QueryRow(db.QueryInsertRestore, backupID, userID).Scan(&restoreID)
	if err != nil {
		utils.LogError("Impossible d'enregistrer la restauration:", err )
		return 0, fmt.Errorf("Impossible d'enregistrer la restauration: %v", err)

	}

	// Récupérer les informations du backup et les informations de connexion à la base
	var databaseID int
	var dbType, host, port, username, password, dbName, backupPath string

	err = db.DB.QueryRow(db.QueryGetBackupByID, backupID, userID).Scan(
		&backupID, &databaseID, &backupPath, &dbType, &host, &port, &username, &password, &dbName)
		
	if err != nil {
		msg := fmt.Sprintf("Backup introuvable ou accès refusé: %v", err)
		db.DB.Exec(db.QueryUpdateRestoreError, msg, restoreID)
		return restoreID, err
	}

	// Construire la commande système en fonction du type de base de données
	var cmd *exec.Cmd
	switch dbType {
	case "postgres":
		// Définir le mot de passe via variable d'environnement pour pg_restore
		os.Setenv("PGPASSWORD", password)
		defer os.Unsetenv("PGPASSWORD")

		// Construction de la commande pg_restore 
		cmd = exec.Command(
			"psql",
			"-h", host,
			"-p", port,
			"-U", username,
			"-d", dbName, // base ciblé
			backupPath, // Destination du fichier du backup
		)
	case "mysql":
		// Construction de la commande pour mysql
		file, err := os.Open(backupPath)
		if err != nil {
			db.DB.Exec(db.QueryUpdateRestoreError, err.Error(), restoreID)
			utils.LogError("Impossible d'ouvrir le fichier de backup mysql", err)
			return restoreID, err
		}
		defer file.Close()
		cmd = exec.Command(
			"mysql",
			"-h", host, // Adresse du serveur MySQL
			"-P", port, // Port de connexion
			"-u", username,
			fmt.Sprintf("-p%s", password), // Mot de passe(format imposé par mysqldump)
			dbName, // Nom de la base à sauvegarder
		)
		cmd.Stdin = file

	default:
		// Type de base non supporté => mise à jour du status en erreur et return
		errMsg := "Type de base non supporté"
		db.DB.Exec(db.QueryUpdateRestoreError, errMsg, restoreID)
		utils.LogError(errMsg, nil)
		return restoreID, fmt.Errorf(errMsg)
	}
	
	// Exécuter la commande de restauration et récupérer la sortie
	output, err := cmd.CombinedOutput()
	if err != nil {
		// En cas d'erreur mise à jour du statut => error
		logMsg := fmt.Sprintf("Erreur pg_restore/mysql: %v, output: %s", err, string(output))
		db.DB.Exec(db.QueryUpdateRestoreError, logMsg, restoreID)
		utils.LogError(logMsg, err)
		return restoreID, err
	}

	// Mise à jour du statut de restauration en "success"
	db.DB.Exec(db.QueryUpdateRestoreSuccess, string(output), restoreID)
	utils.LogInfo(fmt.Sprintf("Restore backup #%d terminé avec succès", backupID))
	return restoreID, nil
}
