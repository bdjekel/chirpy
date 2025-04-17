package main

import (
	"encoding/json"
	"net/http"
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
		Email: params.Email, 
		HashedPassword: hashedPassword,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}


	respondWithJSON(w, 201, UserCreatedResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
