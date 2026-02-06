package services

import (
	"fmt"
	"os"
	"os/exec"
	"safebase/db"
	"safebase/utils"
)

// RestoreBackup exécute de manière synchrone une restauration à partir d'un fichier de backup
func RestoreBackup(backupID, userID int) (restoreID int, err error) {
	// Insérer une entrée dans la table restores avec le statut initial 'pending'
	err = db.DB.QueryRow(db.QueryInsertRestore, backupID, userID).Scan(&restoreID)
	if err != nil {
		utils.LogError("Impossible d'enregistrer la restauration", err)
		return 0, fmt.Errorf("impossible d'enregistrer la restauration: %v", err)
	}

	// Récupérer les informations du backup et les informations de connexion à la base
	var databaseID int
	var dbType, host, port, username, password, dbName, backupPath string
	
	err = db.DB.QueryRow(db.QueryGetBackupByID, backupID, userID).Scan(
		&backupID, &databaseID, &backupPath, &dbType, &host, &port, &username, &password, &dbName)

	if err != nil {
		msg := fmt.Sprintf("Backup introuvable ou accès refusé: %v", err)
		db.DB.Exec(db.QueryUpdateRestoreError, msg, restoreID)
		utils.LogError(msg, err)
		return restoreID, fmt.Errorf(msg)
	}

	// Vérifier que le fichier de backup existe
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		msg := fmt.Sprintf("Fichier de backup introuvable: %s", backupPath)
		db.DB.Exec(db.QueryUpdateRestoreError, msg, restoreID)
		utils.LogError(msg, err)
		return restoreID, fmt.Errorf(msg)
	}

	// Exécuter la restauration selon le type de base de données
	var restoreErr error
	switch dbType {
	case "postgres":
		restoreErr = runRestorePostgres(restoreID, dbName, host, port, username, password, backupPath)
	case "mysql":
		restoreErr = runRestoreMySQL(restoreID, dbName, host, port, username, password, backupPath)
	default:
		errMsg := "Type de base non supporté: " + dbType
		db.DB.Exec(db.QueryUpdateRestoreError, errMsg, restoreID)
		utils.LogError(errMsg, nil)
		return restoreID, fmt.Errorf(errMsg)
	}

	if restoreErr != nil {
		return restoreID, restoreErr
	}

	utils.LogInfo(fmt.Sprintf("Restore backup #%d terminé avec succès (restore_id: %d)", backupID, restoreID))
	return restoreID, nil
}

// runRestorePostgres effectue la restauration PostgreSQL via psql
func runRestorePostgres(restoreID int, dbName, host, port, username, password, backupPath string) error {
	// Définir le mot de passe via variable d'environnement
	os.Setenv("PGPASSWORD", password)
	defer os.Unsetenv("PGPASSWORD")

	// Construction de la commande psql pour restaurer depuis un fichier SQL
	// L'option -f permet de spécifier le fichier à exécuter
	cmd := exec.Command(
		"psql",
		"-h", host,
		"-p", port,
		"-U", username,
		"-d", dbName,
		"-f", backupPath, // Exécuter le fichier SQL
		"--quiet",        // Réduire la verbosité
	)

	// Exécuter la commande et récupérer la sortie
	output, err := cmd.CombinedOutput()
	if err != nil {
		logMsg := fmt.Sprintf("Erreur psql restore: %v\nOutput: %s", err, string(output))
		db.DB.Exec(db.QueryUpdateRestoreError, logMsg, restoreID)
		utils.LogError(logMsg, err)
		return err
	}

	// Mise à jour du statut de restauration en "success"
	logMsg := fmt.Sprintf("Restauration PostgreSQL réussie\nOutput: %s", string(output))
	db.DB.Exec(db.QueryUpdateRestoreSuccess, logMsg, restoreID)
	utils.LogInfo(fmt.Sprintf("Restauration PostgreSQL réussie (restore_id: %d)", restoreID))
	
	return nil
}

// runRestoreMySQL effectue la restauration MySQL en redirigeant le fichier SQL
func runRestoreMySQL(restoreID int, dbName, host, port, username, password, backupPath string) error {
	// Ouvrir le fichier de backup
	file, err := os.Open(backupPath)
	if err != nil {
		logMsg := fmt.Sprintf("Impossible d'ouvrir le fichier de backup: %v", err)
		db.DB.Exec(db.QueryUpdateRestoreError, logMsg, restoreID)
		utils.LogError(logMsg, err)
		return err
	}
	defer file.Close()

	// Construction de la commande mysql
	cmd := exec.Command(
		"mysql",
		"-h", host,
		"-P", port,
		"-u", username,
		fmt.Sprintf("-p%s", password),
		dbName,
	)

	// Rediriger le fichier SQL vers stdin de la commande mysql
	cmd.Stdin = file

	// Exécuter la commande et récupérer la sortie
	output, err := cmd.CombinedOutput()
	if err != nil {
		logMsg := fmt.Sprintf("Erreur mysql restore: %v\nOutput: %s", err, string(output))
		db.DB.Exec(db.QueryUpdateRestoreError, logMsg, restoreID)
		utils.LogError(logMsg, err)
		return err
	}

	// Mise à jour du statut de restauration en "success"
	logMsg := fmt.Sprintf("Restauration MySQL réussie\nOutput: %s", string(output))
	db.DB.Exec(db.QueryUpdateRestoreSuccess, logMsg, restoreID)
	utils.LogInfo(fmt.Sprintf("Restauration MySQL réussie (restore_id: %d)", restoreID))
	
	return nil
}