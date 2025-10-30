package controllers

// Connexion et gestion des bases de données externes

import (
	"encoding/json"
	"net/http"

	"safebase/db"
	"safebase/middleware"
	"safebase/utils"
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

// Fonction pour ajouter une nouvelle base des données
func AddDatabase(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID de l'utilisateur connecté
	UserID := r.Context().Value(middleware.UserIDKey).(int)

	// Déclaration d'une variable pour stocker les données du corps de la requête
	var req AddDatabaseRequest
	// Décodage du JSON envoyé par le frontend
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w,http.StatusBadRequest, "Format JSON invalide")
		return
	}

	// Variable pour récupérer l'ID de la base nouvellement créée
	var dbID int
	// Exécution de la requête SQL d'insertion (définie dans db/queries.go)
	err := db.DB.QueryRow(
		db.QueryInsertDatabase, UserID, req.Name, req.Type, req.Host, req.Port, req.DBUsername, req.DBPassword, 
	).Scan(&dbID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Impossible d'ajouter la base")
		return
	}

	utils.SendSuccess(w, http.StatusCreated, "La base ajouté avec succès !")
}