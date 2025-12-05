package controllers

import (
	"net/http"
	"safebase/db"
	"safebase/middleware"
	"safebase/utils"
)

// Récupère toutes les alertes pour l'utilisateur connecté
func GetUserAlerts(w http.ResponseWriter, r *http.Request) {
	// Récupération de l'ID depuis le middleware
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Exécution de la requête
	rows, err := db.DB.Query(db.QueryGetAllAlerts, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Impossible de récupérer les alerts")
		return
	}

	// Structure renvoyée au front
	type Alert struct {
		ID int `json:"id"`
		Type string `json:"alert_type"`
		Message string `json:"message"`
		CreatedAt string `json:"created_at"`
		IsRead bool `json:"is_read"`
	}

	var result []Alert

	// Parcours des lignes
	for rows.Next() {
		var a Alert
		rows.Scan(&a.ID, &a.Type, &a.Message, &a.CreatedAt, &a.IsRead)
		result = append(result, a)
	}
	// Réponse finale
	utils.SendSuccessWithData(w, http.StatusOK, "Liste des alertes récupérée",result)
}