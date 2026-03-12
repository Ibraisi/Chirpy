package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/ibraisi/chirpy/internal/database"
	"github.com/ibraisi/chirpy/pkg/utils"

	"github.com/google/uuid"
)

var profane = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

type chirpRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}

func chirpFromDB(c database.Chirp) chirpResponse {
	return chirpResponse{
		ID:        c.ID,
		CreatedAt: c.CreatedAt.Time,
		UpdatedAt: c.UpdatedAt.Time,
		Body:      c.Body.String,
		UserID:    c.UserID.UUID.String(),
	}
}

func cleanBody(body string) string {
	words := strings.Fields(body)
	for i, w := range words {
		if _, ok := profane[strings.ToLower(w)]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func (cfg *Config) CreateChirp(w http.ResponseWriter, r *http.Request) {
	id := userIDFromCtx(r)

	var req chirpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if len(req.Body) > 400 {
		utils.ResponseWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Chirp is too long"})
		return
	}

	chirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: sql.NullString{String: cleanBody(req.Body), Valid: true},
		UserID: uuid.NullUUID{
			UUID:  id,
			Valid: true,
		},
	})
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	utils.ResponseWithJSON(w, http.StatusCreated, chirpFromDB(chirp))
}

func (cfg *Config) GetChirps(w http.ResponseWriter, r *http.Request) {
	var (
		chirps []database.Chirp
		err    error
	)

	if authorID := r.URL.Query().Get("author_id"); authorID != "" {
		id, parseErr := uuid.Parse(authorID)
		if parseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		chirps, err = cfg.DB.GetAllChirpsByUser(r.Context(), uuid.NullUUID{UUID: id, Valid: true})
	} else {
		chirps, err = cfg.DB.GetAllChirps(r.Context())
	}
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch chirps"})
		return
	}

	res := make([]chirpResponse, len(chirps))
	for i, c := range chirps {
		res[i] = chirpFromDB(c)
	}

	if r.URL.Query().Get("sort") == "desc" {
		slices.Reverse(res)
	}

	utils.ResponseWithJSON(w, http.StatusOK, res)
}

func (cfg *Config) GetChirpByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	chirp, err := cfg.DB.GetAllChirpsByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, chirpFromDB(chirp))
}

func (cfg *Config) DeleteChirpByID(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromCtx(r)

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chirp, err := cfg.DB.GetAllChirpsByID(r.Context(), chirpID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if chirp.UserID.UUID != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = cfg.DB.DeleteChirpsByID(r.Context(), chirpID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
