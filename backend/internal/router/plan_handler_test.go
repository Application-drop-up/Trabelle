package router_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Application-drop-up/Travellle/internal/router"
	"github.com/Application-drop-up/Travellle/internal/testutil"
)

// createAuthenticatedUser registers a user and logs them in (reading the
// OTP code straight from the DB, same as a real login would read it from
// email), returning their session cookie and ID.
func createAuthenticatedUser(t *testing.T, r http.Handler, db *sql.DB, email string) (*http.Cookie, string) {
	t.Helper()

	registerBody, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": "password123",
		"name":     "Plan Creator",
	})
	registerReq := httptest.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewReader(registerBody))
	registerReq.Header.Set("Content-Type", "application/json")
	registerW := httptest.NewRecorder()
	r.ServeHTTP(registerW, registerReq)
	if registerW.Code != http.StatusCreated {
		t.Fatalf("registration status = %d, want %d, body: %s", registerW.Code, http.StatusCreated, registerW.Body.String())
	}
	var created struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(registerW.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode user response: %v", err)
	}

	startBody, _ := json.Marshal(map[string]string{"email": email, "password": "password123"})
	startReq := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(startBody))
	startReq.Header.Set("Content-Type", "application/json")
	startW := httptest.NewRecorder()
	r.ServeHTTP(startW, startReq)
	if startW.Code != http.StatusOK {
		t.Fatalf("login start status = %d, want %d, body: %s", startW.Code, http.StatusOK, startW.Body.String())
	}

	var code string
	if err := db.QueryRow(
		"SELECT code FROM login_otps WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1",
		created.ID,
	).Scan(&code); err != nil {
		t.Fatalf("failed to read otp code: %v", err)
	}

	verifyBody, _ := json.Marshal(map[string]string{"email": email, "code": code})
	verifyReq := httptest.NewRequest(http.MethodPost, "/api/v1/login/verify", bytes.NewReader(verifyBody))
	verifyReq.Header.Set("Content-Type", "application/json")
	verifyW := httptest.NewRecorder()
	r.ServeHTTP(verifyW, verifyReq)
	if verifyW.Code != http.StatusOK {
		t.Fatalf("login verify status = %d, want %d, body: %s", verifyW.Code, http.StatusOK, verifyW.Body.String())
	}

	for _, cookie := range verifyW.Result().Cookies() {
		if cookie.Name == "session_token" {
			return cookie, created.ID
		}
	}
	t.Fatal("expected a session_token cookie after login verify")
	return nil, ""
}

type planResponse struct {
	ID         string        `json:"id"`
	ShareToken string        `json:"share_token"`
	Title      string        `json:"title"`
	IsPublic   bool          `json:"is_public"`
	Pins       []interface{} `json:"pins"`
}

func TestPlanHandler_CreateAndGet(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"}, false)

	t.Run("creates a plan and returns it", func(t *testing.T) {
		t.Parallel()

		body, _ := json.Marshal(map[string]string{"title": "Trip to Kyoto"})
		req := httptest.NewRequest(http.MethodPost, "/plans", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("POST /plans status = %d, want %d, body: %s", w.Code, http.StatusCreated, w.Body.String())
		}

		var created planResponse
		if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM plans WHERE id = $1", created.ID) })

		if created.Title != "Trip to Kyoto" {
			t.Errorf("Title = %q, want %q", created.Title, "Trip to Kyoto")
		}
		if created.ShareToken == "" {
			t.Error("ShareToken is empty")
		}
		if len(created.Pins) != 0 {
			t.Errorf("Pins = %v, want empty", created.Pins)
		}

		t.Run("fetches the plan by share token", func(t *testing.T) {
			getReq := httptest.NewRequest(http.MethodGet, "/plans/"+created.ShareToken, nil)
			getW := httptest.NewRecorder()

			r.ServeHTTP(getW, getReq)

			if getW.Code != http.StatusOK {
				t.Fatalf("GET /plans/{share_token} status = %d, want %d, body: %s", getW.Code, http.StatusOK, getW.Body.String())
			}

			var fetched planResponse
			if err := json.Unmarshal(getW.Body.Bytes(), &fetched); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if fetched.ID != created.ID {
				t.Errorf("ID = %q, want %q", fetched.ID, created.ID)
			}
			if fetched.Title != created.Title {
				t.Errorf("Title = %q, want %q", fetched.Title, created.Title)
			}
		})
	})

	t.Run("rejects an empty title", func(t *testing.T) {
		t.Parallel()

		body, _ := json.Marshal(map[string]string{"title": ""})
		req := httptest.NewRequest(http.MethodPost, "/plans", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("POST /plans (empty title) status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns 404 for an unknown share token", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/plans/does-not-exist", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("GET /plans/{share_token} (unknown) status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}

func TestPlanHandler_Publish(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"}, false)

	t.Run("makes the plan public", func(t *testing.T) {
		t.Parallel()

		body, _ := json.Marshal(map[string]string{"title": "Trip to Osaka"})
		createReq := httptest.NewRequest(http.MethodPost, "/plans", bytes.NewReader(body))
		createReq.Header.Set("Content-Type", "application/json")
		createW := httptest.NewRecorder()
		r.ServeHTTP(createW, createReq)

		var created planResponse
		if err := json.Unmarshal(createW.Body.Bytes(), &created); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM plans WHERE id = $1", created.ID) })

		if created.IsPublic {
			t.Fatal("newly created plan IsPublic = true, want false")
		}

		publishReq := httptest.NewRequest(http.MethodPost, "/api/v1/plans/"+created.ShareToken+"/publish", nil)
		publishW := httptest.NewRecorder()
		r.ServeHTTP(publishW, publishReq)

		if publishW.Code != http.StatusOK {
			t.Fatalf("POST /api/v1/plans/{share_token}/publish status = %d, want %d, body: %s", publishW.Code, http.StatusOK, publishW.Body.String())
		}

		var published planResponse
		if err := json.Unmarshal(publishW.Body.Bytes(), &published); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if !published.IsPublic {
			t.Error("after Publish(), IsPublic = false, want true")
		}

		getReq := httptest.NewRequest(http.MethodGet, "/plans/"+created.ShareToken, nil)
		getW := httptest.NewRecorder()
		r.ServeHTTP(getW, getReq)

		var fetched planResponse
		if err := json.Unmarshal(getW.Body.Bytes(), &fetched); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if !fetched.IsPublic {
			t.Error("after Publish(), GET /plans/{share_token} IsPublic = false, want true")
		}
	})

	t.Run("returns 404 for an unknown share token", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/plans/does-not-exist/publish", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("POST /api/v1/plans/{share_token}/publish (unknown) status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}

func TestPlanHandler_ListForUser(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"}, false)

	t.Run("returns only the plans the user is a member of", func(t *testing.T) {
		t.Parallel()

		userID := createTestUserForMember(t, r)
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM users WHERE id = $1", userID) })

		memberPlanID := createTestPlan(t, r, "Member Plan")
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM plans WHERE id = $1", memberPlanID) })

		nonMemberPlanID := createTestPlan(t, r, "Non-Member Plan")
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM plans WHERE id = $1", nonMemberPlanID) })

		addResp := addTestMember(t, r, memberPlanID, userID)
		if addResp.Code != http.StatusCreated {
			t.Fatalf("add member status = %d, want %d, body: %s", addResp.Code, http.StatusCreated, addResp.Body.String())
		}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/user/"+userID+"/plans", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d, body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		var plans []planResponse
		if err := json.Unmarshal(w.Body.Bytes(), &plans); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(plans) != 1 {
			t.Fatalf("got %d plans, want 1", len(plans))
		}
		if plans[0].ID != memberPlanID {
			t.Errorf("ID = %q, want %q", plans[0].ID, memberPlanID)
		}
	})

	t.Run("returns an empty list for a user with no memberships", func(t *testing.T) {
		t.Parallel()

		userID := createTestUserForMember(t, r)
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM users WHERE id = $1", userID) })

		req := httptest.NewRequest(http.MethodGet, "/api/v1/user/"+userID+"/plans", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d, body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		var plans []planResponse
		if err := json.Unmarshal(w.Body.Bytes(), &plans); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(plans) != 0 {
			t.Errorf("got %d plans, want 0", len(plans))
		}
	})

	t.Run("rejects an invalid user id", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/user/not-a-uuid/plans", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

func TestPlanHandler_Create_AutoAddsCreatorAsMember(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"}, true)

	t.Run("adds the creator as a member when authenticated", func(t *testing.T) {
		t.Parallel()

		cookie, userID := createAuthenticatedUser(t, r, db, "plan-creator@example.com")
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM users WHERE id = $1", userID) })

		body, _ := json.Marshal(map[string]string{"title": "Authenticated Trip"})
		req := httptest.NewRequest(http.MethodPost, "/plans", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(cookie)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d, body: %s", w.Code, http.StatusCreated, w.Body.String())
		}
		var created planResponse
		if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM plans WHERE id = $1", created.ID) })

		listReq := httptest.NewRequest(http.MethodGet, "/api/v1/user/"+userID+"/plans", nil)
		listW := httptest.NewRecorder()
		r.ServeHTTP(listW, listReq)

		var plans []planResponse
		if err := json.Unmarshal(listW.Body.Bytes(), &plans); err != nil {
			t.Fatalf("failed to decode list response: %v", err)
		}
		if len(plans) != 1 || plans[0].ID != created.ID {
			t.Errorf("GET /user/{id}/plans = %+v, want a single plan with ID %q", plans, created.ID)
		}
	})

	t.Run("still creates the plan anonymously, with no membership", func(t *testing.T) {
		t.Parallel()

		body, _ := json.Marshal(map[string]string{"title": "Anonymous Trip"})
		req := httptest.NewRequest(http.MethodPost, "/plans", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d, body: %s", w.Code, http.StatusCreated, w.Body.String())
		}
		var created planResponse
		if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM plans WHERE id = $1", created.ID) })

		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM plan_members WHERE plan_id = $1", created.ID).Scan(&count); err != nil {
			t.Fatalf("failed to count plan members: %v", err)
		}
		if count != 0 {
			t.Errorf("plan_members count = %d, want 0", count)
		}
	})

	t.Run("still creates the plan when the session cookie is invalid", func(t *testing.T) {
		t.Parallel()

		body, _ := json.Marshal(map[string]string{"title": "Invalid Session Trip"})
		req := httptest.NewRequest(http.MethodPost, "/plans", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "session_token", Value: "unknown-token"})
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d, body: %s", w.Code, http.StatusCreated, w.Body.String())
		}
		var created planResponse
		if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM plans WHERE id = $1", created.ID) })
	})
}
