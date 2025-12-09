package routes

import (
	"net/http"

	"safebase/controllers"
	"safebase/middleware"

	"github.com/go-chi/chi/v5"
)

// SetupRouter configure toutes les routes de l'application
func SetupRouter() http.Handler {
	// Création du routeur principal(chi framework comme express)
	r := chi.NewRouter()

	// Activation du middleware CORS
	r.Use(middleware.CORSHandler())

	// ----- ROUTES PUBLIQUES -----
	r.Post("/register", controllers.Register)
	r.Post("/login", controllers.Login)
	r.Post("/refresh-token", controllers.RefreshToken)

	// ----- ROUTES PROTEGEES ------ par le middleware d'authentification

	// Deconnexion, Profil
	r.With(middleware.AuthMiddleware).Get("/me", controllers.Me)
	r.With(middleware.AuthMiddleware).Post("/logout", controllers.Logout)
	r.With(middleware.AuthMiddleware).Get("/get-profile", controllers.GetProfile)
	r.With(middleware.AuthMiddleware).Put("/update-profile", controllers.UpdateProfile)
	r.With(middleware.AuthMiddleware).Delete("/delete-profile", controllers.DeleteAccount)

	// Base des données
	r.With(middleware.AuthMiddleware).Post("/databases/add", controllers.AddDatabase)
	r.With(middleware.AuthMiddleware).Get("/databases/list", controllers.GetDatabases)
	r.With(middleware.AuthMiddleware).Post("/databases/test", controllers.TestConnection)
	r.With(middleware.AuthMiddleware).Delete("/databases/{id}", controllers.DeleteDatabase)

	// Backups
	r.With(middleware.AuthMiddleware).Post("/backups/{database_id}/create", controllers.CreateBackup)
	r.With(middleware.AuthMiddleware).Get("/backups", controllers.ListAllBackups)
	r.With(middleware.AuthMiddleware).Get("/backups/{database_id}", controllers.ListBackups)
	r.With(middleware.AuthMiddleware).Get("/backups/{backup_id}/download", controllers.DownloadBackup)
	r.With(middleware.AuthMiddleware).Delete("/backups/{backup_id}", controllers.DeleteBackup)

	// CRON
	r.With(middleware.AuthMiddleware).Post("/cron/create", controllers.CreateScheduledTask)
	r.With(middleware.AuthMiddleware).Get("/scheduled-tasks", controllers.GetScheduledTasks)
	r.With(middleware.AuthMiddleware).Put("/scheduled-tasks/toggle", controllers.ToggleTaskActive)
	r.With(middleware.AuthMiddleware).Put("/scheduled-tasks/update-cron", controllers.UpdateCronExpression)
	r.With(middleware.AuthMiddleware).Delete("/scheduled-tasks/{id}", controllers.DeleteScheduledTask)

	// Restaurations
	r.With(middleware.AuthMiddleware).Post("/backups/{backup_id}/restore", controllers.RestoreBackup)

	// Alerts
	r.With(middleware.AuthMiddleware).Get("/alerts", controllers.GetUserAlerts)
	r.With(middleware.AuthMiddleware).Put("/alerts/{alert_id}/read", controllers.MarkAlertAsRead)
	r.With(middleware.AuthMiddleware).Put("/alerts/read-all", controllers.MarkAllAsRead)
	// Retourne le routeur configuré comme 'http.Handler'
	return r
}
