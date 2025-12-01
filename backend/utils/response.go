package utils

import (
	"encoding/json"
	"net/http"
)

// JsonResponse - structure générique pour toutes les responses JSON envoyées au client
type JsonResponse struct {
	Success bool        `json:"success,omitempty"`
	Error   string      `json:"error,omitempty"`
	User    interface{} `json:"user,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

// SendSuccess - envoie une réponse JSON pour une action réussie contenant des informations utilisateur
func SendSuccess(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(JsonResponse{
		Success: true,
		Message: message,
	})
}

// SendUserSuccess - envoie une réponse JSON pour une action réussie simple(sans données spécifiques)
func SendUserSuccess(w http.ResponseWriter, status int, user interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(JsonResponse{
		Success: true,
		User:    user,
		Message: message,
	})
}

// SendError - envoie une reponse JSON en cas d'erreur
func SendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(JsonResponse{
		Success: false,
		Error:   message,
		Message: message,
	})
}

// SendJSON - envoie une réponse JSON générique avec n’importe quelle donnée
func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(JsonResponse{
		Success: true,
		Data:    data,
	})
}
