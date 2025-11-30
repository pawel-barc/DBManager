package cron

import (
	"fmt"
	"safebase/controllers"
	"safebase/utils"
	"time"

	"github.com/robfig/cron/v3"
)

// Fonction démarre le CRON pour exécuter automatiquement les backups
func StartBackupCron() {
	c := cron.New(cron.WithSeconds())
fmt.Println(time.Now())
	_, err := c.AddFunc("0 20 13 * * *", func()  {
		utils.LogInfo("CRON - Je commence le backup automatique pour toutes les bases")
		// Appel de la fonction qui exécute les backups automatique pour toutes les bases
		if err := controllers.RunAllBackups(); err != nil {
			utils.LogError("Erreur lors de l'execution du backup", err)
			return
		} else {
			utils.LogInfo("CRON - Backup crée avec succès")
		}
	})
	if err != nil {
		utils.LogError("Impossible d'ajouter la fonction au CRON", err)
	}
	c.Start()
	utils.LogInfo("CRON - Le programme des backups démarré")

}
