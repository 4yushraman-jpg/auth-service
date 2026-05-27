package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/4yushraman-jpg/auth-service/internal/dto"
	appErrors "github.com/4yushraman-jpg/auth-service/internal/errors"
	"github.com/4yushraman-jpg/auth-service/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(
	authService *service.AuthService,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

func (h *AuthHandler) Register(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.validate.Struct(request); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	err := h.authService.Register(
		r.Context(),
		request.Email,
		request.Password,
	)

	if err != nil {
		if errors.Is(err, appErrors.ErrUserAlreadyExists) {
			http.Error(
				w,
				err.Error(),
				http.StatusConflict,
			)
			return
		}

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)

		return
	}

	response := dto.AuthResponse{
		Message: "user registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}
