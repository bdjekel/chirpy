package main

import (
	"net/http"
	"os"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	//TODO: should the below be "dev" or is the os.Getenv() call correct?

	w.Header().Set("Content-Type", "application/json")
	if cfg.platform != os.Getenv("PLATFORM") {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset requires admin privleges"))
		return
	}
	cfg.fileserverHits.Store(0)
	cfg.DB.DeleteAllUsers(r.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0. Database reset to initial state."))
}