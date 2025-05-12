package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bdjekel/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password 	string `json:"password"`
		Email 		string `json:"email"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}

	type LoginResponse struct {
		ID 			uuid.UUID	`json:"id"` 
		CreatedAt 	time.Time	`json:"created_at"`
		UpdatedAt 	time.Time	`json:"updated_at"`
		Email     	string 		`json:"email"`
		Token		string		`json:"token"`
	}

	// Decode request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}	

	user, err := cfg.DB.UserLogin(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Email is not associated with an account.", err)
		return
	}
	if err := auth.CheckPasswordHash(user.HashedPassword, params.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect Password.", err)
	}

	var expiresIn time.Duration
	switch {
	case params.ExpiresInSeconds < 3600 && params.ExpiresInSeconds > 0:
		expiresIn = time.Duration(params.ExpiresInSeconds) * time.Second
	default:
		expiresIn = time.Duration(3600) * time.Second
	}

	fmt.Printf("\nToken about to be made. Expires in %s seconds.\n", expiresIn)
	
	token, err := auth.MakeJWT(user.ID, os.Getenv("SECRET"), expiresIn)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error creating auth token.", err)
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		ID:        	user.ID,
		CreatedAt: 	user.CreatedAt,
		UpdatedAt: 	user.UpdatedAt,
		Email:     	user.Email,
		Token:		token,
	})
}