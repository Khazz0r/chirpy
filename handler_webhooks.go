package main

import (
	"encoding/json"
	"net/http"

	"github.com/Khazz0r/chirpy/internal/auth"
	"github.com/google/uuid"
)

type Data struct {
	UserID string `json:"user_id"`
}

// handler that upgrades user to Chirpy Red when it receives a request from the Polka webhook
func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  `json:"data"`
	}

	polkaAPIKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not get API key from header", err)
		return
	}

	if polkaAPIKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Polka API key does not match", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing user ID", err)
		return
	}

	_, err = cfg.db.UpgradeUser(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Unable to find user by ID", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
