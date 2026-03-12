package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}
	if !strings.Contains(authHeader, "ApiKey") {
		return "", errors.New("not valid Authorization header")
	}
	apiKey := strings.TrimPrefix(authHeader, "ApiKey ")
	return apiKey, nil
}
