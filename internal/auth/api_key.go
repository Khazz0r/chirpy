package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authParts := strings.Fields(headers.Get("Authorization"))
	if len(authParts) != 2 || strings.ToLower(authParts[0]) != "apikey" {
		return "", errors.New("invalid or missing Authorization header")
	}

	return authParts[1], nil
}
