package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bdjekel/chirpy/internal/auth"
	"github.com/bdjekel/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type UserCreatedResponse struct {
		ID 			uuid.UUID	`json:"id"` 
		CreatedAt 	time.Time	`json:"created_at"`
		UpdatedAt 	time.Time	`json:"updated_at"`
		Email     	string 		`json:"email"`
		IsChirpyRed bool		`json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		HashedPassword: hashedPassword,
		Email: params.Email, 
	})
	fmt.Printf(">>%s<<", err)	
	if err != nil {
		// fmt.Printf(">>%s<<", user)
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}


	respondWithJSON(w, 201, UserCreatedResponse{
		ID:        		user.ID,
		CreatedAt: 		user.CreatedAt,
		UpdatedAt: 		user.UpdatedAt,
		Email:     		user.Email,
		IsChirpyRed: 	user.IsChirpyRed,
	})
}

func (cfg *apiConfig) handlerUpdateCredentials (w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		NewEmail string `json:"email"`
		NewPassword string `json:"password"`
	}

	type UserUpdatedResponse struct {
		ID 			uuid.UUID	`json:"id"` 
		CreatedAt 	time.Time	`json:"created_at"`
		UpdatedAt 	time.Time	`json:"updated_at"`
		Email   	string 		`json:"email"`
		IsChirpyRed bool		`json:"is_chirpy_red"`
	}

	// Validate JWT Access Token
	access_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error retreiving access_token.", err)
		return
	}

	userID, err := auth.ValidateJWT(access_token, os.Getenv("SECRET"))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating access_token.", err)
		return
	}

	// Decode Request Payload
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	// Hash New Password
	hashedPassword, err := auth.HashPassword(params.NewPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	// Update User in Database
	user, err := cfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID: userID,
		Email: params.NewEmail, 
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	// Encode Response Payload
	respondWithJSON(w, 200, UserUpdatedResponse{
		ID:        		user.ID,
		CreatedAt: 		user.CreatedAt,
		UpdatedAt: 		user.UpdatedAt,
		Email:     		user.Email,
		IsChirpyRed: 	user.IsChirpyRed,
	})
}

func (cfg *apiConfig) handlerMembershipUpgrade(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	// Validate API Key
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error retreiving apiKey.", err)
		return
	}

	if apiKey != os.Getenv("POLKA_KEY") {
		respondWithJSON(w, http.StatusUnauthorized, "Invalid ApiKey")
		return
	}

	// Decode request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	//TODO: this endpoint should definitely have auth...

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
	}

	err = cfg.DB.UpgradeUserMembership(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User Membership Upgrade Failed.", err)
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
