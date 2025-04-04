package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {

	mux := http.NewServeMux()
	hitTracker := apiConfig{}
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080"}

	fileServerWithHitTracker := hitTracker.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("/app/", fileServerWithHitTracker)

	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("POST /api/validate_chirp", chirpValidator)
	mux.HandleFunc("GET /admin/metrics", hitTracker.hitCountHandler)
	mux.HandleFunc("POST /admin/reset", hitTracker.hitResetHandler)

	log.Fatal(server.ListenAndServe())
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

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

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func chirpValidator(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	
	type validResponse struct {
		Valid bool `json:"valid"`
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	// Decode request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	// Handle decoding error
	if err != nil {

		errResp := errorResponse{
			Error: "Something went wrong",
		}
		res, err := json.Marshal(errResp)
		if err != nil {
			log.Printf("Something went wrong: %s", err)
			w.WriteHeader(400)
			return
		}
		// Set headers and write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400) // Or appropriate status code
		w.Write(res)
		return
	}
	// Handle too long chrip
	if len(params.Body) > 140 {
		errResp := errorResponse{
			Error: "Chirp is too long",
		}
		
		res, err := json.Marshal(errResp)
		if err != nil {
			log.Printf("Something went wrong: %s", err)
			w.WriteHeader(400)
			return
		}
		// Set headers and write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400) // Or appropriate status code
		w.Write(res)
		return
	}

	validResp := validResponse{
		Valid: true,
	}

	res, err := json.Marshal(validResp)
	if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(400)
			return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(res)
}
