package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Application-drop-up/Travellle/internal/router"
	"github.com/Application-drop-up/Travellle/internal/testutil"
	"github.com/google/uuid"
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
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"}, false)

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

func TestUserHandler_GetByID(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"}, false)

	getReq := func(id string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/user/"+id, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	t.Run("returns the user when it exists", func(t *testing.T) {
		t.Parallel()

		registerRaw, _ := json.Marshal(map[string]string{
			"email":    "getbyid@example.com",
			"password": "password123",
			"name":     "GetByID",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewReader(registerRaw))
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

		got := getReq(created.ID)
		if got.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d, body: %s", got.Code, http.StatusOK, got.Body.String())
		}

		var fetched userResponse
		if err := json.Unmarshal(got.Body.Bytes(), &fetched); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if fetched.ID != created.ID {
			t.Errorf("ID = %q, want %q", fetched.ID, created.ID)
		}
		if fetched.Email != created.Email {
			t.Errorf("Email = %q, want %q", fetched.Email, created.Email)
		}
	})

	t.Run("returns 404 for an unknown id", func(t *testing.T) {
		t.Parallel()

		w := getReq(uuid.New().String())

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("returns 400 for an invalid id", func(t *testing.T) {
		t.Parallel()

		w := getReq("not-a-uuid")

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

func TestUserHandler_Update(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"}, false)

	updateReq := func(id string, body any) *httptest.ResponseRecorder {
		raw, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/user/"+id, bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	registerUser := func(t *testing.T, email string) userResponse {
		t.Helper()
		raw, _ := json.Marshal(map[string]string{
			"email":    email,
			"password": "password123",
			"name":     "Original",
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

	t.Run("updates the name and email", func(t *testing.T) {
		t.Parallel()

		created := registerUser(t, "update@example.com")

		w := updateReq(created.ID, map[string]string{"name": "Updated", "email": "updated@example.com"})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d, body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		var updated userResponse
		if err := json.Unmarshal(w.Body.Bytes(), &updated); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if updated.Name != "Updated" {
			t.Errorf("Name = %q, want %q", updated.Name, "Updated")
		}
		if updated.Email != "updated@example.com" {
			t.Errorf("Email = %q, want %q", updated.Email, "updated@example.com")
		}
	})

	t.Run("rejects an empty name", func(t *testing.T) {
		t.Parallel()

		created := registerUser(t, "update-empty-name@example.com")

		w := updateReq(created.ID, map[string]string{"name": "", "email": created.Email})
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("rejects an empty email", func(t *testing.T) {
		t.Parallel()

		created := registerUser(t, "update-empty-email@example.com")

		w := updateReq(created.ID, map[string]string{"name": "Updated", "email": ""})
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("rejects an invalid email", func(t *testing.T) {
		t.Parallel()

		created := registerUser(t, "update-invalid-email@example.com")

		w := updateReq(created.ID, map[string]string{"name": "Updated", "email": "not-an-email"})
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns 409 when the email is already taken", func(t *testing.T) {
		t.Parallel()

		registerUser(t, "update-taken@example.com")
		created := registerUser(t, "update-taker@example.com")

		w := updateReq(created.ID, map[string]string{"name": "Updated", "email": "update-taken@example.com"})
		if w.Code != http.StatusConflict {
			t.Errorf("status = %d, want %d, body: %s", w.Code, http.StatusConflict, w.Body.String())
		}
	})

	t.Run("returns 404 for an unknown id", func(t *testing.T) {
		t.Parallel()

		w := updateReq(uuid.New().String(), map[string]string{"name": "Updated", "email": "unknown@example.com"})

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("returns 400 for an invalid id", func(t *testing.T) {
		t.Parallel()

		w := updateReq("not-a-uuid", map[string]string{"name": "Updated", "email": "updated@example.com"})

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

func TestUserHandler_Delete(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"}, false)

	deleteReq := func(id string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/user/"+id, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	t.Run("deletes an existing user", func(t *testing.T) {
		t.Parallel()

		registerRaw, _ := json.Marshal(map[string]string{
			"email":    "delete@example.com",
			"password": "password123",
			"name":     "Delete",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewReader(registerRaw))
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

		got := deleteReq(created.ID)
		if got.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d, body: %s", got.Code, http.StatusNoContent, got.Body.String())
		}

		getReq := httptest.NewRequest(http.MethodGet, "/api/v1/user/"+created.ID, nil)
		getW := httptest.NewRecorder()
		r.ServeHTTP(getW, getReq)
		if getW.Code != http.StatusNotFound {
			t.Errorf("GET after delete status = %d, want %d", getW.Code, http.StatusNotFound)
		}
	})

	t.Run("returns 404 for an unknown id", func(t *testing.T) {
		t.Parallel()

		w := deleteReq(uuid.New().String())

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("returns 400 for an invalid id", func(t *testing.T) {
		t.Parallel()

		w := deleteReq("not-a-uuid")

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}
