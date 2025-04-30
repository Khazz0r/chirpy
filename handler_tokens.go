package main

import (
	"net/http"
	"time"

	"github.com/Khazz0r/chirpy/internal/auth"
)

type RefreshToken struct {
	Token string `json:"token"`
}

// handler that will refresh a refresh token to a new one as long as it exists and is valid
func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, req *http.Request) {
	type response struct {
		RefreshToken
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to get bearer token from Authorization header", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token is expired or doesn't exist", err)
		return
	}

	accessToken, err := auth.MakeJWTToken(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		RefreshToken{
			Token: accessToken,
		},
	})
}

// handler that will revoke a token if it is past its expiry time
func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to get bearer token from Authorization header", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to revoke refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
