package api

import (
	"sync/atomic"

	"github.com/ibraisi/chirpy/internal/database"
)

type Config struct {
	Hits      atomic.Int32
	DB        *database.Queries
	Platform  string
	SecretKey string
	PolkaKey  string
}
