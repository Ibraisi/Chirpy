package api

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/ibraisi/chirpy/internal/auth"
)

type contextKey string

const userIDKey contextKey = "userID"

func LoggingMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL)
		next.ServeHTTP(w, r)
	}
}

func (cfg *Config) HitsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Hits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *Config) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID, err := auth.ValidateJWT(token, cfg.SecretKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// userIDFromCtx extracts the authenticated user ID from the request context.
func userIDFromCtx(r *http.Request) uuid.UUID {
	return r.Context().Value(userIDKey).(uuid.UUID)
}
