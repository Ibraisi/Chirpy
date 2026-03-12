package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/ibraisi/chirpy/internal/api"
	"github.com/ibraisi/chirpy/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		panic("not able to load env")
	}

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		panic("not able to open sql")
	}

	cfg := &api.Config{
		DB:        database.New(db),
		Platform:  os.Getenv("PLATFORM"),
		SecretKey: os.Getenv("SECRET_KEY"),
	}

	mux := http.NewServeMux()
	wrap := func(h http.HandlerFunc) http.Handler { return api.LoggingMiddleware(h) }

	fileServer := http.StripPrefix("/app/", http.FileServer(http.Dir("../")))
	mux.Handle("GET /app/", cfg.HitsMiddleware(fileServer))

	mux.Handle("GET /admin/metrics", wrap(cfg.Metrics))
	mux.Handle("POST /admin/reset", wrap(cfg.Reset))

	mux.Handle("GET /api/healthz", wrap(api.Readiness))
	mux.Handle("POST /api/login", wrap(cfg.Login))
	mux.Handle("POST /api/refresh", wrap(cfg.RefreshToken))
	mux.Handle("POST /api/revoke", wrap(cfg.RevokeRefreshToken))
	mux.Handle("POST /api/users", wrap(cfg.CreateUser))
	mux.Handle("PUT /api/users", wrap(cfg.UpdateUser))

	mux.Handle("POST /api/chirps", wrap(cfg.CreateChirp))
	mux.Handle("GET /api/chirps", wrap(cfg.GetChirps))
	mux.Handle("GET /api/chirps/{chirpID}", wrap(cfg.GetChirpByID))
	mux.Handle("DELETE /api/chirps/{chirpID}", wrap(cfg.DeleteChirpByID))

	server := &http.Server{Addr: ":8080", Handler: mux}
	log.Fatal(server.ListenAndServe())
}
