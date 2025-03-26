package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080"}

	mux.Handle("/assets/logo.png", http.FileServer(http.Dir(".")))

	log.Fatal(server.ListenAndServe())
}
