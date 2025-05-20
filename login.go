package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/bdjekel/chirpy/internal/auth"
	"github.com/bdjekel/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password 	string `json:"password"`
		Email 		string `json:"email"`
	}

	type LoginResponse struct {
		ID 				uuid.UUID	`json:"id"` 
		CreatedAt 		time.Time	`json:"created_at"`
		UpdatedAt 		time.Time	`json:"updated_at"`
		Email     		string 		`json:"email"`
		Token			string		`json:"token"`
		RefreshToken 	string		`json:"refresh_token"`
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

	expiresIn := 3600 * time.Second
	
	access_token, err := auth.MakeJWT(user.ID, os.Getenv("SECRET"), expiresIn)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error creating access_token.", err)
	}


	refresh_token_string, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh_token_string", err)
		return
	}

	refresh_token, err := cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refresh_token_string,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh_token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		ID:        	user.ID,
		CreatedAt: 	user.CreatedAt,
		UpdatedAt: 	user.UpdatedAt,
		Email:     	user.Email,
		Token:		access_token,
		RefreshToken: refresh_token.Token,
	})
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	// ValidateJWT
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retreiving refresh_token.", err)
		return
	}
	
	refresh_token_data, err := cfg.DB.GetRefreshToken(r.Context(), refresh_token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token is expired or does not exist.", err)
	}
	
	respondWithJSON(w, http.StatusOK, refresh_token_data.Token)

}