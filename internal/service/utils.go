package service

import (
	"encoding/json"
	"errors"
	"net/http"
)

const (
	ContentTypeJSON = "application/json"
	ContentType     = "Content-Type"
)

var (
	ErrFailedToEncodeResponse = errors.New("failed to encode response")
	ErrFailedToWriteResponse  = errors.New("failed to write response")
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set(ContentType, ContentTypeJSON)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, ErrFailedToEncodeResponse.Error(), http.StatusInternalServerError)
		return
	}
}

func writeNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func writeString(w http.ResponseWriter, status int, message string) {
	w.Header().Set(ContentType, ContentTypeJSON)
	w.WriteHeader(status)
	if _, err := w.Write([]byte(message)); err != nil {
		http.Error(w, ErrFailedToWriteResponse.Error(), http.StatusInternalServerError)
		return
	}
}
