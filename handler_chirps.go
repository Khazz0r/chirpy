package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"errors"

	"github.com/Khazz0r/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

// handler to create a chirp to the database
func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
	}

	// ensure chirp body fits all rules before creating it
	chirpBody, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
	}

	type response struct {
		Chirp
	}

	chirp, err := cfg.db.CreateChirp(context.Background(), database.CreateChirpParams{
		Body:   chirpBody,
		UserID: UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp", err)
	}

	respondWithJSON(w, http.StatusCreated, response{
		Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		},
	})
}

// helper function to ensure Chirps are valid under the rules when a chirp is created
func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	cleanedBody := badWordReplacer(body)

	return cleanedBody, nil
}

// helper function for validate chirp to clean up "bad" words
func badWordReplacer(body string) string {
	sliceBody := strings.Split(body, " ")

	for i, word := range sliceBody {
		lowerWord := strings.ToLower(word)
		if lowerWord == "kerfuffle" || lowerWord == "sharbert" || lowerWord == "fornax" {
			word = "****"
			sliceBody[i] = word
		}
	}

	cleanedBody := strings.Join(sliceBody, " ")

	return cleanedBody
}
