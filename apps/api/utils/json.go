package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

func WriteError(w http.ResponseWriter, status int, err error) {
	code := http.StatusText(status)
	if apiErr, ok := err.(APIError); ok {
		code = apiErr.Code
	}

	_ = WriteJSON(w, status, APIError{
		Code:    code,
		Message: err.Error(),
	})
}
