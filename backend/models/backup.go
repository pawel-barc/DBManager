package models

import "time"

type Backup struct {
	ID         uint      `json:"id"`
	DatabaseID uint      `json:"database_id"`
	Name       string    `json:"name"`
	FilePath   string    `json:"file_path"`
	FileSize   float64   `json:"file_size"`             // double precision
	BackupDate time.Time `json:"backup_date,omitempty"` // peut être null tant que non terminé
	Status     string    `json:"status"`                // queued | running | done | error
	Version    string    `json:"version"`               // ex: "pg_dump 16.3"
	Log        string    `json:"log"`                   // log texte du dump
	CreatedAt  time.Time `json:"created_at,omitempty"`  // si tu veux garder une trace de création (optionnel côté DB)
}
