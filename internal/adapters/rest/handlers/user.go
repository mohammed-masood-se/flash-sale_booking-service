package handlers

import (
	"booking-service/internal/adapters/rest/handlers/dto"
	"booking-service/internal/adapters/rest/util"
	"booking-service/internal/core/ports"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (handler *UserHandler) GetRouter() chi.Router {
	router := chi.NewRouter()
	router.Post("/", handler.RegisterUser)
	router.Post("/verify", handler.VerifyUser)
	return router
}

func (handler *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var body dto.RegisterUserRequest
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		util.BadRequest(w, "invalid json body")
		return
	}

	insertedID, err := handler.userService.RegisterUser(ctx, body.Email, body.Password)

	if err != nil {
		util.HandleError(w, err, func(clientError error) {
			util.BadRequest(w, clientError.Error())
		})
		return
	}

	util.WriteJSON(w, http.StatusCreated, dto.RegisterUserResponse{
		RegistrationID: insertedID,
		Message:        fmt.Sprintf("registration code sent to '%v'", body.Email),
	})
}

func (handler *UserHandler) VerifyUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var body dto.VerfiyUserRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		util.BadRequest(w, "invalid json body")
		return
	}

	insertedID, err := handler.userService.VerifyUser(ctx, body.Email, body.Code)

	if err != nil {
		util.HandleError(w, err, func(clientError error) {
			util.BadRequest(w, clientError.Error())
		})
		return
	}

	util.WriteJSON(w, http.StatusCreated, dto.VerfiyUserResponse{
		UserID:  insertedID,
		Message: "your account has been successfully verified",
	})
}
