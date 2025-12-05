package controllers

import (
	"net/http"
	"safebase/middleware"
	"safebase/services"
	"safebase/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type RestoreBackupResponse struct {
	RestoreID int `json:"restore_id"`
	Message string `json:"message"`
}

func RestoreBackup(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	vars := chi.URLParam(r, "backup_id")
	backupID, err := strconv.Atoi(vars)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ID backup invalide")
		utils.LogError( "ID backup invalide", err)
		return
	}

	restoreID, dbName, err := services.RestoreBackup(backupID, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Backup introuvable ou accès refusé ")
		utils.LogError( "Backup introuvable ou accès refusé", err)
		utils.RestoreError(userID, backupID, err.Error())
		return
	}
	utils.RestoreSuccess(userID, backupID, dbName)
	utils.SendSuccessWithData(w, http.StatusCreated, "Restauration créée avec succès", RestoreBackupResponse {
		RestoreID: restoreID,
		Message: "Restauration créée avec succès",
	})
}