package main

//TODO: Clean up main() by moving certain functions into their own files.

import (
	"database/sql"
	"encoding/json" // move
	"fmt"           // needed?
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time" // move

	"github.com/bdjekel/chirpy/internal/database"
	"github.com/google/uuid" // move
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	DB database.Queries
	platform string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	// Retrieve env variables
	godotenv.Load()
	
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	// Connect to database
	dbConnection, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Connection to database could not be established: %s", err)
	}

	dbQueries := database.New(dbConnection)


	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		DB: *dbQueries,
		platform: platform,
	}

	mux := http.NewServeMux()

	// Middleware
	fileServerWithHitTracker := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fileServerWithHitTracker)

	// api endpoints
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirps)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	
	// Admin endpoints
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}



// STRUCTS

type chirpParameters struct {
	Body string `json:"body"`
}

type ValidResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type UserRequest struct {
	Email string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

// FUNCTIONS

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// TODO: what is best practice? It seems the returned html should be pulled from a file instead of being sprinted in...
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	// api endpoints
	hits := cfg.fileserverHits.Load()
	r.Header.Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", hits)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != os.Getenv("PLATFORM") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(403)
		w.Write([]byte("Forbidden"))
		return
	}
	cfg.fileserverHits.Store(0)
	cfg.DB.DeleteAllUsers(r.Context())
}


func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := UserRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Something went wrong: %s", err)
		w.WriteHeader(400)
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Something went wrong: %s", err)
		w.WriteHeader(400)
		return
	}

	responseUser := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	respondWithJSON(w, 201, responseUser)
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}


//TODO: Move respondWithError and respond withJSON to their own file.

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {

	// Encode response
	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 500-level error: %s", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errorResponse{ Error: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	// Encode response
	res, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(res)
}