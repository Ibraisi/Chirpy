package auth

import (
	"github.com/alexedwards/argon2id"
	"github.com/goforj/godump"
)

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CheckPasswordHash(password, hash string) (bool, error) {
	godump.Dump(password)
	godump.Dump(hash)
	ok, _, err := argon2id.CheckHash(password, hash)
	godump.Dump(ok)
	return ok, err
}
