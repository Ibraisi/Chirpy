package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ibraisi/chirpy/internal/auth"
	"github.com/ibraisi/chirpy/internal/database"
	"github.com/ibraisi/chirpy/pkg/utils"

	"github.com/google/uuid"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshTokenRes struct {
	Token string `json:"token"`
}

type loginResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *Config) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	user, err := cfg.DB.GetUserByEmail(r.Context(), sql.NullString{String: req.Email, Valid: true})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ok, _ := auth.CheckPasswordHash(req.Password, user.HashedPassword)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.SecretKey, time.Duration(time.Hour))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	refreshToken := auth.MakeRefreshToken()
	err = cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
		ExpiresAt: sql.NullTime{
			Time:  time.Now().Add(time.Hour * 24 * 60),
			Valid: true,
		},
		RevokedAt: sql.NullTime{},
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, loginResponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		Email:        user.Email.String,
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func (cfg *Config) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	refToken, err := cfg.DB.GetRefreshToken(r.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if refToken.ExpiresAt.Valid && time.Now().After(refToken.ExpiresAt.Time) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newToken, err := auth.MakeJWT(refToken.UserID.UUID, cfg.SecretKey, time.Hour)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, refreshTokenRes{
		Token: newToken,
	})

}

func (cfg *Config) RevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err := cfg.DB.RevokeToken(r.Context(), token); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
