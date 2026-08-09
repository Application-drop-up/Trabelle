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

type userResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

func TestUserHandler_Register(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"})

	registerReq := func(body any) *httptest.ResponseRecorder {
		raw, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	t.Run("registers a user and returns it without the password", func(t *testing.T) {
		t.Parallel()

		w := registerReq(map[string]string{
			"email":    "taro@example.com",
			"password": "password123",
			"name":     "Taro",
		})

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d, body: %s", w.Code, http.StatusCreated, w.Body.String())
		}

		var created userResponse
		if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM users WHERE id = $1", created.ID) })

		if created.Email != "taro@example.com" {
			t.Errorf("Email = %q, want %q", created.Email, "taro@example.com")
		}
		if created.Name != "Taro" {
			t.Errorf("Name = %q, want %q", created.Name, "Taro")
		}
		if created.ID == "" {
			t.Error("ID is empty")
		}
		if created.CreatedAt == "" {
			t.Error("CreatedAt is empty")
		}
		if bytes.Contains(w.Body.Bytes(), []byte("password")) {
			t.Errorf("response body leaks password field: %s", w.Body.String())
		}
	})

	t.Run("rejects a missing field", func(t *testing.T) {
		t.Parallel()

		w := registerReq(map[string]string{
			"email":    "missing-password@example.com",
			"password": "",
			"name":     "No Password",
		})

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("rejects an invalid request body", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewReader([]byte("not json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("rejects a duplicate email", func(t *testing.T) {
		t.Parallel()

		body := map[string]string{
			"email":    "dup@example.com",
			"password": "password123",
			"name":     "First",
		}

		first := registerReq(body)
		if first.Code != http.StatusCreated {
			t.Fatalf("first registration status = %d, want %d, body: %s", first.Code, http.StatusCreated, first.Body.String())
		}
		var created userResponse
		if err := json.Unmarshal(first.Body.Bytes(), &created); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM users WHERE id = $1", created.ID) })

		second := registerReq(body)
		if second.Code != http.StatusConflict {
			t.Errorf("second registration status = %d, want %d, body: %s", second.Code, http.StatusConflict, second.Body.String())
		}
	})
}
