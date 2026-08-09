package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	domain "github.com/Application-drop-up/Travellle/internal/domain/user"
	useruc "github.com/Application-drop-up/Travellle/internal/usecase/user"
)

type AuthHandler struct {
	uc *useruc.UseCase
}

func NewAuthHandler(uc *useruc.UseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type userResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

func toUserResponse(dto *useruc.UserDto) userResponse {
	return userResponse{
		ID:        dto.ID.String(),
		Email:     dto.Email,
		Name:      dto.Name,
		CreatedAt: dto.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func (ah *AuthHandler) Register(rw http.ResponseWriter, req *http.Request) {
	var body registerRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Email == "" || body.Password == "" || body.Name == "" {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	dto, err := ah.uc.Register(req.Context(), useruc.RegisterCommand{
		Email:    body.Email,
		Password: body.Password,
		Name:     body.Name,
	})
	if errors.Is(err, domain.ErrEmailTaken) {
		writeError(rw, http.StatusConflict, "email already taken")
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(rw, http.StatusCreated, toUserResponse(dto))
}
