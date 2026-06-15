package rest

import (
	"booking-service/internal/adapters/rest/handlers"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RestServerConfig struct {
	Port string
}

type RestServer struct {
	config *RestServerConfig
	router chi.Router
}

func NewRestServer(config RestServerConfig) (*RestServer, error) {
	healthHandler := handlers.NewHealthHandler()

	router := chi.NewRouter()

	router.Route("/api/v1", func(router chi.Router) {
		router.Mount("/health", healthHandler.GetRouter())
	})

	return &RestServer{
		config: &config,
		router: router,
	}, nil
}

func (restServer *RestServer) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%v", restServer.config.Port), restServer.router)
}
