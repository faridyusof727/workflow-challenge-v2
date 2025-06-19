package render

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
)

func JSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

func Error(w http.ResponseWriter, r *http.Request, status int, err error, log *slog.Logger) {
	if log != nil {
		log.Error("API Error",
			"endpoint", r.URL.Path,
			"method", r.Method,
			"status", status,
		)
	}

	errorJSON := map[string]string{"error": GetAPIError(err).Error()}
	JSON(w, r, status, errorJSON)
}
