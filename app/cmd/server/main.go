package main

import (
	"log"
	"net/http"
	"time"
)

const address = ":8080"

func main() {
	server := &http.Server{
		Addr:              address,
		Handler:           http.NewServeMux(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("http-server-projeto-korp listening on %s", address)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
