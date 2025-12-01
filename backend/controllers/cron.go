package controllers

import (
	"encoding/json"
	"net/http"
	"safebase/db"
	"safebase/middleware"
	"safebase/models"
	"safebase/utils"
)

// Créer une nouvelle tâche planifiée
func CreateScheduledTask (w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Structure permettant de récupérer les données de la requête
	var req struct {
		DatabaseID int `json:"database_id"`
		CronExpression string `json:"cron_expression"`
	}
	// Lecture et décodage du JSON envoyé par le client
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Format JSON invalide")
		return
	}
	// Variable pour stocker l'ID de la tâche nouvellement créée
	var taskID int

	// Insertion de la tâche dans la base de données
	err := db.DB.QueryRow(db.QueryCreateScheduledTask, userID, req.DatabaseID, req.CronExpression).Scan(&taskID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Impossible de créer une tâche")
		return
	}

	utils.SendSuccess(w, http.StatusCreated, "Tâche programmée créée")

}
// Récupérer toutes les tâches planifiées d'un utilisateur
func GetScheduledTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Exécution de la requête qui return toutes les tâches liées à cet utilisateur
	rows, err := db.DB.Query(db.QueryGetUserScheduledTasks, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Erreur lors du chargement")
		return
	}

	defer rows.Close()

	// List dans laquelle seront stockées les tâches récupérées 
	tasks := []models.ScheduledTask{}
	for rows.Next() {
		var t models.ScheduledTask
		// Lecture d'une ligne de résultat
		err := rows.Scan(&t.ID, &t.UserID, &t.DatabaseID, &t.CronExpression, &t.IsActive, &t.LastRunAt)
		if err == nil {
			tasks = append(tasks, t)
		}
	}
// Réponse contenant toutes les tâches de l'utilisateur
utils.SendSuccessWithData(w, http.StatusOK, "Récupération réusie", tasks)
}

// Mettre à jour l'expression CRON d'une tâche
func UpdateCronExpression(w http.ResponseWriter, r *http.Request) {
	// Structure pour stocker les données envoyées par le client
	var req struct {
		CronExpression string `json:"cron_expression"`
		TaskID int `json:"id"`
	}

	// Décodage du JSON reçu
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Format JSON invalide")
		return
	}

	// Exécution de la mise à jour dans la base de données
	_, err := db.DB.Exec(db.QueryUpdateCronExpression, req.CronExpression, req.TaskID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Erreur lors de la mise à jour de la tâche ")
		return
	}
utils.LogInfo("La tâche à été mise à jour")

}

// Activer ou désactiver une tâche planifiée
func ToggleTaskActive(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TaskID int `json:"task_id"`
		Active bool `json:"active"`
	}

	// Décodage du JSON reçu
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Format JSON invalide")
		return
	}

	// Mise à jour de l'état actif/inactif
	_, err := db.DB.Exec(db.QueryToggleTaskActive, req.Active, req.TaskID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Erreur lors de la mise à jour de la tâche")
		return
	}
	utils.LogInfo("L'état' à été mise à jour")
}