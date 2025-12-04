package models

import "time"

type Backup struct {
ID int `json:"id"`
DatabaseID int `json:"database_id"`
Name string `json:"name"`
FilePath string `json:"file_path"`
FileSize float64 `json:"file_size"`
BackupDate *time.Time `json:"backup_date"`
Status string `json:"status"`
Version string `json:"version"`
Log string `json:"log"`
}