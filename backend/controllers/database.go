package controllers

// Connexion et gestion des bases de données externes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"safebase/db"
	"safebase/middleware"
	"safebase/utils"
)

// ---------- Types de requêtes/réponses ----------

// Payload reçu du frontend pour créer une DB
type AddDatabaseRequest struct {
	Name       string `json:"name"`
	Type       string `json:"type"` // ex: "PostgreSQL" | "MySQL"
	Host       string `json:"host"`
	Port       int    `json:"port"`
	DBUsername string `json:"db_username"`
	DBPassword string `json:"db_password"`
}

// Représentation "publique" d'une DB (jamais de password exposé)
type DatabasePublic struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Host       string    `json:"host"`
	Port       int       `json:"port"`
	DBUsername string    `json:"db_username"`
	CreatedAt  time.Time `json:"created_at"`
}

// ---------- Helpers ----------

var (
	// nom: lettres/chiffres/underscore/tiret, 3..64
	reDBName = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,64}$`)
	// host très permissif (FQDN ou IP v4 simple). Ajuste si besoin.
	reHost = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
)

func getUserIDFromCtx(r *http.Request) (int, error) {
	val := r.Context().Value(middleware.UserIDKey)
	if val == nil {
		return 0, errors.New("utilisateur non authentifié")
	}
	id, ok := val.(int)
	if !ok {
		return 0, errors.New("identifiant utilisateur invalide")
	}
	return id, nil
}

func normalizeDBType(t string) string {
	t = strings.TrimSpace(strings.ToLower(t))
	switch t {
	case "postgres", "postgresql":
		return "PostgreSQL"
	case "mysql":
		return "MySQL"
	default:
		// Laisse passer d'autres moteurs si besoin, mais normalise la casse.
		if t == "" {
			return ""
		}
		return strings.Title(t)
	}
}

func validateAddDB(req *AddDatabaseRequest) error {
	if !reDBName.MatchString(req.Name) {
		return errors.New("name invalide (3-64, lettres/chiffres/_/-)")
	}
	if req.Type = normalizeDBType(req.Type); req.Type == "" {
		return errors.New("type requis (ex: PostgreSQL | MySQL)")
	}
	if !reHost.MatchString(req.Host) {
		return errors.New("host invalide")
	}
	if req.Port < 1 || req.Port > 65535 {
		return errors.New("port invalide (1-65535)")
	}
	if strings.TrimSpace(req.DBUsername) == "" {
		return errors.New("db_username requis")
	}
	// db_password peut être vide si la cible accepte des méthodes alternatives,
	// mais dans la plupart des cas on le requiert :
	if strings.TrimSpace(req.DBPassword) == "" {
		return errors.New("db_password requis")
	}
	return nil
}

// ---------- Handlers ----------

// POST /api/databases
// Ajoute une base pour l'utilisateur courant
func AddDatabase(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromCtx(r)
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var req AddDatabaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Format JSON invalide")
		return
	}
	if err := validateAddDB(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Insertion
	var newID int
	// Utilise la requête centralisée si tu l'as définie dans db/queries.go.
	// Ici, on écrit la requête inline pour rester autonome.
	q := `
		INSERT INTO databases (user_id, name, type, host, port, db_username, db_password, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,NOW())
		RETURNING id;
	`
	err = db.DB.QueryRowContext(r.Context(), q,
		userID, req.Name, req.Type, req.Host, req.Port, req.DBUsername, req.DBPassword,
	).Scan(&newID)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Impossible d'ajouter la base")
		return
	}

	// Réponse JSON avec l'objet créé (public, sans password)
	resp := DatabasePublic{
		ID:         newID,
		UserID:     userID,
		Name:       req.Name,
		Type:       req.Type,
		Host:       req.Host,
		Port:       req.Port,
		DBUsername: req.DBUsername,
		CreatedAt:  time.Now().UTC(),
	}
	// 201 Created
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// GET /api/databases
// Liste les bases de l'utilisateur courant
func ListDatabases(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromCtx(r)
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, err.Error())
		return
	}

	q := `
		SELECT id, user_id, name, type, host, port, db_username, created_at
		FROM databases
		WHERE user_id = $1
		ORDER BY created_at DESC, id DESC;
	`

	rows, err := db.DB.QueryContext(r.Context(), q, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Impossible de charger les bases")
		return
	}
	defer rows.Close()

	out := make([]DatabasePublic, 0, 16)
	for rows.Next() {
		var d DatabasePublic
		if err := rows.Scan(&d.ID, &d.UserID, &d.Name, &d.Type, &d.Host, &d.Port, &d.DBUsername, &d.CreatedAt); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Erreur de lecture des bases")
			return
		}
		out = append(out, d)
	}
	if err := rows.Err(); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Erreur de curseur")
		return
	}

	// Retourne un tableau JSON (ce que ton front attend)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(out)
}

// (Optionnel) GET /api/databases/:id – détail (sans password)
func GetDatabase(w http.ResponseWriter, r *http.Request, id int) {
	userID, err := getUserIDFromCtx(r)
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, err.Error())
		return
	}

	q := `
		SELECT id, user_id, name, type, host, port, db_username, created_at
		FROM databases
		WHERE id = $1 AND user_id = $2
		LIMIT 1;
	`

	var d DatabasePublic
	err = db.DB.QueryRowContext(r.Context(), q, id, userID).Scan(
		&d.ID, &d.UserID, &d.Name, &d.Type, &d.Host, &d.Port, &d.DBUsername, &d.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		utils.SendError(w, http.StatusNotFound, "Base introuvable")
		return
	}
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Impossible de charger la base")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(d)
}
