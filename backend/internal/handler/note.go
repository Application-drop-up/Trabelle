package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	pindomain "github.com/Application-drop-up/Travellle/internal/domain/pin"
	domain "github.com/Application-drop-up/Travellle/internal/domain/note"
	noteuc "github.com/Application-drop-up/Travellle/internal/usecase/note"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type NoteHandler struct {
	uc *noteuc.UseCase
}

func NewNoteHandler(uc *noteuc.UseCase) *NoteHandler {
	return &NoteHandler{uc: uc}
}

type createNoteRequest struct {
	Content string `json:"content"`
}

type updateNoteRequest struct {
	Content string `json:"content"`
}

type noteResponse struct {
	ID        string `json:"id"`
	PinID     string `json:"pin_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func toNoteResponse(n *domain.Note) noteResponse {
	return noteResponse{
		ID:        n.ID.String(),
		PinID:     n.PinID.String(),
		Content:   n.Content,
		CreatedAt: n.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: n.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	pinID, err := uuid.Parse(chi.URLParam(r, "pin_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pin_id")
		return
	}

	var req createNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	n, err := h.uc.CreateNote(r.Context(), noteuc.CreateInput{
		PinID:   pinID,
		Content: req.Content,
	})
	if errors.Is(err, pindomain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "pin not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, toNoteResponse(n))
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	pinID, err := uuid.Parse(chi.URLParam(r, "pin_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pin_id")
		return
	}
	noteID, err := uuid.Parse(chi.URLParam(r, "note_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid note_id")
		return
	}

	var req updateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	n, err := h.uc.UpdateNote(r.Context(), pinID, noteID, noteuc.UpdateInput{Content: req.Content})
	if errors.Is(err, domain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, toNoteResponse(n))
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	pinID, err := uuid.Parse(chi.URLParam(r, "pin_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pin_id")
		return
	}
	noteID, err := uuid.Parse(chi.URLParam(r, "note_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid note_id")
		return
	}

	if err := h.uc.DeleteNote(r.Context(), pinID, noteID); errors.Is(err, domain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "note not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
