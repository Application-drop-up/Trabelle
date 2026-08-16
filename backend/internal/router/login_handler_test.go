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

func TestLoginHandler(test *testing.T) {
	test.Parallel()

	db := testutil.NewTestDB(test)
	mux := router.New(db, "test-api-key", []string{"http://localhost:3000"})

	registerUser := func(test *testing.T, email, password string) userResponse {
		test.Helper()
		raw, _ := json.Marshal(map[string]string{
			"email":    email,
			"password": password,
			"name":     "Taro",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusCreated {
			test.Fatalf("registration status = %d, want %d, body: %s", recorder.Code, http.StatusCreated, recorder.Body.String())
		}
		var created userResponse
		if err := json.Unmarshal(recorder.Body.Bytes(), &created); err != nil {
			test.Fatalf("failed to decode response: %v", err)
		}
		test.Cleanup(func() { _, _ = db.Exec("DELETE FROM users WHERE id = $1", created.ID) })
		return created
	}

	otpCodeFor := func(test *testing.T, userID string) string {
		test.Helper()
		var code string
		err := db.QueryRow(
			"SELECT code FROM login_otps WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1",
			userID,
		).Scan(&code)
		if err != nil {
			test.Fatalf("failed to read otp code: %v", err)
		}
		return code
	}

	loginStartReq := func(body any) *httptest.ResponseRecorder {
		raw, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)
		return recorder
	}

	loginVerifyReq := func(body any) *httptest.ResponseRecorder {
		raw, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login/verify", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)
		return recorder
	}

	test.Run("LoginStart sends an otp for valid credentials", func(test *testing.T) {
		test.Parallel()

		created := registerUser(test, "login-start@example.com", "password123")

		recorder := loginStartReq(map[string]string{"email": created.Email, "password": "password123"})
		if recorder.Code != http.StatusOK {
			test.Fatalf("status = %d, want %d, body: %s", recorder.Code, http.StatusOK, recorder.Body.String())
		}
		test.Cleanup(func() { _, _ = db.Exec("DELETE FROM login_otps WHERE user_id = $1", created.ID) })

		if code := otpCodeFor(test, created.ID); code == "" {
			test.Error("expected an otp to be created")
		}
	})

	test.Run("LoginStart returns 401 for an unknown email", func(test *testing.T) {
		test.Parallel()

		recorder := loginStartReq(map[string]string{"email": "unknown@example.com", "password": "password123"})
		if recorder.Code != http.StatusUnauthorized {
			test.Errorf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
		}
	})

	test.Run("LoginStart returns 401 for a wrong password", func(test *testing.T) {
		test.Parallel()

		created := registerUser(test, "login-wrong-password@example.com", "password123")

		recorder := loginStartReq(map[string]string{"email": created.Email, "password": "wrong-password"})
		if recorder.Code != http.StatusUnauthorized {
			test.Errorf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
		}
	})

	test.Run("LoginStart returns 400 for a missing field", func(test *testing.T) {
		test.Parallel()

		recorder := loginStartReq(map[string]string{"email": "missing-password@example.com", "password": ""})
		if recorder.Code != http.StatusBadRequest {
			test.Errorf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
		}
	})

	test.Run("LoginVerify issues a session cookie for a valid code", func(test *testing.T) {
		test.Parallel()

		created := registerUser(test, "login-verify@example.com", "password123")
		test.Cleanup(func() { _, _ = db.Exec("DELETE FROM login_otps WHERE user_id = $1", created.ID) })
		test.Cleanup(func() { _, _ = db.Exec("DELETE FROM sessions WHERE user_id = $1", created.ID) })

		start := loginStartReq(map[string]string{"email": created.Email, "password": "password123"})
		if start.Code != http.StatusOK {
			test.Fatalf("login start status = %d, want %d, body: %s", start.Code, http.StatusOK, start.Body.String())
		}
		code := otpCodeFor(test, created.ID)

		recorder := loginVerifyReq(map[string]string{"email": created.Email, "code": code})
		if recorder.Code != http.StatusOK {
			test.Fatalf("status = %d, want %d, body: %s", recorder.Code, http.StatusOK, recorder.Body.String())
		}

		cookies := recorder.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session_token" {
				sessionCookie = cookie
			}
		}
		if sessionCookie == nil {
			test.Fatal("expected a session_token cookie to be set")
		}
		if !sessionCookie.HttpOnly {
			test.Error("session_token cookie should be HttpOnly")
		}
		if sessionCookie.Value == "" {
			test.Error("session_token cookie value is empty")
		}
	})

	test.Run("LoginVerify returns 401 for an unknown email", func(test *testing.T) {
		test.Parallel()

		recorder := loginVerifyReq(map[string]string{"email": "unknown@example.com", "code": "123456"})
		if recorder.Code != http.StatusUnauthorized {
			test.Errorf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
		}
	})

	test.Run("LoginVerify returns 401 for a wrong code", func(test *testing.T) {
		test.Parallel()

		created := registerUser(test, "login-verify-wrong-code@example.com", "password123")
		test.Cleanup(func() { _, _ = db.Exec("DELETE FROM login_otps WHERE user_id = $1", created.ID) })

		start := loginStartReq(map[string]string{"email": created.Email, "password": "password123"})
		if start.Code != http.StatusOK {
			test.Fatalf("login start status = %d, want %d, body: %s", start.Code, http.StatusOK, start.Body.String())
		}

		recorder := loginVerifyReq(map[string]string{"email": created.Email, "code": "000000"})
		if recorder.Code != http.StatusUnauthorized {
			test.Errorf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
		}
	})

	test.Run("LoginVerify returns 400 for a missing field", func(test *testing.T) {
		test.Parallel()

		recorder := loginVerifyReq(map[string]string{"email": "missing-code@example.com", "code": ""})
		if recorder.Code != http.StatusBadRequest {
			test.Errorf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
		}
	})

	test.Run("LoginVerify returns 400 for an invalid request body", func(test *testing.T) {
		test.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login/verify", bytes.NewReader([]byte("not json")))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			test.Errorf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
		}
	})
}
