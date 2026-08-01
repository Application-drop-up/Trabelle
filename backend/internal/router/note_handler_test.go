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

// createTestPlanForNote and createTestPinForNote use locally-scoped
// response types to avoid colliding with planResponse/pinResponse types
// declared by Plan/Pin handler tests.
func createTestPlanForNote(t *testing.T, r http.Handler, title string) string {
	t.Helper()

	body, _ := json.Marshal(map[string]string{"title": title})
	req := httptest.NewRequest(http.MethodPost, "/plans", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("POST /plans status = %d, want %d, body: %s", w.Code, http.StatusCreated, w.Body.String())
	}

	var plan struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &plan); err != nil {
		t.Fatalf("failed to decode plan response: %v", err)
	}
	return plan.ID
}

func createTestPinForNote(t *testing.T, r http.Handler, planID string) string {
	t.Helper()

	body, _ := json.Marshal(map[string]any{
		"name":      "Tokyo Tower",
		"latitude":  35.6586,
		"longitude": 139.7454,
		"category":  "sightseeing",
		"colour":    "#FF5733",
	})
	req := httptest.NewRequest(http.MethodPost, "/plans/"+planID+"/pins", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("POST /plans/{plan_id}/pins status = %d, want %d, body: %s", w.Code, http.StatusCreated, w.Body.String())
	}

	var pin struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &pin); err != nil {
		t.Fatalf("failed to decode pin response: %v", err)
	}
	return pin.ID
}

type noteResponse struct {
	ID      string `json:"id"`
	PinID   string `json:"pin_id"`
	Content string `json:"content"`
}

func TestNoteHandler_CreateUpdateDelete(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"})

	planID := createTestPlanForNote(t, r, "Note Handler Test Plan")
	t.Cleanup(func() { _, _ = db.Exec("DELETE FROM plans WHERE id = $1", planID) })
	pinID := createTestPinForNote(t, r, planID)

	notesPath := "/plans/" + planID + "/pins/" + pinID + "/notes"

	t.Run("rejects empty content", func(t *testing.T) {
		t.Parallel()

		body, _ := json.Marshal(map[string]string{"content": ""})
		req := httptest.NewRequest(http.MethodPost, notesPath, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns 404 for a nonexistent pin", func(t *testing.T) {
		t.Parallel()

		body, _ := json.Marshal(map[string]string{"content": "hello"})
		req := httptest.NewRequest(http.MethodPost, "/plans/"+planID+"/pins/00000000-0000-0000-0000-000000000000/notes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d, body: %s", w.Code, http.StatusNotFound, w.Body.String())
		}
	})

	t.Run("creates, updates, and deletes a note", func(t *testing.T) {
		createBody, _ := json.Marshal(map[string]string{"content": "Great view at sunset"})
		createReq := httptest.NewRequest(http.MethodPost, notesPath, bytes.NewReader(createBody))
		createReq.Header.Set("Content-Type", "application/json")
		createW := httptest.NewRecorder()
		r.ServeHTTP(createW, createReq)
		if createW.Code != http.StatusCreated {
			t.Fatalf("POST note status = %d, want %d, body: %s", createW.Code, http.StatusCreated, createW.Body.String())
		}
		var note noteResponse
		if err := json.Unmarshal(createW.Body.Bytes(), &note); err != nil {
			t.Fatalf("failed to decode create response: %v", err)
		}

		updateBody, _ := json.Marshal(map[string]string{"content": "Updated content"})
		updateReq := httptest.NewRequest(http.MethodPatch, notesPath+"/"+note.ID, bytes.NewReader(updateBody))
		updateReq.Header.Set("Content-Type", "application/json")
		updateW := httptest.NewRecorder()
		r.ServeHTTP(updateW, updateReq)
		if updateW.Code != http.StatusOK {
			t.Fatalf("PATCH note status = %d, want %d, body: %s", updateW.Code, http.StatusOK, updateW.Body.String())
		}
		var updated noteResponse
		if err := json.Unmarshal(updateW.Body.Bytes(), &updated); err != nil {
			t.Fatalf("failed to decode update response: %v", err)
		}
		if updated.Content != "Updated content" {
			t.Errorf("after update, Content = %q, want %q", updated.Content, "Updated content")
		}

		deleteReq := httptest.NewRequest(http.MethodDelete, notesPath+"/"+note.ID, nil)
		deleteW := httptest.NewRecorder()
		r.ServeHTTP(deleteW, deleteReq)
		if deleteW.Code != http.StatusNoContent {
			t.Fatalf("DELETE note status = %d, want %d", deleteW.Code, http.StatusNoContent)
		}

		deleteAgainReq := httptest.NewRequest(http.MethodDelete, notesPath+"/"+note.ID, nil)
		deleteAgainW := httptest.NewRecorder()
		r.ServeHTTP(deleteAgainW, deleteAgainReq)
		if deleteAgainW.Code != http.StatusNotFound {
			t.Errorf("DELETE (already deleted) status = %d, want %d", deleteAgainW.Code, http.StatusNotFound)
		}
	})
}
