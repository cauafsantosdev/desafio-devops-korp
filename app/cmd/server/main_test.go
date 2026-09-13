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

func TestHealthEndpoint(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recorder := httptest.NewRecorder()

	newMux().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response healthResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != "ok" {
		t.Fatalf("expected health status %q, got %q", "ok", response.Status)
	}
}

func TestMetricsEndpoint(t *testing.T) {
	mux := newMux()

	request := httptest.NewRequest(http.MethodGet, "/projeto-korp", nil)
	mux.ServeHTTP(httptest.NewRecorder(), request)

	metricsRequest := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metricsRecorder := httptest.NewRecorder()
	mux.ServeHTTP(metricsRecorder, metricsRequest)

	if metricsRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, metricsRecorder.Code)
	}

	body := metricsRecorder.Body.String()

	for _, metric := range []string{
		"http_server_requests_total",
		"http_server_request_duration_seconds",
	} {
		if !strings.Contains(body, metric) {
			t.Errorf("expected metrics response to contain %q", metric)
		}
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
