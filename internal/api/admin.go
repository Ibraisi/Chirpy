package api

import (
	"fmt"
	"log"
	"net/http"
)

func Readiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Printf("Method %s Path %s failed to write a response", r.Method, r.URL)
	}
}

func (cfg *Config) Metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprintf(w, `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.Hits.Load()); err != nil {
		log.Printf("Method %s Path %s failed to write a response", r.Method, r.URL)
	}
}

func (cfg *Config) Reset(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err := cfg.DB.DeleteAllUsers(r.Context()); err != nil {
		log.Printf("reset: delete all users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cfg.Hits.Store(0)
	w.WriteHeader(http.StatusOK)
}
