package cron

import (
	"fmt"
	"safebase/db"
	"safebase/services"
	"safebase/utils"

	"github.com/robfig/cron/v3"
)

// StartCronScheduler démarre le scheduler et enregistre toutes les tâches actives.
func StartCronScheduler() *cron.Cron { // Retourne l'instance du cron pour les actions éventuelles
// Parser standard sans secondes (format 5 champs: minute, heure, jour, moi et jour de la semaine)
parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	c := cron.New(cron.WithParser(parser))

	rows, err := db.DB.Query(db.QueryGetAllActiveTasks)
	if err != nil {
		utils.LogError("Erreur lors du chargement des tâches CRON actives", err)
		return c
	}
	defer rows.Close()

	for rows.Next() {
		var id, userID, databaseID int
		var expr string

		if err := rows.Scan(&id, &userID, &databaseID, &expr); err != nil {
			utils.LogError("Erreur lors du scan d'une tâche CRON", err)
			continue
		}

		// Capture des variables locales dans un scope séparé pour éviter la closure problematique 
		taskID := id
		uid := userID
		dbID := databaseID
		schedule := expr

		// Fabrication d'une closure qui capture taskID, dbID, uID
		job := func(tid, did, uid int) func() {
			return func() {
				utils.LogInfo(fmt.Sprintf("CRON: exécution tâche %d (db = %d user = %d)", tid, did, uid))
				services.ExecuteScheduledBackup(tid)
			}
		}(taskID, dbID, uid)

		// Ajout de la tâche au scheduler
		if _, err :=  c.AddFunc(schedule, job); err != nil {
			utils.LogError(fmt.Sprintf("Erreur lors de l'enregistrement de la tâche CRON id=%d, expr=%s",taskID, schedule), err)
			continue
		}

		utils.LogInfo(fmt.Sprintf("Tâche CRON enregistée: id=%d, expr:%s", taskID, schedule))
	}
	// Démarrer le scheduler en arrière-plan
	c.Start()
	utils.LogInfo("Scheduler CRON démarré")
	return c
}