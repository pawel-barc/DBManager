package controllers

// import(
// 	"encoding/json"
// 	"net/http"
// 	"safebase/db"
// 	"safebase/middleware"
// 	"safebase/models"
// 	"safebase/utils"
// )

// // Ajouter une nouvelle tâche planifiée
// func CreateScheduledTask (w http.ResponseWriter, r *http.Request) {
// 	userID := r.Context().Value(middleware.UserIDKey).(int)

// 	var req struct {
// 		DatabaseID int `json:"database_is"`
// 		CronExpression string `json:"cron_expression"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.SendError(w, http.StatusBadRequest, "Format JSON invalide")
// 		return
// 	}

// 	var taskID int

// 	err := db.DB.QueryRow(db.QueryCreateScheduledTask, userID, req.DatabaseID, req.CronExpression).Scan(&taskID)
// 	if err != nil {
// 		utils.SendError(w, http.StatusInternalServerError, "Impossible de créer une tâche")
// 		return
// 	}

// 	utils.SendSuccess(w, http.StatusCreated, "Tâche programmée créée")

// }
