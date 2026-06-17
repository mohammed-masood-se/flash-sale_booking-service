package util

import (
	"booking-service/internal/core/domain"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, statusCode int, body any) {

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(bodyBytes)
}

func BadRequest(w http.ResponseWriter, error string) {
	WriteJSON(w, http.StatusBadRequest, struct {
		Error string `json:"error"`
	}{
		Error: error,
	})
}

func InternalServerError(w http.ResponseWriter) {
	WriteJSON(w, http.StatusInternalServerError, struct {
		Error string `json:"error"`
	}{
		Error: "internal server error",
	})
}

type ClientErrorHandler func(clientError error)

func HandleError(w http.ResponseWriter, err error, fn ClientErrorHandler) {
	if err != nil {

		if domain.IsClientError(err) {
			fn(err)
			return
		}

		if domain.IsServiceError(err) {
			HandleServiceError(w, err)
			return
		}

		HandleUnhandledErrors(w, err)
		return
	}
}

func HandleServiceError(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) {
		w.WriteHeader(499)
		return
	}

	if errors.Is(err, context.DeadlineExceeded) {
		WriteJSON(w, http.StatusServiceUnavailable, struct{ Error string }{Error: "time limit exceeded, please try again"})
		return
	}

	log.Printf("[rest-server] SERVER-ERROR: %v", err)
	InternalServerError(w)
}

func HandleUnhandledErrors(w http.ResponseWriter, err error) {
	log.Printf("[rest-server] UNHANDLED-ERROR: %v", err)
	InternalServerError(w)
}
