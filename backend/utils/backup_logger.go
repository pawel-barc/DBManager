package utils

import (
	"fmt"
	"os"
	"time"
)

// Cette fonction crée un fichier log pour un backup donné, enregistre le contenu fourni et retourne le chemin du fichier
func WriteBackupLog(backupID int, content string) (string, error) {
	os.MkdirAll("logs/backups", 0755)

	filename := fmt.Sprintf("logs/backups/backup_%d_%s.log",
	backupID, time.Now().Format("20060102_150405"),)

	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
	return "", err
	}
	return filename, nil
}