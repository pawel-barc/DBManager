package controllers

import (
	"encoding/json"
	"net/http"
	"safebase/db"
	"safebase/middleware"
	"safebase/utils"
	"safebase/validation"

	"golang.org/x/crypto/bcrypt"
)

// Struct utilisée pour la mis à jour des données d'utilisateur
type UpdateProfileRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

// ----- GET PROFILE ----- récupère le profil de l'utilisateur connecté
func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var user struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
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

// Mise à jour du profil
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Format JSON invalide")
		return
	}

	// Validation des champs
	if err := validation.ValidateUsername(req.Username); err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validation.ValidateEmail(req.Email); err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Password != "" {
		if err := validation.ValidatePassword(req.Password); err != nil {
			utils.SendError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	// Hash du mot de passe si fourni
	var hashedPassword string
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Erreur interne lors du hash du mot de passe")
			return
		}
		hashedPassword = string(hash)
	}

	// Mise à jour SQL
	var err error
	if req.Password == "" {
		_, err = db.DB.Exec(
			"UPDATE users SET username=$1, email=$2 WHERE id=$3",
			req.Username, req.Email, userID,
		)
	} else {
		_, err = db.DB.Exec(
			"UPDATE users SET username=$1, email=$2, password=$3 WHERE id=$4",
			req.Username, req.Email, hashedPassword, userID,
		)
	}

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Erreur lors de la mise à jour du profil")
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Profil mis à jour avec succès")
}
