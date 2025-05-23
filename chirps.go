package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/bdjekel/chirpy/internal/auth"
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
	type parameters struct {
		Body string `json:"body"`
	}

	// Decode request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// ValidateJWT
	access_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retreiving access_token.", err)
		return
	}

	userID, err := auth.ValidateJWT(access_token, os.Getenv("SECRET"))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating access_token.", err)
		return
	}

	// Handle too long chrip
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Max Chirp length exceeded.", nil)
		return
	}

	// Handle Profanity
	params.Body = profaneWordHandler(params.Body)

	// Add chirp to database
	chirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: params.Body,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp", err)
		return
	}

	//Respond with JSON
	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	s := r.URL.Query().Get("author_id")
	sortBy := r.URL.Query().Get("sort")

	if s != "" {
		authorID, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Author ID invalid.", err)
		}

		chirps, err := cfg.DB.GetChirpsByAuthor(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error retrieving chirps", err)
		}

		if sortBy == "desc" {
			sort.Slice(chirps, func (i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt)})
		}

		respondWithJSON(w, http.StatusOK, chirps)		
		return
	}	

	chirps, err := cfg.DB.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving chirps", err)
	}
	
	if sortBy == "desc" {
		sort.Slice(chirps, func (i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt)})
	}
	
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	
	chirpID, err := uuid.Parse(r.PathValue("id"))
    if err != nil {
        http.Error(w, "Invalid chirp ID", http.StatusBadRequest)
        return
    }

	chirp, err := cfg.DB.GetChirpByID(r.Context(), chirpID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Chirp does not exist.", err)
			return
		}

	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
    if err != nil {
        http.Error(w, "Invalid chirp ID", http.StatusForbidden)
        return
    }

	// Validate JWT Access Token
	access_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error retreiving access_token.", err)
		return
	}

	userID, err := auth.ValidateJWT(access_token, os.Getenv("SECRET"))
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Error validating access_token.", err)
		return
	}

	// Find Chirp in Database
	chirp_data, err := cfg.DB.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error retrieving chirp.", err)
		return
	}

	if chirp_data.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Unauthorized DELETE Request.", err)
		return
	}

	if err := cfg.DB.DeleteChirp(r.Context(), chirpID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
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