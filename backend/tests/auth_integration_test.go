// backend/tests/auth_integration_test.go
package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"safebase/controllers"
	"safebase/db"

	"github.com/gorilla/mux"
)

// setupRouter crée un mini router pour nos tests d'intégration
func setupRouter() http.Handler {
	r := mux.NewRouter()

	// Routes d'authentification
	r.HandleFunc("/register", controllers.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", controllers.Login).Methods(http.MethodPost)
	r.HandleFunc("/refresh", controllers.RefreshToken).Methods(http.MethodPost)
	r.HandleFunc("/logout", controllers.Logout).Methods(http.MethodPost)

	return r
}

// Optionnel : helper pour nettoyer la table users entre les tests
func cleanUsersTable(t *testing.T) {
	if db.DB == nil {
		t.Log("db.DB est nil, impossible de nettoyer la table users (pense à initialiser la DB avant les tests)")
		return
	}
	_, err := db.DB.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("échec du nettoyage de la table users: %v", err)
	}
}

// ----- TEST 1 : inscription + connexion OK, cookies posés -----

func TestRegisterAndLogin_SetCookies(t *testing.T) {
	router := setupRouter()

	// Nettoyage de la table users pour partir sur une base propre (optionnel)
	cleanUsersTable(t)

	// --- Inscription ---
	registerBody := []byte(`{
		"username": "alice",
		"email": "alice@mail.com",
		"password": "Password1!"
	}`)

	reqReg := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(registerBody))
	reqReg.Header.Set("Content-Type", "application/json")

	recReg := httptest.NewRecorder()
	router.ServeHTTP(recReg, reqReg)

	if recReg.Code != http.StatusCreated {
		t.Fatalf("expected status %d for register, got %d. Body: %s",
			http.StatusCreated, recReg.Code, recReg.Body.String())
	}

	// --- Connexion ---
	loginBody := []byte(`{
		"email": "alice@mail.com",
		"password": "Password1!"
	}`)

	reqLogin := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")

	recLogin := httptest.NewRecorder()
	router.ServeHTTP(recLogin, reqLogin)

	if recLogin.Code != http.StatusOK {
		t.Fatalf("expected status %d for login, got %d. Body: %s",
			http.StatusOK, recLogin.Code, recLogin.Body.String())
	}

	// Vérifier que les cookies access_token et refresh_token sont bien posés
	resp := recLogin.Result()
	defer resp.Body.Close()

	var accessCookie *http.Cookie
	var refreshCookie *http.Cookie

	for _, c := range resp.Cookies() {
		if c.Name == "access_token" {
			accessCookie = c
		}
		if c.Name == "refresh_token" {
			refreshCookie = c
		}
	}

	if accessCookie == nil {
		t.Errorf("expected access_token cookie to be set")
	}
	if refreshCookie == nil {
		t.Errorf("expected refresh_token cookie to be set")
	}
}

// ----- TEST 2 : mauvais mot de passe -> 401 et aucun cookie d'auth -----

func TestLoginWithWrongPassword_NoCookies(t *testing.T) {
	router := setupRouter()

	// Nettoyage de la table users (optionnel)
	cleanUsersTable(t)

	// On crée d'abord un utilisateur valide
	registerBody := []byte(`{
		"username": "bob",
		"email": "bob@mail.com",
		"password": "Password1!"
	}`)

	reqReg := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(registerBody))
	reqReg.Header.Set("Content-Type", "application/json")

	recReg := httptest.NewRecorder()
	router.ServeHTTP(recReg, reqReg)

	if recReg.Code != http.StatusCreated {
		t.Fatalf("expected status %d for register, got %d. Body: %s",
			http.StatusCreated, recReg.Code, recReg.Body.String())
	}

	// Tentative de connexion avec mauvais mot de passe
	loginBody := []byte(`{
		"email": "bob@mail.com",
		"password": "WrongPassword123!"
	}`)

	reqLogin := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")

	recLogin := httptest.NewRecorder()
	router.ServeHTTP(recLogin, reqLogin)

	if recLogin.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d for login with wrong password, got %d. Body: %s",
			http.StatusUnauthorized, recLogin.Code, recLogin.Body.String())
	}

	// Vérifie qu'aucun cookie d'auth n'est posé
	resp := recLogin.Result()
	defer resp.Body.Close()

	for _, c := range resp.Cookies() {
		if c.Name == "access_token" {
			t.Errorf("did not expect access_token cookie on failed login")
		}
		if c.Name == "refresh_token" {
			t.Errorf("did not expect refresh_toke759652n cookie on failed login")
		}
	}
}
