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

func toNoteResponse(note *domain.Note) noteResponse {
	return noteResponse{
		ID:        note.ID.String(),
		PinID:     note.PinID.String(),
		Content:   note.Content,
		CreatedAt: note.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: note.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func (nh *NoteHandler) Create(rw http.ResponseWriter, req *http.Request) {
	pinID, err := uuid.Parse(chi.URLParam(req, "pin_id"))
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid pin_id")
		return
	}

	var body createNoteRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil || body.Content == "" {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	note, err := nh.uc.CreateNote(req.Context(), noteuc.CreateInput{
		PinID:   pinID,
		Content: body.Content,
	})
	if errors.Is(err, pindomain.ErrNotFound) {
		writeError(rw, http.StatusNotFound, "pin not found")
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(rw, http.StatusCreated, toNoteResponse(note))
}

func (nh *NoteHandler) Update(rw http.ResponseWriter, req *http.Request) {
	pinID, err := uuid.Parse(chi.URLParam(req, "pin_id"))
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid pin_id")
		return
	}
	noteID, err := uuid.Parse(chi.URLParam(req, "note_id"))
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid note_id")
		return
	}

	var body updateNoteRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil || body.Content == "" {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	note, err := nh.uc.UpdateNote(req.Context(), pinID, noteID, noteuc.UpdateInput{Content: body.Content})
	if errors.Is(err, domain.ErrNotFound) {
		writeError(rw, http.StatusNotFound, "note not found")
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(rw, http.StatusOK, toNoteResponse(note))
}

func (nh *NoteHandler) Delete(rw http.ResponseWriter, req *http.Request) {
	pinID, err := uuid.Parse(chi.URLParam(req, "pin_id"))
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid pin_id")
		return
	}
	noteID, err := uuid.Parse(chi.URLParam(req, "note_id"))
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid note_id")
		return
	}

	if err := nh.uc.DeleteNote(req.Context(), pinID, noteID); errors.Is(err, domain.ErrNotFound) {
		writeError(rw, http.StatusNotFound, "note not found")
		return
	} else if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
