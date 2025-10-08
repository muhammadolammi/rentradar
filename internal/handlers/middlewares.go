package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/muhammadolammi/rentradar/internal/helpers"
)

// Middleware to check for the API key in the authorization header for all requests.
func (apiConfig *Config) VerifyApiKey() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			api_key := r.Header.Get("API-KEY")
			if api_key == "" {
				helpers.RespondWithError(w, http.StatusUnauthorized, "missing API-KEY header")
				return
			}
			if api_key != apiConfig.APIKEY {
				helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid API-KEY key")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Middleware to check for the AUTHORIZATION in user only enpoints in the authorization header for all requests.
func (apiConfig *Config) AuthMiddleware(requireSudo bool, jwtKey []byte, next func(http.ResponseWriter, *http.Request, User)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Missing or invalid token")
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		authclaims := &jwt.RegisteredClaims{}

		authJwt, err := jwt.ParseWithClaims(
			tokenString,
			authclaims,
			func(token *jwt.Token) (interface{}, error) { return []byte(apiConfig.JWTKEY), nil },
		)

		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error parsing jwt claims, err: %v", err))
			return
		}

		if authclaims.ExpiresAt != nil && authclaims.ExpiresAt.Time.Before(time.Now().UTC()) {
			helpers.RespondWithError(w, http.StatusUnauthorized, "auth token expired")
			return
		}

		userId, err := authJwt.Claims.GetIssuer()
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting issuer from jwt claims, err: %v", err))
			return
		}
		id, err := uuid.Parse(userId)
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error parsing id, err: %v", err))
			return
		}
		user, err := apiConfig.DB.GetUser(r.Context(), id)
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting user, err: %v", err))
			return
		}

		// --- Optional SUDO Verification ---
		if requireSudo {
			sudoKey := r.Header.Get("SUDO-KEY")
			if sudoKey == "" {
				helpers.RespondWithError(w, http.StatusUnauthorized, "missing SUDO-KEY header")
				return
			}
			if sudoKey != apiConfig.SUDOKEY {

				helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid SUDO-API key")
				return
			}
		}
		next(w, r, DbUserToModelsUser(user))
	})
}
