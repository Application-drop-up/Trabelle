package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	domain "github.com/Application-drop-up/Travellle/internal/domain/user"
	useruc "github.com/Application-drop-up/Travellle/internal/usecase/user"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AuthHandler struct {
	useCase *useruc.UseCase
}

func NewAuthHandler(useCase *useruc.UseCase) *AuthHandler {
	return &AuthHandler{useCase: useCase}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type updateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
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

func (authHandler *AuthHandler) Register(rw http.ResponseWriter, req *http.Request) {
	var body registerRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Email == "" || body.Password == "" || body.Name == "" {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	dto, err := authHandler.useCase.Register(req.Context(), useruc.RegisterCommand{
		Email:    body.Email,
		Password: body.Password,
		Name:     body.Name,
	})
	if errors.Is(err, domain.ErrEmailTaken) {
		writeError(rw, http.StatusConflict, "email already taken")
		return
	}
	if errors.Is(err, domain.ErrInvalidEmail) || errors.Is(err, domain.ErrPasswordTooShort) {
		writeError(rw, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(rw, http.StatusCreated, toUserResponse(dto))
}

func (authHandler *AuthHandler) GetByID(rw http.ResponseWriter, req *http.Request) {
	id, err := uuid.Parse(chi.URLParam(req, "id"))
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid id")
		return
	}

	dto, err := authHandler.useCase.GetUserByID(req.Context(), id)
	if errors.Is(err, domain.ErrNotFound) {
		writeError(rw, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(rw, http.StatusOK, toUserResponse(dto))
}

func (authHandler *AuthHandler) Update(rw http.ResponseWriter, req *http.Request) {
	id, err := uuid.Parse(chi.URLParam(req, "id"))
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid id")
		return
	}

	var body updateRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Name == "" || body.Email == "" {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	dto, err := authHandler.useCase.UpdateUser(req.Context(), id, useruc.UpdateCommand{
		Name:  body.Name,
		Email: body.Email,
	})
	if errors.Is(err, domain.ErrNotFound) {
		writeError(rw, http.StatusNotFound, "user not found")
		return
	}
	if errors.Is(err, domain.ErrEmailTaken) {
		writeError(rw, http.StatusConflict, "email already taken")
		return
	}
	if errors.Is(err, domain.ErrInvalidEmail) {
		writeError(rw, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(rw, http.StatusOK, toUserResponse(dto))
}
