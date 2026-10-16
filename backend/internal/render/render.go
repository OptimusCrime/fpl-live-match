// Package render provides HTTP response helpers.
package render

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/OptimusCrime/fpl-live-match/backend/internal/resterr"
)

// JSON writes v as a JSON response. If v is a resterr.Resterr, the appropriate
// HTTP status code and error message are written instead.
func JSON(w http.ResponseWriter, r *http.Request, v any) {
	w.Header().Set("Content-Type", "application/json")

	if re, ok := v.(resterr.Resterr); ok {
		level := slog.LevelWarn
		if re.StatusCode >= 500 {
			level = slog.LevelError
		}
		slog.Log(r.Context(), level, "request error", append(requestAttrs(r), "status", re.StatusCode, "err", re.Err)...)
		w.WriteHeader(re.StatusCode)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": re.Err.Error()}); err != nil {
			slog.Error("failed to write error response", append(requestAttrs(r), "err", err)...)
		}
		return
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		slog.Error("failed to encode response", append(requestAttrs(r), "err", err)...)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(buf.Bytes()); err != nil {
		slog.Error("failed to write response", append(requestAttrs(r), "err", err)...)
	}
}

// requestAttrs returns log attributes identifying the request.
func requestAttrs(r *http.Request) []any {
	return []any{
		"method", r.Method,
		"path", r.URL.Path,
		"request_id", middleware.GetReqID(r.Context()),
	}
}
