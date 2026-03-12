package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:   "chirpy-acess",
			Subject:  userID.String(),
			Audience: jwt.ClaimStrings{},
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(expiresIn).UTC(),
			},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now().UTC(),
			},
		},
	)

	tokenDat, err := token.SignedString([]byte(tokenSecret))
	return string(tokenDat), err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, err
	}
	return uuid.Parse(claims.Subject)
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}
	if !strings.Contains(authHeader, "Bearer") {
		return "", errors.New("not valid Authorization header")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token, nil
}

func MakeRefreshToken() string {
	tokenDER := make([]byte, 32)
	_, err := rand.Read(tokenDER)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(tokenDER)
}
