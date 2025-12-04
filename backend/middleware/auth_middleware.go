package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"safebase/config"
	"safebase/utils"

	"github.com/golang-jwt/jwt/v5"
)

type key string

const UserIDKey key = "user_id"  // Clé pour stocker l'ID de l'utilisateur dans le contexte

// ----- MIDDLEWARE ----- pour la verification du token JWT du cookie
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Récupération du cookie contenant le token d'accès
		cookie, err := r.Cookie("access_token")
		if err != nil {
			// Aucun token trouvé -> accès refusé
			utils.SendError(w, http.StatusUnauthorized, "TOKEN_EXPIRED")
			return
		}

		// Nettoyage de la valeur du token
		tokenString := strings.TrimSpace(cookie.Value)

		// Parsing et validation du token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil // Clé secrète pour vérifier le token
		})

		// Si erreur ou token invalide -> accès refusé
		if err != nil || !token.Valid {
			utils.SendError(w, http.StatusUnauthorized, "Token invalide")
			return
		}
		
		// Extraction du payload du token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.SendError(w, http.StatusUnauthorized, "Format du token invalide")
			return
		}

		// Verification de la date d'expiration du token
		exp := int64(claims["exp"].(float64))
		if time.Now().Unix() > exp {
			utils.SendError(w, http.StatusUnauthorized, "TOKEN_EXPIRED")
			return
		}
		
		// Extraction de l'ID  utilisateur depuis les claims
		userIDFloat, ok := claims["user_id"].(float64) // JWT sérialise les nombres en float64
		if !ok {
			utils.SendError(w, http.StatusUnauthorized, "ID utilisateur invalide")
			return
		}

		// Conversion en int pour utilisation dans l'application
		userID := int(userIDFloat)

		// Stocakge de l'ID utilisateur dans le contexte de la requête
		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		// Appel du handler suivant en passant le contexte enrichi
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}