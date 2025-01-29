package handler

import (
	"encoding/json"
	"net/http"

	"github.com/anilsoylu/answer-backend/internal/models"
	"github.com/anilsoylu/answer-backend/internal/service"
	"github.com/anilsoylu/answer-backend/internal/utils/errors"
	"github.com/anilsoylu/answer-backend/internal/utils/response"
	"github.com/anilsoylu/answer-backend/internal/utils/validator"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Register handles user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.ErrInvalidInput)
		return
	}

	// Validate request
	if err := validator.ValidateStruct(req); err != nil {
		response.Error(w, err)
		return
	}

	// Register user
	if err := h.service.Register(&req); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{
		"message": "User registered successfully",
	})
}

// Login handles user login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.ErrInvalidInput)
		return
	}

	// Validate request
	if err := validator.ValidateStruct(req); err != nil {
		response.Error(w, err)
		return
	}

	// Login user
	authResponse, err := h.service.Login(&req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, authResponse)
} 