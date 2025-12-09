package controllers

import (
	"net/http"
	"safebase/db"
	"safebase/middleware"
	"safebase/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
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

// Marque un alerte comme lu
func MarkAlertAsRead(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	alertIDstr := chi.URLParam(r, "alert_id")
	alertID, err := strconv.Atoi(alertIDstr)
	if err != nil || alertID <= 0 {
		utils.SendError(w, http.StatusBadRequest, "ID d'alerte invalide")
		return
	}

	// Mettre à jour le statut dans la base
	_, err = db.DB.Exec(`UPDATE alerts SET is_read = TRUE WHERE id=$1 AND user_id=$2`, alertID, userID)
	if err != nil {
		utils.LogError("Impossible de mettre à jour l'alerte comme lue", err)
		utils.SendError(w, http.StatusInternalServerError, "Impossible de mettre à jour l'alerte")
		return
	}
	utils.SendSuccess(w, http.StatusOK, "Alerte marquée comme lue")

}


// Marquer toutes les alertes de l'utilisateur comme lues
func MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserIDKey).(int)
	
	_, err := db.DB.Exec(`UPDATE alerts SET is_read = TRUE WHERE user_id=$1`, userId)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Impossible de marquer toutes les alertes comme lues")
		return
	}
	utils.SendSuccess(w, http.StatusOK, "Toutes les alertes ont été marqueés comme lues")
}