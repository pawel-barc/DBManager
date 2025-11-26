package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"safebase/db"
	"safebase/middleware"
	"safebase/utils"

	"github.com/go-chi/chi/v5"
)

// Request / response structure
type CreateBackupRequest struct {
	DatabaseID int `json:"database_id"`
	Name string `json:"name"`
	Version string `json:"version,omitempty"`
}

type CreateBackupResponse struct {
	BackupID int `json:"backup_id"`
	Status string `json:"status"`
	Message string `json:"message"`
}

// Helper vérifie si l'utilisateur est le propriétaire de la base
func assertDatabaseOwnership(userID, databaseID int) error {
	var count int
	// Exécuter la requête SQL pour vérifier la possession de la base
	err := db.DB.QueryRow(db.QueryCheckDatabaseOwnership, databaseID, userID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Exécuter une sauvegarde
func CreateBackup(w http.ResponseWriter, r *http.Request) {
	// Récupération de l'ID utilisateur depuis le contexte
	userID := r.Context().Value(middleware.UserIDKey).(int)
	var req CreateBackupRequest
	// Décodage du JSON reçu
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Format JSON invalide")
		return
	}

	// Vérification des champs obligatoires
	if req.DatabaseID == 0 || req.Name == "" {
		utils.SendError(w, http.StatusBadRequest, "database_id et le name sont obligatoires")
		return
	}

	// Vérification de la possession de la base
	if err := assertDatabaseOwnership(userID, req.DatabaseID); err != nil {
		if err == sql.ErrNoRows {
			utils.SendError(w, http.StatusForbidden, "Accès refusé à cette base")
			return
		}

		utils.SendError(w, http.StatusInternalServerError, "Vérification d'accès impossible")
		return
	}

	// Création du dossier de backup s'il n'existe pas
	backupDir := "./backups"
	if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Impossible de créer le dossier de backup")
		return
	}

	// Création du nom de fichier unique avec date et heure
	filename := req.Name + "_" + time.Now().Format("20060102_150405") + ".sql"
	filePath := filepath.Join(backupDir, filename)

	// Création du fichier de backup
	f, err := os.Create(filePath)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Impossible de créeer le fichier backup")
		return
	}
	defer f.Close() // Fermeture automatique du fichier à la fin de la fonction

	// Insertion du backup dans la base de récupération de l'ID
	var backupID int
	err = db.DB.QueryRow(db.QueryInsertBackup, req.DatabaseID, req.Name, req.Version, filePath).Scan(&backupID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Impossible de créer le backup")
		return
	}

	// Réponse réussie avec l'ID du backup
	utils.SendSuccessWithData(w, http.StatusCreated, "Backup crée avec succès", CreateBackupResponse {
		BackupID: backupID,
		Status: "completed",
		Message: "Backup créé et enregistré",
	})
}