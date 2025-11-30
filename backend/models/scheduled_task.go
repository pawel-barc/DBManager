package models

import (
	"time"
)

type ScheduledTask struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	DatabaseID int `json:"database_id"`
	CronExpression string `json:"cron_expression"`
	IsActive bool `json:"is_active"`
	LastRunAt *time.Time `json:"last_run_at"`
}