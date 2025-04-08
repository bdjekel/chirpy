package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/bdjekel/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID 			uuid.UUID 	`json:"id"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
	Body		string		`json:"body"`
	UserID		string		`json:"user_id"`
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	// Request struct
	type parameters struct {
		Body string `json:"body"`
		ID string `json:"user_id"`
	}

	// Decode request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Handle too long chrip
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Max Chirp length exceeded.", nil)
		return
	}

	// Handle Profanity
	chirpText := profaneWordHandler(params.Body)

	// Add chirp to database
	chirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: chirpText,
		UserID: sql.NullString{String: params.ID, Valid: params.ID != ""},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp: %s", err)
		return
	}

	//Respond with JSON
	respondWithJSON(w, http.StatusOK, chirp)
}

func profaneWordHandler(body string) string {
	// use map so that lookup is O(1)
	profanities := map[string]struct{}{	
		"kerfuffle": {}, 
		"sharbert": {}, 
		"fornax": {},
	}
	words := strings.Fields(body)

	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := profanities[loweredWord]; ok {
			words[i] = "****"
		}
	}

	cleanedBody := strings.Join(words, " ")

	return cleanedBody
}