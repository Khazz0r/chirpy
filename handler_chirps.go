package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Khazz0r/chirpy/internal/auth"
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

type response struct {
	Chirp
}

// handler to create a chirp to the database
func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	// obtain token for verifying if user is authorized
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get bearer token from Authorization header", err)
		return
	}

	// validate to ensure an access token is used
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not authorized to create a chirp", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	// ensure chirp body fits all rules before creating it
	chirpBody, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   chirpBody,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp", err)
		return
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

// helper function for creating chirps to ensure Chirps are valid
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

// handler that gets all Chirps if requested
func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving all chirps from database", err)
		return
	}

	structuredChirps := []Chirp{}

	for _, chirp := range chirps {
		structuredChirps = append(structuredChirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, structuredChirps)
}

// handler that retrieves a Chirp based on the ID passed in through the request
func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {
	chirpIDStr := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID format", err)
		return
	}
	chirp, err := cfg.db.GetChirpByID(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	type response struct {
		Chirp
	}

	respondWithJSON(w, http.StatusOK, response{
		Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		},
	})
}
