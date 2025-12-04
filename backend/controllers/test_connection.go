package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"safebase/utils"
)

// Structure représentant les données existant dans la base des données
type TestConnectionRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Host string `json:"host"`
	Port int `json:"port"`
	DBUsername string `json:"db_username"`
	DBPassword string `json:"db_password"`
}

func TestConnection(w http.ResponseWriter, r *http.Request) {
var req TestConnectionRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	utils.SendError(w, http.StatusBadRequest, "Format JSON invalide")
	return
}

var dsn string
var driver string

if req.Type == "postgres" {
	driver = "postgres"
	dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=disable",
	req.DBUsername, req.DBPassword, req.Host, req.Port)
} else if req.Type == "mysql" {
	driver = "mysql"
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/",
	req.DBUsername, req.DBPassword, req.Host, req.Port)
} else {
	utils.SendError(w, http.StatusBadRequest, "Type de base non supporté")
	return
}

db, err := sql.Open(driver, dsn)
if err != nil {
	utils.SendError(w, 500, "Erreur lors de l'ouverture de la connexion")
	return
}
defer db.Close()
if err := db.Ping(); err != nil {
	utils.SendError(w, 400, "Connexion impossible : "+err.Error())
	return
}

utils.SendSuccess(w, 200, "Connexion réussie !")
}