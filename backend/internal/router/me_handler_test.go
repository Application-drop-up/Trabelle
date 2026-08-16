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

func TestMeHandler(t *testing.T) {
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

	loginSessionCookie := func(t *testing.T, email, password string) (*http.Cookie, userResponse) {
		t.Helper()
		created := registerUser(t, email, password)
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM login_otps WHERE user_id = $1", created.ID) })
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM sessions WHERE user_id = $1", created.ID) })

		startRaw, _ := json.Marshal(map[string]string{"email": email, "password": password})
		startReq := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(startRaw))
		startReq.Header.Set("Content-Type", "application/json")
		startW := httptest.NewRecorder()
		r.ServeHTTP(startW, startReq)
		if startW.Code != http.StatusOK {
			t.Fatalf("login start status = %d, want %d, body: %s", startW.Code, http.StatusOK, startW.Body.String())
		}

		verifyRaw, _ := json.Marshal(map[string]string{"email": email, "code": otpCodeFor(t, created.ID)})
		verifyReq := httptest.NewRequest(http.MethodPost, "/api/v1/login/verify", bytes.NewReader(verifyRaw))
		verifyReq.Header.Set("Content-Type", "application/json")
		verifyW := httptest.NewRecorder()
		r.ServeHTTP(verifyW, verifyReq)
		if verifyW.Code != http.StatusOK {
			t.Fatalf("login verify status = %d, want %d, body: %s", verifyW.Code, http.StatusOK, verifyW.Body.String())
		}

		cookies := verifyW.Result().Cookies()
		for _, cookie := range cookies {
			if cookie.Name == "session_token" {
				return cookie, created
			}
		}
		t.Fatal("expected a session_token cookie after login verify")
		return nil, created
	}

	meReq := func(cookie *http.Cookie) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/user/me", nil)
		if cookie != nil {
			req.AddCookie(cookie)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	t.Run("returns the current user for a valid session", func(t *testing.T) {
		t.Parallel()

		cookie, created := loginSessionCookie(t, "me@example.com", "password123")

		w := meReq(cookie)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d, body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		var got userResponse
		if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.ID != created.ID {
			t.Errorf("ID = %q, want %q", got.ID, created.ID)
		}
		if got.Email != created.Email {
			t.Errorf("Email = %q, want %q", got.Email, created.Email)
		}
	})

	t.Run("returns 401 when no session cookie is present", func(t *testing.T) {
		t.Parallel()

		w := meReq(nil)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})

	t.Run("returns 401 for an unknown session token", func(t *testing.T) {
		t.Parallel()

		w := meReq(&http.Cookie{Name: "session_token", Value: "unknown-token"})
		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})
}
