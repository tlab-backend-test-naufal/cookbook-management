package libhttp

import (
	"encoding/json"
	"net/http"
)

type BaseError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Base is the base object of all responses
type Base struct {
	Message *string      `json:"message,omitempty"`
	Data    *interface{} `json:"data,omitempty"`
	Error   *BaseError   `json:"error,omitempty"`
}

// NoContent sends a response without any content
func NoContent(w http.ResponseWriter) {
	respond(w, http.StatusNoContent, nil)
}

// WithMessage sends a response with a simple text message
func WithMessage(w http.ResponseWriter, code int, message string) {
	respond(w, code, Base{Message: &message})
}

// WithJSON sends a response containing a JSON object
func WithJSON(w http.ResponseWriter, code int, jsonPayload interface{}, message ...string) {
	if len(message) != 0 {
		respond(w, code, Base{Data: &jsonPayload, Message: &message[0]})
		return
	}
	respond(w, code, Base{Data: &jsonPayload})
}

func WithError(w http.ResponseWriter, code int, err error) {
	errMsg := err.Error()
	respond(w, code, Base{Error: &BaseError{Message: errMsg}})
}

func respond(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}
