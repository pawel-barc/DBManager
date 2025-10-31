package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"safebase/db"
	"safebase/middleware"
	"safebase/utils"
)

// --- Structures d'entrée/sortie ---

type CreateBackupRequest struct {
	DatabaseID int    `json:"database_id"`
	Name       string `json:"name"`
	Version    string `json:"version,omitempty"`
}

type CreateBackupResponse struct {
	BackupID int    `json:"backup_id"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

type BackupRow struct {
	ID         int        `json:"id"`
	DatabaseID int        `json:"database_id"`
	Name       string     `json:"name"`
	FilePath   string     `json:"file_path"`
	FileSize   float64    `json:"file_size"`
	BackupDate *time.Time `json:"backup_date,omitempty"`
	Status     string     `json:"status"`
	Version    string     `json:"version"`
	Log        string     `json:"log"`
}

type BackupStatusResponse struct {
	ID       int     `json:"id"`
	Status   string  `json:"status"`
	Progress *int    `json:"progress,omitempty"`
	Error    *string `json:"error,omitempty"`
}

type BackupLogResponse struct {
	ID  int    `json:"id"`
	Log string `json:"log"`
}

// --- Vérifie que l'utilisateur possède la base associée ---

func assertDatabaseOwnership(userID int, databaseID int) error {
	var count int
	err := db.DB.QueryRow(db.QueryCheckDatabaseOwnership, databaseID, userID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// --- Création d'un backup (POST /api/backups) ---

func CreateBackup(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var req CreateBackupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Format JSON invalide")
		return
	}
	if req.DatabaseID == 0 || req.Name == "" {
		utils.SendError(w, http.StatusBadRequest, "database_id et name sont obligatoires")
		return
	}

	if err := assertDatabaseOwnership(userID, req.DatabaseID); err != nil {
		if err == sql.ErrNoRows {
			utils.SendError(w, http.StatusForbidden, "Accès refusé à cette base")
			return
		}
		utils.SendError(w, http.StatusInternalServerError, "Vérification d'accès impossible")
		return
	}

	var backupID int
	err := db.DB.QueryRow(db.QueryInsertBackup, req.DatabaseID, req.Name, req.Version).Scan(&backupID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Impossible de créer le backup")
		return
	}

	utils.SendJSON(w, http.StatusCreated, CreateBackupResponse{
		BackupID: backupID,
		Status:   "queued",
		Message:  "Backup créé et en file d'attente",
	})
}

// --- Liste des backups (GET /api/backups?database_id=) ---

func ListBackups(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	dbIDStr := r.URL.Query().Get("database_id")
	if dbIDStr == "" {
		utils.SendError(w, http.StatusBadRequest, "Paramètre database_id requis")
		return
	}
	databaseID, err := strconv.Atoi(dbIDStr)
	if err != nil || databaseID <= 0 {
		utils.SendError(w, http.StatusBadRequest, "database_id invalide")
		return
	}

	if err := assertDatabaseOwnership(userID, databaseID); err != nil {
		if err == sql.ErrNoRows {
			utils.SendError(w, http.StatusForbidden, "Accès refusé à cette base")
			return
		}
		utils.SendError(w, http.StatusInternalServerError, "Vérification d'accès impossible")
		return
	}

	rows, err := db.DB.Query(db.QueryListBackupsByDatabase, databaseID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Impossible de charger les backups")
		return
	}
	defer rows.Close()

	var list []BackupRow
	for rows.Next() {
		var b BackupRow
		if err := rows.Scan(
			&b.ID, &b.DatabaseID, &b.Name, &b.FilePath, &b.FileSize, &b.BackupDate,
			&b.Status, &b.Version, &b.Log,
		); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Erreur de lecture des backups")
			return
		}
		list = append(list, b)
	}

	utils.SendJSON(w, http.StatusOK, list)
}

// --- Statut d'un backup (GET /api/backups/{id}/status) ---

func GetBackupStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	vars := mux.Vars(r)
	idStr := vars["id"]
	backupID, err := strconv.Atoi(idStr)
	if err != nil || backupID <= 0 {
		utils.SendError(w, http.StatusBadRequest, "id invalide")
		return
	}

	var id int
	var status string
	err = db.DB.QueryRow(db.QueryGetBackupStatusWithOwnership, backupID, userID).Scan(&id, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendError(w, http.StatusNotFound, "Backup introuvable")
			return
		}
		utils.SendError(w, http.StatusInternalServerError, "Impossible de récupérer le statut")
		return
	}

	utils.SendJSON(w, http.StatusOK, BackupStatusResponse{
		ID:     id,
		Status: status,
	})
}

// --- Log d'un backup (GET /api/backups/{id}/log) ---

func GetBackupLog(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	vars := mux.Vars(r)
	idStr := vars["id"]
	backupID, err := strconv.Atoi(idStr)
	if err != nil || backupID <= 0 {
		utils.SendError(w, http.StatusBadRequest, "id invalide")
		return
	}

	var id int
	var logText string
	err = db.DB.QueryRow(db.QueryGetBackupLogWithOwnership, backupID, userID).Scan(&id, &logText)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendError(w, http.StatusNotFound, "Backup introuvable")
			return
		}
		utils.SendError(w, http.StatusInternalServerError, "Impossible de récupérer le log")
		return
	}

	utils.SendJSON(w, http.StatusOK, BackupLogResponse{
		ID:  id,
		Log: logText,
	})
}

// --- Téléchargement du fichier (GET /api/backups/{id}/download) ---

func DownloadBackup(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	vars := mux.Vars(r)
	idStr := vars["id"]
	backupID, err := strconv.Atoi(idStr)
	if err != nil || backupID <= 0 {
		utils.SendError(w, http.StatusBadRequest, "id invalide")
		return
	}

	var filePath, name string
	if err := db.DB.QueryRow(db.QueryGetBackupPathWithOwnership, backupID, userID).Scan(&filePath, &name); err != nil {
		if err == sql.ErrNoRows {
			utils.SendError(w, http.StatusNotFound, "Fichier introuvable ou backup non terminé")
			return
		}
		utils.SendError(w, http.StatusInternalServerError, "Impossible de récupérer le fichier")
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Erreur ouverture du fichier")
		return
	}
	defer f.Close()

	filename := name
	if ext := filepath.Ext(filePath); ext != "" && filepath.Ext(filename) == "" {
		filename = filename + ext
	}

	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(w, r, filename, time.Now(), f)
}
