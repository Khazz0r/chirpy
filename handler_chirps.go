package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
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
		respondWithError(w, http.StatusUnauthorized, "Unable to get bearer token from Authorization header", err)
		return
	}

	// validate to ensure an access token matches
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
	var (
		chirps []database.Chirp
		err    error
		userID uuid.UUID
	)

	authorID := req.URL.Query().Get("author_id")
	sortType := req.URL.Query().Get("sort")

	if authorID == "" {
		chirps, err = cfg.db.GetAllChirps(req.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error retrieving all chirps from database", err)
			return
		}
	} else {
		userID, err = uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to parse author ID", err)
			return
		}

		chirps, err = cfg.db.GetChirpsByAuthorID(req.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Unable to find chirpys by that author ID", err)
			return
		}
	}

	if sortType == "desc" {
		slices.SortFunc(chirps, func(chirp1, chirp2 database.Chirp) int {
			if chirp1.CreatedAt.After(chirp2.CreatedAt) {
				return -1
			} else if chirp1.CreatedAt.Before(chirp2.CreatedAt) {
				return 1
			}
			return 0
		})
	} else {
		// by default, sort in ascending order if desc is not specified
		slices.SortFunc(chirps, func(chirp1, chirp2 database.Chirp) int {
			if chirp1.CreatedAt.Before(chirp2.CreatedAt) {
				return -1
			} else if chirp1.CreatedAt.After(chirp2.CreatedAt) {
				return 1
			}
			return 0
		})
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

// handler that deletes a chirp as long as the user is authorized
func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
	chirpIDStr := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID format", err)
		return
	}

	// obtain token for verifying if user is authorized
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to get bearer token from Authorization header", err)
		return
	}

	// validate to ensure an access token matches
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Not authorized to create a chirp", err)
		return
	}

	chirp, err := cfg.db.GetChirpByID(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "Not author of chirp, can't delete", err)
	}

	err = cfg.db.DeleteChirp(req.Context(), database.DeleteChirpParams{
		ID:     chirp.ID,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find chirp by ID provided", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
