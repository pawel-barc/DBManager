package controllers

import (
	"net/http"

	"safebase/db"
	"safebase/middleware"
	"safebase/utils"
)

// Struct utilisée pour la mis à jour des données d'utilisateur
type UpdateProfileRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password,omitempty"`
}

// ----- GET PROFILE ----- récupère le profil de l'utilisateur connecté
func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var user struct {
		ID uint `json:"id"`
		Username string `json:"username"`
		Email string `json:"email"`
	}

	// Requête pour récupérer les informations utilisateur
	err := db.DB.QueryRow(
		"SELECT id, username, email FROM users WHERE id=$1",
		userID,
	).Scan(&user.ID, &user.Username, &user.Email)

	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Utilisateur introuvable")
		return
	}

	// Response JSON standarisée avec message
	utils.SendUserSuccess(w, http.StatusOK, user, "Profil récupéré avec succès")
}