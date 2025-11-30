package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"safebase/db"
	"safebase/middleware"
	"safebase/models"
	"safebase/services"
	"safebase/utils"

	"github.com/go-chi/chi/v5"
)

// Request / response structure
type CreateBackupRequest struct {
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
		utils.LogError("Format JSON invalide lors de la création du backup", err)
		utils.SendError(w, http.StatusBadRequest, "Format JSON invalide")
		return
	}

	vars := chi.URLParam(r, "database_id")
	databaseID, _ := strconv.Atoi(vars)

	// Vérification des champs obligatoires
	if databaseID == 0 || req.Name == "" {
		utils.LogError("Champs obligatoires manquants pour la création du backup", nil)
		utils.SendError(w, http.StatusBadRequest, "database_id et le name sont obligatoires")
		return
	}

	// Vérification de la possession de la base
	if err := assertDatabaseOwnership(userID, databaseID); err != nil {
		utils.LogError("Echec de la vérification de possession de la base", err)
		if err == sql.ErrNoRows {
			utils.SendError(w, http.StatusForbidden, "Accès refusé à cette base")
			return
		}

		utils.SendError(w, http.StatusInternalServerError, "Vérification d'accès impossible")
		return
	}

	// Récupération informations database
	var (
		dbType, host, username, password, dbName string
		port string
	)
	err := db.DB.QueryRow(db.QuerySelectDatabaseInfo, databaseID, userID).Scan(&dbType, &host, &port, &username, &password, &dbName)
	if err != nil {
		utils.LogError("Impossible de récupérer les infos de la database", err)
		utils.SendError(w, http.StatusInternalServerError, "Impossible de lire la base")
		return
	}

	// Création du dossier de backup s'il n'existe pas
	backupDir := "./backups"
	if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
		utils.LogError("Impossible de créer le dossier de backup", err)
		utils.SendError(w, http.StatusInternalServerError, "Impossible de créer le dossier de backup")
		return
	}

	// Création du nom de fichier unique avec date et heure
	filename := req.Name + "_" + time.Now().Format("20060102_150405") + ".sql"
	filePath := filepath.Join(backupDir, filename)

	// Insertion du backup dans la base de récupération de l'ID
	var backupID int
	err = db.DB.QueryRow(db.QueryInsertBackup, databaseID, req.Name, req.Version, filePath).Scan(&backupID)
	if err != nil {
		utils.LogError("Impossible de créer l'entrée backup dans la base", err)
		utils.SendError(w, http.StatusInternalServerError, "Impossible de créer le backup")
		return
	}
	utils.LogInfo(fmt.Sprintf("Backup #%d crée (PENDING)", backupID))

	// Exécution du backup réel

	var backupErr error

	switch dbType {
	case "postgres":
		backupErr = services.RunBackupPostgres(backupID, dbName, host, port, username, password, filePath)
	case "mysql":
		backupErr = services.RunBackupMySQL(backupID, dbName, host, port, username, password, filePath)
	default:
		utils.LogError("Type de base non supporté: " +dbType, nil)
		utils.SendError(w, http.StatusBadRequest, "Type de base non supporté")	
		return
	}
	if backupErr != nil {
		utils.SendError(w, http.StatusInternalServerError, "Echec du backup")
		return
	}

	// Réponse réussie avec l'ID du backup
	utils.SendSuccessWithData(w, http.StatusCreated, "Backup crée avec succès", CreateBackupResponse {
		BackupID: backupID,
		Status: "completed",
		Message: "Backup créé et enregistré",
	})
}

func CreateBackupForDb(databaseID int, backupName, version string) error {
	var dbType, host, username, password, dbName string
	var port string
	err := db.DB.QueryRow(db.QuerySelectDatabaseInfoNoUser, databaseID).Scan(&dbType, &host, &port, &username, &password, &dbName)
	if err != nil {
		utils.LogError("Impossible de récupérer les infos de la database", err)
		return err
	}

	// Création du dossier de backup s'il n'existe pas
	backupDir := "./backups"
	if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
		utils.LogError("Impossible de créer le dossier de backup", err)
		return err
	}

	// Création du nom de fichier unique avec date et heure
	filename := backupName + "_" + time.Now().Format("20060102_150405") + ".sql"
	filePath := filepath.Join(backupDir, filename)

	// Insertion du backup dans la base de récupération de l'ID
	var backupID int
	err = db.DB.QueryRow(db.QueryInsertBackup, databaseID, backupName, version, filePath).Scan(&backupID)
	if err != nil {
		utils.LogError("Impossible de créer l'entrée backup dans la base", err)
		return err
	}
	utils.LogInfo(fmt.Sprintf("Backup #%d crée (PENDING)", backupID))

	// Exécution du backup réel

	var backupErr error

	switch dbType {
	case "postgres":
		backupErr = services.RunBackupPostgres(backupID, dbName, host, port, username, password, filePath)
	case "mysql":
		backupErr = services.RunBackupMySQL(backupID, dbName, host, port, username, password, filePath)
	default:
		utils.LogError("Type de base non supporté: " +dbType, nil)
		return err
	}
	if backupErr != nil {
		return err
	}
return nil
}

// Récupération de tout les backups d'un utilisateur
func ListAllBackups(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	rows, err := db.DB.Query(db.QuerySelectAllBackups, userID)
	if err != nil {
		utils.LogError("Impossible de récupérer tout les backups", err)
		utils.SendError(w, http.StatusInternalServerError, "Erreur du chargement du backups")
		return
	}
	// Fermeture automatique du curseur de résultats
	defer rows.Close()

	var backups []models.Backup
	for rows.Next() {
		var b models.Backup
		var name, version, log sql.NullString
		var fileSize sql.NullFloat64
		if err := rows.Scan(&b.ID, &b.DatabaseID, &name, &b.FilePath, &fileSize, &b.BackupDate, &b.Status, &version, &log); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Erreur de lecture des backups")
			return
		}
		// En cas NULL
		if name.Valid {
			b.Name = name.String
		} else {
			b.Name = ""
		}
		if fileSize.Valid {
			b.FileSize = fileSize.Float64
		} else {
			b.FileSize = 0
		}
		if version.Valid {
			b.Version = version.String
		} else {
			b.Version = ""
		}
		if log.Valid {
			b.Log = log.String
		} else {
			b.Log = ""
		}
		backups = append(backups, b)
	}
	utils.SendSuccessWithData(w, http.StatusOK, "Liste complète des backups", backups)
}

// Récupération de la liste des backups d'une base de données
func ListBackups(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Récupération du paramètre database_id depuis l'URL
	dbIDStr := chi.URLParam(r, "database_id")

	if dbIDStr == "" {
		utils.SendError(w, http.StatusBadRequest, "Paramètre database_id requis")
		return
	}

	// Convertion du paramètre en entier
	databaseID, err := strconv.Atoi(dbIDStr)
	if err != nil || databaseID <= 0 {
		utils.SendError(w, http.StatusBadRequest, "database_id invalide")
		return
	}

	// Vérification que l'utilisateur est bien le propriétaire de la base 
	if err := assertDatabaseOwnership(userID, databaseID); err != nil {
		if err == sql.ErrNoRows {
			utils.SendError(w, http.StatusForbidden, "Accès refusé à cette base")
			return
		}
		utils.SendError(w, http.StatusInternalServerError, "Vérification d'accès impossible")
		return
	}

	// Exécution de la requête SQL pour récupérer les backups
	rows, err := db.DB.Query(db.QuerySelectBackupsByDatabase, databaseID )
	if err != nil {
		utils.LogError(fmt.Sprintf("Impossible de charger les backups pour database_id=%d", databaseID), err)
		utils.SendError(w, http.StatusInternalServerError, "Impossible de charger les backups")
		return
	}

	// Fermeture automatique du curseur de résultats
	defer rows.Close()

	// Lecture des résultats ligne par ligne
	var backups []models.Backup

	for rows.Next() {
		var b models.Backup
			var filePath sql.NullString
		// Remplissage de la structure Backup avec les colonnes SQL
		if err := rows.Scan(&b.ID, &b.DatabaseID, &b.Name, &filePath, &b.FileSize, &b.BackupDate, &b.Status, &b.Version, &b.Log); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Erreur de lecture des backups")
			return
		}
		if filePath.Valid {
			b.FilePath = filePath.String
		} else {
			b.FilePath = ""
		}
		backups = append(backups, b)
	}
	// Envoi de la réponse avec la liste complète
	utils.LogInfo(fmt.Sprintf("Liste des backups récupérée pour la database_id=%d, %d items", databaseID, len(backups)))
	utils.SendSuccessWithData(w, http.StatusOK, "Liste des backups", backups)

}
// Télécharger un backup 
func DownloadBackup(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	// Extraction de l'ID du backup depuis l'URL 
	backupIDStr := chi.URLParam(r, "backup_id")
	// Conversion en entier et validation
	backupID, err := strconv.Atoi(backupIDStr)
	if err != nil || backupID <= 0 {
		utils.SendError(w, http.StatusBadRequest, "id invalide")
		return
	}

	// Variables qui recevrons le chemin du fichier et le nom du backup 
	var filePath, name string
	// Vérifier que le backup appartient bien à l'utilisateur et récupérer file_path et name depuis la base
	err = db.DB.QueryRow(db.QueryGetBackupPathWithOwnership, backupID, userID).Scan(&filePath, &name)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendError(w, http.StatusNotFound, "Fichier introuvable ou backup non terminé")
			return
		}
		utils.LogError(fmt.Sprintf("Impossible de récupérer le chemin du backup_id=%d", backupID), err)
		utils.SendError(w, http.StatusInternalServerError, "Impossible de récupérer le fichier")
		return
	}

	// Ouverture du fichier physique sur le disque
	f, err := os.Open(filePath)
	if err != nil {
		utils.LogError(fmt.Sprintf("Impossible d'ouvrir le fichier backup_id=%d", backupID), err)
		utils.SendError(w, http.StatusInternalServerError, "Erreur ouverture du fichier")
		return
	}
	defer f.Close()

	// Configuration des en-têtes HTTP pour forcer le téléchargement 
	w.Header().Set("Content-Disposition", `attachment; filename="`+name+`"`)
	w.Header().Set("Content-Type", "application/octet-stream")
	// Envoi du fichier au client
	http.ServeContent(w, r, name, time.Now(), f)
	utils.LogInfo(fmt.Sprintf("Backup_id=%d téléchargé par user_id=%d, fichier=%s", backupID, userID, filePath))

}

// Suppression d'un backup 
func DeleteBackup(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Extraction de l'ID du backup depuis l'URL
	backupIDStr := chi.URLParam(r, "backup_id")
	backupID, err := strconv.Atoi(backupIDStr)
	if err != nil || backupID <= 0 {
		utils.SendError(w, http.StatusBadRequest, "Backup ID invalide")
		return
	}
	// Variable qui recevra le chemin du fichier 
	var filePath string
	// Vérification que le backup appartient bien à l'utilisateur
	err = db.DB.QueryRow(db.QueryFindBackupPathFile, backupID, userID).Scan(&filePath)
	if err != nil {
		if err == sql.ErrNoRows{
			utils.SendError(w, http.StatusNotFound, "Le backup n'existe pas ou vous n'y avez pas d'accès")
			return
		}
		utils.LogError("Erreur lors de la récupération du backup ", err)
		utils.SendError(w, http.StatusInternalServerError, "Erreur lors de la récupération du backup")
		return
	}

	// Suppression de l'enregistrement du backup dans la base de données
	_, err = db.DB.Exec(`DELETE FROM backups WHERE id = $1`, backupID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Impossible de supprimer ce backup")
		return
	}

	// Suppression du fichier physique si le chemin existe
	if filePath != "" {
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			utils.LogError(fmt.Sprintf("Impossible de supprimer le fichier backup_id=%d, path=%s", backupID, filePath), err)
		}
	}

	// Réponse de succès
	utils.LogInfo(fmt.Sprintf("Backup_id=%d supprimé par user_id=%d", backupID, userID))
	utils.SendSuccess(w, http.StatusOK, "Le backup a été supprimé avec succès")
	
}

// Backups pour toutes les bases de données (CRON)
func RunAllBackups() error {
	rows, err := db.DB.Query("SELECT id FROM databases")
	if err != nil {
		utils.LogError("Erreur lors de la récupération des bases", err)
	}

	defer rows.Close()

	for rows.Next() {
		var dbID int
		if err := rows.Scan(&dbID); err != nil {
			utils.LogError("Erreur lors de la récupération de l'ID de la base", err)
			continue
		}
		backupName := fmt.Sprintf("auto-%s", time.Now().Format("2006-01-02"))
		if err := CreateBackupForDb(dbID, backupName, "v1"); err != nil {
			utils.LogError("Impossible de créer un backup", err)
		} else {
			utils.LogInfo("Le backup de la base créé avec succès")
		}
	}
	return nil
}