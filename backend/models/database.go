package models

import "time"

// Modèle de connexion à une base de données externe

type Database struct {
	ID uint `json:"id"`
	UserID uint `json:"user_id"`
	Name string `json:"name"`
	Type string `json:"type"`// ex. MySQL / PostgreSQL
	Host string `json:"host"`
	Port int `json:"port"`
	DBUsername string `json:"db_username"`
	DBPassword string `json:"db_password"`
	CreatedAt time.Time `json:"created_at"`
}