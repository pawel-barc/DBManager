package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"safebase/db"
	"safebase/middleware"
	"safebase/models"
	"safebase/utils"
	"safebase/validation"

	"golang.org/x/crypto/bcrypt"
)

// Response struct pour envoyer des message simples
type Response struct {
	Message string `json:"message"`
}

// ------ INSCRIPTION UTILISATEUR ------
func Register(w http.ResponseWriter, r *http.Request) {
	var req models.User
	// Décodage du JSON envoyé par le client
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Format JSON invalide")
		return
	}

	// ------VALIDATION------
	if err := validation.ValidatePassword(req.Password); err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validation.ValidateEmail(req.Email); err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validation.ValidateUsername(req.Username); err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Hash du mot de passe
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Erreur interne du serveur")
		return
	}
	req.Password = string(hashed)

	// Insertion dans la base des données
	_, err = db.DB.Exec(
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3)",
		req.Username, req.Email, req.Password,
	)

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "L'adresse e-mail existe déjà dans la base des données")
		return
	}

	// Response JSON pour confirmation
	utils.SendSuccess(w, http.StatusCreated, "Utilisateur enregistré avec succès!")
}

// ----- CONNEXION -----
func Login(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Format JSON invalide")
		return
	}

	var user models.User
	// Récupération de l'utilisateur depuis la base des données
	err := db.DB.QueryRow("SELECT id, username, email, password FROM users WHERE email=$1", 
	req.Email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, "Identifiants invalides")
		return
	}

	// Vérification du mot de passe
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		utils.SendError(w, http.StatusUnauthorized, "Identifiants invalides")
		return
	}

	// Création des tokens JWT
	accessToken, _ := utils.CreateToken(user.ID, user.Username, time.Minute*15)
	refreshToken, _ := utils.CreateToken(user.ID, user.Username, time.Hour*24*7)

	// Mis à jour du refresh token dans la base
	_, err = db.DB.Exec("UPDATE users SET refresh_token=$1 WHERE id=$2", refreshToken, user.ID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Erreur interne du serveur!")
		return
	}

	// Envoi des cookies
	http.SetCookie(w, &http.Cookie{
		Name: "access_token",
		Value: accessToken,
		Path: "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name: "refresh_token",
		Value: refreshToken,
		Path: "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, 
	})

	// Response JSON avec les informations de l'utilisateur
	utils.SendUserSuccess(w, http.StatusOK, map[string]interface{}{
		"id": user.ID,
		"username": user.Username,
		"email": user.Email,
	}, "Connexion réussie")
}

// -----RAFRAÎCHISSEMENT DU TOKEN D'ACCES-------
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, "Aucun token de rafraîchissement trouvé",)
		return
	}

	var user models.User
	err = db.DB.QueryRow("SELECT id, username FROM users WHERE refresh_token=$1", 
	cookie.Value).Scan(&user.ID, &user.Username)
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, "Token de rafraîchissement invalide")
		return
	}

	newAccessToken, _ := utils.CreateToken(user.ID, user.Username, time.Minute*15)

	http.SetCookie(w, &http.Cookie{
		Name: "access_token",
		Value: newAccessToken,
		Path: "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	utils.SendSuccess(w, http.StatusOK, "Token mis à jour avec succès")
}

// ------ ENDPOINT PROTEGE ----- pour récupérer les infos de l'utilisateur
func Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var user models.User
	err := db.DB.QueryRow(
		"SELECT id, username, email FROM users WHERE id=$1", 
		userID,
	).Scan(&user.ID, &user.Username, &user.Email)

	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Utilisateur introuvable")
		return
	}

	utils.SendUserSuccess(w, http.StatusOK, user, "That's me")

}

// ----- DECONNEXION DE L'UTILISATEUR ------
func Logout(w http.ResponseWriter, r *http.Request) {
	// Suppression des cookies
	http.SetCookie(w, &http.Cookie{
		Name: "access_token",
		Value: "",
		Path: "/",
		HttpOnly: true,
		MaxAge: -1,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name: "refresh_token",
		Value: "",
		Path: "/",
		HttpOnly: true,
		MaxAge: -1,
		SameSite: http.SameSiteLaxMode,
	})

	utils.SendSuccess(w, http.StatusOK, "Déconnexion réussie")
}