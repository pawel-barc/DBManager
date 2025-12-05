package utils

import (
	"fmt"
	"safebase/db"
)

// Types des alerts possibles

const (
	AlertBackupSuccess =  "backup_success"
	AlertBackupError = "backup_error"
	AlertRestoreSuccess = "restore_success"
	AlertRestoreError = "restore_error"
	AlertNotBackupWarning = "no_backup_warning"
)

// Enregistre une alerte de succès lors d'un backup terminé
func BackupSuccess(userID, backupID int, dbName string) {
	message := fmt.Sprintf("Backup #%d pour '%s' a été réalisé avec succès", backupID, dbName)
	insertAlert(userID, AlertBackupSuccess, message)
}
// Enregistre une alerte d'échec lors d'un backup
func BackupError(userID, backupID int, errMsg string) {
	message := fmt.Sprintf("Echec du backup #%d : %s", backupID, errMsg)
	insertAlert(userID, AlertBackupError, message)
}
// Enregistre une alerte de succès lors d'une restauration
func RestoreSuccess(userID, backupID int, dbName string) {
	message := fmt.Sprintf("Restauration du backup: #%d pour '%s' a été réalisé avec succès", backupID, dbName)
	insertAlert(userID, AlertRestoreSuccess, message)
}
// Enregistre une alerte d'échec lors d'une restauration
func RestoreError(userID, backupID int, errMsg string) {
	message := fmt.Sprintf("Echec de la restauration: #%d : %s", backupID, errMsg)
	insertAlert(userID, AlertRestoreError, message)
}
// Alerte lorsqu'aucun backup n'a été réalisé depuis plus de sept jours
func NoBackupWarning(userID int, dbName string) {
	message := fmt.Sprintf("La base '%s' n'a pas eu de backup depuis plus de 7 jours", dbName)
	insertAlert(userID, AlertNotBackupWarning, message)
}
// Fonction interne générique qui insère une alerte dans la base
func insertAlert(userId int, alertType, message string) {
	db.DB.Exec(`INSERT INTO alerts (user_id, alert_type, message) VALUES ($1, $2, $3)`, userId, alertType, message )
}