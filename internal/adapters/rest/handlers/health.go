package handlers

import (
	"booking-service/internal/adapters/rest/util"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (handler *HealthHandler) GetRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/", handler.Health)
	return router
}

func (handler *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	util.WriteJSON(w, http.StatusOK, struct {
		RestServerStatus string `json:"restServerStatus"`
	}{
		RestServerStatus: "online",
	})
}
