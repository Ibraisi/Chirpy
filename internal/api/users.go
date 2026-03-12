package api

import (
	"database/sql"
	"encoding/json"
	"github.com/ibraisi/chirpy/internal/auth"
	"github.com/ibraisi/chirpy/internal/database"
	"github.com/ibraisi/chirpy/pkg/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func userFromDB(u database.User) userResponse {
	return userResponse{
		ID:        u.ID,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
		Email:     u.Email.String,
	}
}

func (cfg *Config) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	hashed, err := auth.HashPassword(req.Password)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to process password"})
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:          sql.NullString{String: req.Email, Valid: true},
		HashedPassword: hashed,
	})
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	utils.ResponseWithJSON(w, http.StatusCreated, userFromDB(user))
}
