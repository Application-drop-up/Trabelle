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

type planMemberResponse struct {
	ID     string `json:"id"`
	PlanID string `json:"plan_id"`
	UserID string `json:"user_id"`
}

func createTestUserForMember(t *testing.T, r http.Handler) string {
	t.Helper()

	body, _ := json.Marshal(map[string]string{
		"email":    uuid.New().String() + "@example.com",
		"password": "password123",
		"name":     "Member Test User",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("user registration status = %d, want %d, body: %s", w.Code, http.StatusCreated, w.Body.String())
	}
	var created struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode user response: %v", err)
	}
	return created.ID
}

func addTestMember(t *testing.T, r http.Handler, planID, userID string) *httptest.ResponseRecorder {
	t.Helper()

	body, _ := json.Marshal(map[string]string{"user_id": userID})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/plans/"+planID+"/members", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestPlanMemberHandler_AddListRemove(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"}, false)

	planID := createTestPlan(t, r, "Plan Member Handler Test Plan")
	t.Cleanup(func() { _, _ = db.Exec("DELETE FROM plans WHERE id = $1", planID) })

	t.Run("rejects an invalid plan_id", func(t *testing.T) {
		t.Parallel()

		userID := createTestUserForMember(t, r)
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM users WHERE id = $1", userID) })

		w := addTestMember(t, r, "not-a-uuid", userID)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("rejects an invalid request body", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/plans/"+planID+"/members", bytes.NewReader([]byte("not json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns 404 for a nonexistent plan", func(t *testing.T) {
		t.Parallel()

		userID := createTestUserForMember(t, r)
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM users WHERE id = $1", userID) })

		w := addTestMember(t, r, uuid.New().String(), userID)
		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d, body: %s", w.Code, http.StatusNotFound, w.Body.String())
		}
	})

	t.Run("returns 404 for a nonexistent user", func(t *testing.T) {
		t.Parallel()

		w := addTestMember(t, r, planID, uuid.New().String())
		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d, body: %s", w.Code, http.StatusNotFound, w.Body.String())
		}
	})

	t.Run("adds, lists, and removes a member; rejects a duplicate add", func(t *testing.T) {
		t.Parallel()

		userID := createTestUserForMember(t, r)
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM users WHERE id = $1", userID) })

		addW := addTestMember(t, r, planID, userID)
		if addW.Code != http.StatusCreated {
			t.Fatalf("add status = %d, want %d, body: %s", addW.Code, http.StatusCreated, addW.Body.String())
		}
		var added planMemberResponse
		if err := json.Unmarshal(addW.Body.Bytes(), &added); err != nil {
			t.Fatalf("failed to decode add response: %v", err)
		}
		if added.UserID != userID {
			t.Errorf("UserID = %q, want %q", added.UserID, userID)
		}

		dupW := addTestMember(t, r, planID, userID)
		if dupW.Code != http.StatusConflict {
			t.Errorf("duplicate add status = %d, want %d, body: %s", dupW.Code, http.StatusConflict, dupW.Body.String())
		}

		listReq := httptest.NewRequest(http.MethodGet, "/api/v1/plans/"+planID+"/members", nil)
		listW := httptest.NewRecorder()
		r.ServeHTTP(listW, listReq)
		if listW.Code != http.StatusOK {
			t.Fatalf("list status = %d, want %d, body: %s", listW.Code, http.StatusOK, listW.Body.String())
		}
		var members []planMemberResponse
		if err := json.Unmarshal(listW.Body.Bytes(), &members); err != nil {
			t.Fatalf("failed to decode list response: %v", err)
		}
		if len(members) != 1 || members[0].UserID != userID {
			t.Errorf("List() = %+v, want a single member with UserID %q", members, userID)
		}

		deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/plans/"+planID+"/members/"+userID, nil)
		deleteW := httptest.NewRecorder()
		r.ServeHTTP(deleteW, deleteReq)
		if deleteW.Code != http.StatusNoContent {
			t.Fatalf("delete status = %d, want %d, body: %s", deleteW.Code, http.StatusNoContent, deleteW.Body.String())
		}

		deleteAgainReq := httptest.NewRequest(http.MethodDelete, "/api/v1/plans/"+planID+"/members/"+userID, nil)
		deleteAgainW := httptest.NewRecorder()
		r.ServeHTTP(deleteAgainW, deleteAgainReq)
		if deleteAgainW.Code != http.StatusNotFound {
			t.Errorf("delete (already removed) status = %d, want %d", deleteAgainW.Code, http.StatusNotFound)
		}
	})
}
