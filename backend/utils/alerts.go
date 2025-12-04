package utils

import (
	"safebase/db"
)

func CreateAlert(userID int, message string) {
	db.DB.Exec(`INSERT INTO alerts (user_id, alert_type, message) VALUES ($1, 'backup_error', $2)`, userID, message) 
}