package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const address = ":8080"

type projetoKorpResponse struct {
	Nome    string `json:"nome"`
	Horario string `json:"horario"`
}

type healthResponse struct {
	Status string `json:"status"`
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle(
		"GET /projeto-korp",
		observeHTTP("/projeto-korp", http.HandlerFunc(projetoKorpHandler)),
	)
	mux.Handle(
		"GET /healthz",
		observeHTTP("/healthz", http.HandlerFunc(healthHandler)),
	)
	mux.Handle("GET /metrics", metricsHandler())

	return mux
}

func projetoKorpHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := projetoKorpResponse{
		Nome:    "Projeto Korp",
		Horario: time.Now().UTC().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(healthResponse{Status: "ok"}); err != nil {
		log.Printf("failed to encode health response: %v", err)
	}
}

func main() {
	server := &http.Server{
		Addr:              address,
		Handler:           newMux(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("http-server-projeto-korp listening on %s", address)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
