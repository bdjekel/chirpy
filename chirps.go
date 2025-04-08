package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID 			uuid.UUID 	`json:"id"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
	Body		string		`json:"body"`
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	// Request struct
	type parameters struct {
		Body string `json:"body"`
	}

	// Response struct
	type response struct {
		Chirp string `json:"cleaned_body`
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
	chirp := profaneWordHandler(params.Body)

	respondWithJSON(w, http.StatusOK, response{Chirp: chirp})
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