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
	// Code is only populated in non-production environments.
	Code string `json:"code,omitempty"`
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

	code, err := authHandler.useCase.LoginStart(req.Context(), body.Email, body.Password)
	if errors.Is(err, domain.ErrInvalidCredentials) {
		writeError(rw, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	response := messageResponse{Message: "verification code sent"}
	if authHandler.isDev {
		response.Code = code
	}
	writeJSON(rw, http.StatusOK, response)
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

func (authHandler *AuthHandler) Me(rw http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(sessionCookieName)
	if err != nil {
		writeError(rw, http.StatusUnauthorized, "not authenticated")
		return
	}

	dto, err := authHandler.useCase.CurrentUser(req.Context(), cookie.Value)
	if errors.Is(err, useruc.ErrSessionNotFound) {
		writeError(rw, http.StatusUnauthorized, "not authenticated")
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(rw, http.StatusOK, toUserResponse(dto))
}

func (authHandler *AuthHandler) Logout(rw http.ResponseWriter, req *http.Request) {
	if cookie, err := req.Cookie(sessionCookieName); err == nil {
		if err := authHandler.useCase.Logout(req.Context(), cookie.Value); err != nil {
			writeError(rw, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	rw.WriteHeader(http.StatusNoContent)
}
