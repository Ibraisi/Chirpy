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

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type upgradeUserRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	}
}

type userResponse struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func userFromDB(u database.User) userResponse {
	return userResponse{
		ID:          u.ID,
		CreatedAt:   u.CreatedAt.Time,
		UpdatedAt:   u.UpdatedAt.Time,
		Email:       u.Email.String,
		IsChirpyRed: u.IsChirpyRed,
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

func (cfg *Config) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromCtx(r)

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	hashed, err := auth.HashPassword(req.Password)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to process password"})
		return
	}

	user, err := cfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		Email: sql.NullString{
			String: req.Email,
			Valid:  true,
		},
		HashedPassword: hashed,
		ID:             userID,
	})
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, userFromDB(user))
}

func (cfg *Config) UpgradeUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.PolkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req upgradeUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	id, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id format"})
		return
	}

	if err = cfg.DB.UpgradeUser(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
