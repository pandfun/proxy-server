package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Used to send JSON response
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// Send error messages. Calls @ WriteJSON
func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func WriteResponse(w http.ResponseWriter, resp *http.Response) {

	// Copy the response headers
	for headerKey, headerValue := range resp.Header {
		w.Header()[headerKey] = headerValue
	}

	// Copy the status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	_, err := io.Copy(w, resp.Body)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to write response"))
		return
	}
}