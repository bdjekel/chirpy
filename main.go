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

	type error struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		// How do I return the error correctly here? Should I return an error anywhere else? REad lesson and docs
		w.WriteHeader(500) // confirm error code is correct
	}

	// How do I return struct here properly? Review lesson, read docs
}

func handler(w http.ResponseWriter, r *http.Request){
    type parameters struct {
        // these tags indicate how the keys in the JSON should be mapped to the struct fields
        // the struct fields must be exported (start with a capital letter) if you want them parsed
        Name string `json:"name"`
        Age int `json:"age"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        // an error will be thrown if the JSON is invalid or has the wrong types
        // any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }
    // params is a struct with data populated successfully
    // ...
}