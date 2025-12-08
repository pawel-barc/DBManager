package cron

import (
	"fmt"
	"safebase/utils"

	"github.com/robfig/cron/v3"
)

// Vérifie chaque jour à 03:00 si certaines bases n'ont pas été sauvegardées
func StartBackupHealthCron() {
	c := cron.New()
	_, err := c.AddFunc("59 19 * * *", func() {
		fmt.Println("CRON - Vérification des backups > 7 jours...")
		utils.CheckOldBackups()
	})

	if err != nil {
		utils.LogError("Erreur CRON backup_health", err)
	}

	c.Start()
	utils.LogInfo("CRON des vérifications de backups démarré")
}