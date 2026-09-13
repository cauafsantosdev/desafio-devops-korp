package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_server_requests_total",
			Help: "Total number of HTTP requests processed by the service.",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_server_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(statusCode int) {
	if w.statusCode != 0 {
		return
	}

	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *statusResponseWriter) Write(body []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}

	return w.ResponseWriter.Write(body)
}

func observeHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		startedAt := time.Now()
		recorder := &statusResponseWriter{ResponseWriter: w}

		next.ServeHTTP(recorder, r)

		statusCode := recorder.statusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		httpRequestsTotal.
			WithLabelValues(r.Method, metricPath(r.URL.Path), strconv.Itoa(statusCode)).
			Inc()

		httpRequestDuration.
			WithLabelValues(r.Method, metricPath(r.URL.Path)).
			Observe(time.Since(startedAt).Seconds())
	})
}

func metricPath(path string) string {
	switch path {
	case "/projeto-korp", "/healthz":
		return path
	default:
		return "unmatched"
	}
}

func metricsHandler() http.Handler {
	return promhttp.Handler()
}
