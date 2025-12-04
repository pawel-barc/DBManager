package controllers

// Connexion et gestion des bases de données externes

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"safebase/db"
	"safebase/middleware"
	"safebase/utils"

	"github.com/go-chi/chi/v5"
)

// Structure représentant les données envoyées par le frontend
type AddDatabaseRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Host string `json:"host"`
	Port int `json:"port"`
	DBUsername string `json:"db_username"`
	DBPassword string `json:"db_password"`
}

// Structure représentant les données existant dans la base des données
type DatabasesItem struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Host string `json:"host"`
	Port int `json:"port"`
	DBUsername string `json:"db_username"`
}

// Fonction pour ajouter une nouvelle base des données
func AddDatabase(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID de l'utilisateur connecté
	UserID := r.Context().Value(middleware.UserIDKey).(int)

	// Déclaration d'une variable pour stocker les données du corps de la requête
	var req AddDatabaseRequest
	// Décodage du JSON envoyé par le frontend
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Format JSON invalide")
		return
	}

	// Vérification des doublons
	var exists bool
	checkQuerry := `SELECT EXISTS(
	SELECT 1 FROM databases WHERE user_id = $1 AND name = $2 AND type = $3 AND host = $4 AND port = $5 AND db_username = $6)`

	err := db.DB.QueryRow(checkQuerry, UserID, req.Name, req.Type, req.Host, req.Port, req.DBUsername).Scan(&exists)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Erreur interne de vérification")
		return
	}

	if exists {
		utils.SendError(w, http.StatusBadRequest, "Cette base est déjà enregistrée")
		return
	}

	// Variable pour récupérer l'ID de la base nouvellement créée
	var dbID int
	// Exécution de la requête SQL d'insertion (définie dans db/queries.go)
	err = db.DB.QueryRow(
		db.QueryInsertDatabase, UserID, req.Name, req.Type, req.Host, req.Port, req.DBUsername, req.DBPassword, 
	).Scan(&dbID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Impossible d'ajouter la base")
		return
	}

	utils.SendSuccess(w, http.StatusCreated, "La base ajouté avec succès !")
}

// Récupération des bases existantes
func GetDatabases(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	rows, err := db.DB.Query(db.QuerySelectDatabases, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Erreur lors de la récupération des bases")
		return
	}
	defer rows.Close()

	var databases []DatabasesItem

	for rows.Next() {
		var d DatabasesItem
		if err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.Type,
			&d.Host,
			&d.Port,
			&d.DBUsername,
		); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Erreur lors de la lecture des données")
			return
		}
	databases = append(databases, d)	
	}

	utils.SendSuccessWithData(w, http.StatusOK, "Liste des bases récupérée", databases)
}

// Suppression d'une base uniquement dans l'application
func DeleteDatabase(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Récupération de l'ID de la base
	dbIDStr := chi.URLParam(r, "id")
	dbID, err := strconv.Atoi(dbIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ID invalide")
		return
	}

	// Vérification que la base appartient bien à l'utilisateur
	var exists bool

	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM databases WHERE id = $1 AND user_id = $2)", dbID, userID).Scan(&exists)
	
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Erreur lors de la vérification")
		return
	}

	if !exists {
		utils.SendError(w, http.StatusNotFound, "Base non trouvée")
		return
	}

	// Suppression des fichiers de backup associés
	rows, err := db.DB.Query("SELECT file_path FROM backups WHERE database_id = $1", dbID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Erreur lors de la récupération des backups")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var filePath string
		if err := rows.Scan(&filePath); err == nil && filePath != "" {
			if removeErr := os.Remove(filePath); removeErr != nil && !os.IsNotExist(removeErr) {
				utils.LogError("Impossible de supprimer un fichier backup "+filePath, removeErr)
			}
		}
	}
	// Suppression des backups liés
	_, err = db.DB.Exec("DELETE FROM backups WHERE database_id = $1", dbID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Erreur lors de la suppression des backups associés")
		return
	}
	// Suppression de la base dans l'application
	_, err = db.DB.Exec("DELETE FROM databases WHERE id=$1 AND user_id=$2", dbID, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Erreur interne lors de la suppression de la base")
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Base supprimée de l'application")
	utils.LogInfo("La base supprimée avec suucès")
}