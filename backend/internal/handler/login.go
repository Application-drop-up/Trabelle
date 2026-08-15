package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	domain "github.com/Application-drop-up/Travellle/internal/domain/user"
	useruc "github.com/Application-drop-up/Travellle/internal/usecase/user"
)

const sessionCookieName = "session_token"

type loginStartRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginVerifyRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type messageResponse struct {
	Message string `json:"message"`
}

func (authHandler *AuthHandler) LoginStart(rw http.ResponseWriter, req *http.Request) {
	var body loginStartRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Email == "" || body.Password == "" {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	err := authHandler.useCase.LoginStart(req.Context(), body.Email, body.Password)
	if errors.Is(err, domain.ErrInvalidCredentials) {
		writeError(rw, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(rw, http.StatusOK, messageResponse{Message: "verification code sent"})
}

func (authHandler *AuthHandler) LoginVerify(rw http.ResponseWriter, req *http.Request) {
	var body loginVerifyRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Email == "" || body.Code == "" {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	session, err := authHandler.useCase.LoginVerify(req.Context(), body.Email, body.Code)
	if errors.Is(err, domain.ErrInvalidCredentials) || errors.Is(err, useruc.ErrInvalidOTP) {
		writeError(rw, http.StatusUnauthorized, "invalid email or code")
		return
	}
	if errors.Is(err, useruc.ErrOTPExpired) {
		writeError(rw, http.StatusUnauthorized, "verification code expired")
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     sessionCookieName,
		Value:    session.Token,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	writeJSON(rw, http.StatusOK, messageResponse{Message: "login successful"})
}
