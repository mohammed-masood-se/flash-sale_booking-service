package rest

import (
	"booking-service/internal/adapters/rest/handlers"
	"booking-service/internal/core/ports"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RestServerConfig struct {
	Port string
}

type Services struct {
	UserService ports.UserService
}

type RestServer struct {
	config *RestServerConfig
	router chi.Router
}

func NewRestServer(config RestServerConfig, services *Services) (*RestServer, error) {
	healthHandler := handlers.NewHealthHandler()

	userHandler := handlers.NewUserHandler(services.UserService)

	router := chi.NewRouter()

	router.Route("/api/v1", func(router chi.Router) {
		router.Mount("/health", healthHandler.GetRouter())
		router.Mount("/users", userHandler.GetRouter())
	})

	return &RestServer{
		config: &config,
		router: router,
	}, nil
}

func (restServer *RestServer) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%v", restServer.config.Port), restServer.router)
}
