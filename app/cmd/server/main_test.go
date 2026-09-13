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

	newHandler().ServeHTTP(recorder, request)

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

	newHandler().ServeHTTP(recorder, request)

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
	handler := newHandler()

	requests := []struct {
		method string
		path   string
		status int
	}{
		{http.MethodGet, "/projeto-korp", http.StatusOK},
		{http.MethodPost, "/projeto-korp", http.StatusMethodNotAllowed},
		{http.MethodGet, "/rota-inexistente", http.StatusNotFound},
	}

	for _, testCase := range requests {
		request := httptest.NewRequest(testCase.method, testCase.path, nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)

		if recorder.Code != testCase.status {
			t.Fatalf(
				"%s %s: expected status %d, got %d",
				testCase.method,
				testCase.path,
				testCase.status,
				recorder.Code,
			)
		}
	}

	metricsRequest := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metricsRecorder := httptest.NewRecorder()
	handler.ServeHTTP(metricsRecorder, metricsRequest)

	if metricsRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, metricsRecorder.Code)
	}

	body := metricsRecorder.Body.String()

	for _, expected := range []string{
		"http_server_requests_total",
		"http_server_request_duration_seconds",
		`http_server_requests_total{method="GET",path="/projeto-korp",status="200"}`,
		`http_server_requests_total{method="POST",path="/projeto-korp",status="405"}`,
		`http_server_requests_total{method="GET",path="unmatched",status="404"}`,
	} {
		if !strings.Contains(body, expected) {
			t.Errorf("expected metrics response to contain %q", expected)
		}
	}

	if strings.Contains(body, `path="/metrics"`) {
		t.Error("expected /metrics to be excluded from application request metrics")
	}
}

func TestProjetoKorpRejectsPost(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/projeto-korp", nil)
	recorder := httptest.NewRecorder()

	newHandler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusMethodNotAllowed,
			recorder.Code,
		)
	}
}
