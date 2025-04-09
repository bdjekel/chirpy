package main

//TODO: Clean up main() by moving certain functions into their own files.

import (
	"database/sql" // move
	// needed?
	"log"
	"net/http"
	"os"
	"sync/atomic"

	// move
	"github.com/bdjekel/chirpy/internal/database" // move
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