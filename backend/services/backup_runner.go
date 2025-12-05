package services

import (
	"fmt"
	"os"
	"os/exec"

	"safebase/db"
	"safebase/utils"
)

// Exécute un backup PostgreSQL en utilisant pg_dump et met à jour le statut dans la base
func RunBackupPostgres(backupID int, dbName, host, port, username, password, backupPath string) error {
	// Définition du mot de passe via la variable d'environnement PGPASSWORD et nettoyage après l'exécution 
	os.Setenv("PGPASSWORD", password)
	defer os.Unsetenv("PGPASSWORD")

	// Construction de la commande pg_dump 
	cmd := exec.Command(
		"pg_dump",
		"-h", host,
		"-p", port,
		"-U", username,
		"-F", "p", // Format sql
		"-f", backupPath, // Destination du fichier du backup
		dbName, // Nom de la base à sauvegarder
	)

	// Exécution de la commande et récupération de stdout et stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Message d'erreur détaillé à enregistrer
		logMsg := fmt.Sprintf("pg_dump failed: %s\nOutput: %s", err.Error(), string(output))
		utils.LogError(logMsg, err)

		// Mise à jour du backup en échec dans la base de données (status = 'error')
		db.DB.Exec(db.QueryUpdateBackupFail, logMsg, backupID)
		return err
	}

	// Récupération de la taille du fichier généré en KB
	info, _ := os.Stat(backupPath)
	fileSize := float64(info.Size()) / 1024

	// Mise à jour du backup en succès dans la base de données (status = 'success')
	db.DB.Exec(db.QueryUpdateBackupSuccess, backupPath, fileSize, "Backup réussi", backupID)
	utils.LogInfo(fmt.Sprintf("Backup postgres réussi: %s, size=%.2f KB", backupPath, fileSize))
	return nil
}

// Exécute un backup  MySQL via mysqldump, redirige la sortie dans un fichier et met à jour le statut du backup dans la DB
func RunBackupMySQL(backupID int, dbName, host, port, username, password, backupPath string) error {
	// Construction de la commande pg_dump 
	cmd := exec.Command(
		"mysqldump",
		"-h", host, // Adresse du serveur MySQL
		"-P", port, // Port de connexion
		"-u", username,
		fmt.Sprintf("-p%s", password), // Mot de passe(format imposé par mysqldump)
		dbName, // Nom de la base à sauvegarder
	)

	// Création du fichier de sortie qui recevra le contenu du dump
	outfile, _ := os.Create(backupPath)
	defer outfile.Close()

	// Redirection de la sortie standard de mysqldump vers le fichier
	cmd.Stdout = outfile

	// Exécution de la commande systeme
	if err := cmd.Run(); err != nil {
		utils.LogError("mysqldump failed", err)
		db.DB.Exec(db.QueryUpdateBackupFail, err.Error(), backupID)
		return err

	}

	// Récupération de la taille du fichier généré
	info, _ := os.Stat(backupPath)
	fileSize := float64(info.Size()) / 1024

	// Mise à jour du statut en base : succès + infos sur le fichier
	db.DB.Exec(db.QueryUpdateBackupSuccess, backupPath, fileSize, "Backup MySQL réussi", backupID)
	// Log d'information
	utils.LogInfo(fmt.Sprintf("Backup mysql réussi: %s, size=%.2f KB", backupPath, fileSize))
	return nil
}
