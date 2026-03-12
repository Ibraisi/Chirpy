package api

import (
	"database/sql"
	"encoding/json"
	"github.com/ibraisi/chirpy/internal/auth"
	"github.com/ibraisi/chirpy/pkg/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type loginRequest struct {
	Email        string         `json:"email"`
	Password     string         `json:"password"`
	ExpiresEfter *time.Duration `json:"expires_in_seconds"`
}

type loginResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
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
	tokenExpiry := time.Duration(time.Hour)
	if req.ExpiresEfter != nil {
		tokenExpiry = *req.ExpiresEfter
	}
	token, err := auth.MakeJWT(user.ID, cfg.SecretKey, tokenExpiry)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	loginRes := loginResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email.String,
		Token:     token,
	}

	utils.ResponseWithJSON(w, http.StatusOK, loginRes)
}
