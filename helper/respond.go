package helper

import (
	"encoding/json"
	"net/http"
	"time"
)

// Helper functions for HTTP responses
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, map[string]interface{}{
		"error":     message,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
