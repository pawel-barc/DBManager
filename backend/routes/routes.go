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

	// Base des données
	r.With(middleware.AuthMiddleware).Post("/databases/add", controllers.AddDatabase)
	r.With(middleware.AuthMiddleware).Get("/databases/list", controllers.GetDatabases)
	r.With(middleware.AuthMiddleware).Post("/databases/test", controllers.TestConnection)

	// Retourne le routeur configuré comme 'http.Handler'
	return r
}