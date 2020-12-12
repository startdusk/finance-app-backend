package utils

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// GenericError - represent error structure for generic error
// (we need this to make all our error response the same)
type GenericError struct {
	Code  int         `json:"code"`
	Error string      `json:"error"`
	Data  interface{} `json:"data,omitempty"`
}

// WriteError returns a JSON error message and HTTP status code
func WriteError(w http.ResponseWriter, code int, message string, data interface{}) {
	response := GenericError{
		Code:  code,
		Error: message,
		Data:  data,
	}

	WriteJSON(w, code, response)
}

// WriteJSON returns a JSON data and HTTP status code
func WriteJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logrus.WithError(err).Warn("Error writing response.")
	}
}
