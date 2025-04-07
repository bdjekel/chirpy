package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/bdjekel/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Connection to database could not be established")
	}

	dbQueries := database.New(db)


	mux := http.NewServeMux()
	hitTracker := apiConfig{
		DB: *dbQueries,
	}
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080"}

	fileServerWithHitTracker := hitTracker.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("/app/", fileServerWithHitTracker)

	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)
	mux.HandleFunc("GET /admin/metrics", hitTracker.hitCountHandler)
	mux.HandleFunc("POST /admin/reset", hitTracker.hitResetHandler)

	log.Fatal(server.ListenAndServe())
}



// STRUCTS

type apiConfig struct {
	fileserverHits atomic.Int32
	DB database.Queries
}

type parameters struct {
	Body string `json:"body"`
}

type ValidResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// FUNCTIONS

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// TODO: what is best practice? It seems the returned html should be pulled from a file instead of being sprinted in...
func (cfg *apiConfig) hitCountHandler(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileserverHits.Load()
	r.Header.Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", hits)))
}

func (cfg *apiConfig) hitResetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}


// HANDLER FUNCTIONS

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func validateChirp(w http.ResponseWriter, r *http.Request) {

	// Decode request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Something went wrong: %s", err)
		w.WriteHeader(400)
		return
	}

	// Handle too long chrip
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Max Chirp length exceeded.")
		return
	}

	// Handle Profanity
	res := ValidResponse{CleanedBody: profaneWordHandler(params.Body)}

	respondWithJSON(w, 200, res)
}


// HELPER FUNCTIONS

func respondWithError(w http.ResponseWriter, code int, msg string) {

	// Encode response
	res, err := json.Marshal(errorResponse{Error: msg})
	if err != nil {
		log.Printf("Something went wrong: %s", err)
		w.WriteHeader(400)
		return
	}
	
	// Set headers and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	// Encode response
	res, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Something went wrong: %s", err)
		return
	}

	// Set headers and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}

func profaneWordHandler(body string) string {
	cleaned_words := ""
	profanities := []string{	
		"kerfuffle", 
		"sharbert", 
		"fornax"}
	words := strings.Fields(body)

	for _, word := range(words) {
		for _, profanity := range(profanities) {
			if strings.ToLower(word) == profanity {
				word = "****"
				break
			}
		}
		cleaned_words += word + " "
	}

	return strings.TrimSpace(cleaned_words)
}