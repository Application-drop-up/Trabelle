package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Application-drop-up/Travellle/internal/router"
	"github.com/Application-drop-up/Travellle/internal/testutil"
)

func TestLoginHandler(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"})

	registerUser := func(t *testing.T, email, password string) userResponse {
		t.Helper()
		raw, _ := json.Marshal(map[string]string{
			"email":    email,
			"password": password,
			"name":     "Taro",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("registration status = %d, want %d, body: %s", w.Code, http.StatusCreated, w.Body.String())
		}
		var created userResponse
		if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM users WHERE id = $1", created.ID) })
		return created
	}

	otpCodeFor := func(t *testing.T, userID string) string {
		t.Helper()
		var code string
		err := db.QueryRow(
			"SELECT code FROM login_otps WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1",
			userID,
		).Scan(&code)
		if err != nil {
			t.Fatalf("failed to read otp code: %v", err)
		}
		return code
	}

	loginStartReq := func(body any) *httptest.ResponseRecorder {
		raw, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	loginVerifyReq := func(body any) *httptest.ResponseRecorder {
		raw, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login/verify", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	t.Run("LoginStart sends an otp for valid credentials", func(t *testing.T) {
		t.Parallel()

		created := registerUser(t, "login-start@example.com", "password123")

		w := loginStartReq(map[string]string{"email": created.Email, "password": "password123"})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d, body: %s", w.Code, http.StatusOK, w.Body.String())
		}
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM login_otps WHERE user_id = $1", created.ID) })

		if code := otpCodeFor(t, created.ID); code == "" {
			t.Error("expected an otp to be created")
		}
	})

	t.Run("LoginStart returns 401 for an unknown email", func(t *testing.T) {
		t.Parallel()

		w := loginStartReq(map[string]string{"email": "unknown@example.com", "password": "password123"})
		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})

	t.Run("LoginStart returns 401 for a wrong password", func(t *testing.T) {
		t.Parallel()

		created := registerUser(t, "login-wrong-password@example.com", "password123")

		w := loginStartReq(map[string]string{"email": created.Email, "password": "wrong-password"})
		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})

	t.Run("LoginStart returns 400 for a missing field", func(t *testing.T) {
		t.Parallel()

		w := loginStartReq(map[string]string{"email": "missing-password@example.com", "password": ""})
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("LoginVerify issues a session cookie for a valid code", func(t *testing.T) {
		t.Parallel()

		created := registerUser(t, "login-verify@example.com", "password123")
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM login_otps WHERE user_id = $1", created.ID) })
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM sessions WHERE user_id = $1", created.ID) })

		start := loginStartReq(map[string]string{"email": created.Email, "password": "password123"})
		if start.Code != http.StatusOK {
			t.Fatalf("login start status = %d, want %d, body: %s", start.Code, http.StatusOK, start.Body.String())
		}
		code := otpCodeFor(t, created.ID)

		w := loginVerifyReq(map[string]string{"email": created.Email, "code": code})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d, body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		cookies := w.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session_token" {
				sessionCookie = cookie
			}
		}
		if sessionCookie == nil {
			t.Fatal("expected a session_token cookie to be set")
		}
		if !sessionCookie.HttpOnly {
			t.Error("session_token cookie should be HttpOnly")
		}
		if sessionCookie.Value == "" {
			t.Error("session_token cookie value is empty")
		}
	})

	t.Run("LoginVerify returns 401 for an unknown email", func(t *testing.T) {
		t.Parallel()

		w := loginVerifyReq(map[string]string{"email": "unknown@example.com", "code": "123456"})
		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})

	t.Run("LoginVerify returns 401 for a wrong code", func(t *testing.T) {
		t.Parallel()

		created := registerUser(t, "login-verify-wrong-code@example.com", "password123")
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM login_otps WHERE user_id = $1", created.ID) })

		start := loginStartReq(map[string]string{"email": created.Email, "password": "password123"})
		if start.Code != http.StatusOK {
			t.Fatalf("login start status = %d, want %d, body: %s", start.Code, http.StatusOK, start.Body.String())
		}

		w := loginVerifyReq(map[string]string{"email": created.Email, "code": "000000"})
		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})

	t.Run("LoginVerify returns 400 for a missing field", func(t *testing.T) {
		t.Parallel()

		w := loginVerifyReq(map[string]string{"email": "missing-code@example.com", "code": ""})
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("LoginVerify returns 400 for an invalid request body", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login/verify", bytes.NewReader([]byte("not json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}
