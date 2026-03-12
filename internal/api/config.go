package api

import (
	"github.com/ibraisi/chirpy/internal/database"
	"sync/atomic"
)

type Config struct {
	Hits      atomic.Int32
	DB        *database.Queries
	Platform  string
	SecretKey string
}
