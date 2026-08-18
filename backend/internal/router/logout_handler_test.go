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

func TestLogoutHandler(test *testing.T) {
	test.Parallel()

	db := testutil.NewTestDB(test)
	mux := router.New(db, "test-api-key", []string{"http://localhost:3000"}, true)

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

	loginSessionCookie := func(test *testing.T, email, password string) *http.Cookie {
		test.Helper()
		created := registerUser(test, email, password)
		test.Cleanup(func() { _, _ = db.Exec("DELETE FROM login_otps WHERE user_id = $1", created.ID) })
		test.Cleanup(func() { _, _ = db.Exec("DELETE FROM sessions WHERE user_id = $1", created.ID) })

		startRaw, _ := json.Marshal(map[string]string{"email": email, "password": password})
		startReq := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(startRaw))
		startReq.Header.Set("Content-Type", "application/json")
		startRecorder := httptest.NewRecorder()
		mux.ServeHTTP(startRecorder, startReq)
		if startRecorder.Code != http.StatusOK {
			test.Fatalf("login start status = %d, want %d, body: %s", startRecorder.Code, http.StatusOK, startRecorder.Body.String())
		}

		verifyRaw, _ := json.Marshal(map[string]string{"email": email, "code": otpCodeFor(test, created.ID)})
		verifyReq := httptest.NewRequest(http.MethodPost, "/api/v1/login/verify", bytes.NewReader(verifyRaw))
		verifyReq.Header.Set("Content-Type", "application/json")
		verifyRecorder := httptest.NewRecorder()
		mux.ServeHTTP(verifyRecorder, verifyReq)
		if verifyRecorder.Code != http.StatusOK {
			test.Fatalf("login verify status = %d, want %d, body: %s", verifyRecorder.Code, http.StatusOK, verifyRecorder.Body.String())
		}

		for _, cookie := range verifyRecorder.Result().Cookies() {
			if cookie.Name == "session_token" {
				return cookie
			}
		}
		test.Fatal("expected a session_token cookie after login verify")
		return nil
	}

	logoutReq := func(cookie *http.Cookie) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/logout", nil)
		if cookie != nil {
			req.AddCookie(cookie)
		}
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)
		return recorder
	}

	meReq := func(cookie *http.Cookie) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/user/me", nil)
		if cookie != nil {
			req.AddCookie(cookie)
		}
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)
		return recorder
	}

	test.Run("logs out, clears the cookie, and invalidates the session", func(test *testing.T) {
		test.Parallel()

		cookie := loginSessionCookie(test, "logout@example.com", "password123")

		recorder := logoutReq(cookie)
		if recorder.Code != http.StatusNoContent {
			test.Fatalf("status = %d, want %d, body: %s", recorder.Code, http.StatusNoContent, recorder.Body.String())
		}

		var cleared *http.Cookie
		for _, c := range recorder.Result().Cookies() {
			if c.Name == "session_token" {
				cleared = c
			}
		}
		if cleared == nil {
			test.Fatal("expected a session_token cookie to be set (cleared)")
		}
		if cleared.MaxAge >= 0 {
			test.Errorf("MaxAge = %d, want negative (cookie deletion)", cleared.MaxAge)
		}

		meRecorder := meReq(cookie)
		if meRecorder.Code != http.StatusUnauthorized {
			test.Errorf("GET /user/me after logout status = %d, want %d", meRecorder.Code, http.StatusUnauthorized)
		}
	})

	test.Run("returns 204 when no session cookie is present", func(test *testing.T) {
		test.Parallel()

		recorder := logoutReq(nil)
		if recorder.Code != http.StatusNoContent {
			test.Errorf("status = %d, want %d", recorder.Code, http.StatusNoContent)
		}
	})

	test.Run("returns 204 for an unknown session token", func(test *testing.T) {
		test.Parallel()

		recorder := logoutReq(&http.Cookie{Name: "session_token", Value: "unknown-token"})
		if recorder.Code != http.StatusNoContent {
			test.Errorf("status = %d, want %d", recorder.Code, http.StatusNoContent)
		}
	})
}
