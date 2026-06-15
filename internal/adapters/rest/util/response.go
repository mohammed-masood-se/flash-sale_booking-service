package util

import (
	"encoding/json"
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
