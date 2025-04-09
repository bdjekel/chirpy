package main

import (
	"encoding/json"
	"log"
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
	UserID		uuid.UUID	`json:"user_id"`
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {

	// Decode request
	decoder := json.NewDecoder(r.Body)
	params := database.CreateChirpParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	log.Println("---UserID---")
	log.Println(params.UserID)
	log.Println("---Body---")
	log.Println(params.Body)

	// Handle too long chrip
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Max Chirp length exceeded.", nil)
		return
	}

	// Handle Profanity
	params.Body = profaneWordHandler(params.Body)
	log.Println("---Updated Body---")
	log.Println(params.Body)
	// Add chirp to database
	chirp, err := cfg.DB.CreateChirp(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp: %s", err)
		return
	}

	//Respond with JSON
	respondWithJSON(w, http.StatusCreated, chirp)
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