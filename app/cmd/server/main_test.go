package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestProjetoKorpEndpoint(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/projeto-korp", nil)
	recorder := httptest.NewRecorder()

	newMux().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var response projetoKorpResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Nome != "Projeto Korp" {
		t.Errorf("expected nome %q, got %q", "Projeto Korp", response.Nome)
	}

	horario, err := time.Parse(time.RFC3339, response.Horario)
	if err != nil {
		t.Fatalf("horario is not valid RFC3339: %v", err)
	}

	if !strings.HasSuffix(response.Horario, "Z") {
		t.Errorf("expected UTC timestamp ending in Z, got %q", response.Horario)
	}

	if time.Since(horario) > 2*time.Second {
		t.Errorf("expected current timestamp, got %q", response.Horario)
	}
}

func TestProjetoKorpRejectsPost(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/projeto-korp", nil)
	recorder := httptest.NewRecorder()

	newMux().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusMethodNotAllowed,
			recorder.Code,
		)
	}
}
