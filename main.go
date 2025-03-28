package main

import (
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

	mux.HandleFunc("/healthz", readinessHandler)
	mux.HandleFunc("/metrics", hitTracker.hitCountHandler)
	mux.HandleFunc("/reset", hitTracker.hitResetHandler)

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

func (cfg *apiConfig) hitCountHandler(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileserverHits.Load()
	w.Write([]byte(fmt.Sprintf("Hits: %v", hits)))
}

func (cfg *apiConfig) hitResetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
