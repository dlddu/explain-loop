// Package server wires the HTTP routes.
package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/dlddu/explain-loop/internal/mcp"
	"github.com/dlddu/explain-loop/internal/version"
)

type health struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// App builds the HTTP handler.
//
// 주의: /mcp 는 아직 인증이 없다. 공개 주소(explain-loop.${PUBLIC_DOMAIN})로 노출되므로
// 데이터를 읽고 쓰는 도구를 붙이기 전에 homelab-k3s-mcp 처럼 OAuth Bearer 검증을 먼저 넣는다.
func App() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", root)
	mux.HandleFunc("GET /healthz", healthz)
	mux.HandleFunc("GET /readyz", readyz)
	mux.Handle("POST /mcp", mcp.NewHandler())
	return logging(mux)
}

func root(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(version.Name))
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, health{Status: "ok", Version: version.Version})
}

func readyz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, health{Status: "ready", Version: version.Version})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		if r.URL.Path == "/healthz" || r.URL.Path == "/readyz" {
			return
		}
		slog.Info("request", "method", r.Method, "path", r.URL.Path, "status", rec.status, "duration", time.Since(start))
	})
}
