package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	domain "github.com/Application-drop-up/Travellle/internal/domain/plan"
	noteuc "github.com/Application-drop-up/Travellle/internal/usecase/note"
	pinuc "github.com/Application-drop-up/Travellle/internal/usecase/pin"
	planuc "github.com/Application-drop-up/Travellle/internal/usecase/plan"
	"github.com/go-chi/chi/v5"
)

type PlanHandler struct {
	uc    *planuc.UseCase
	pinUC *pinuc.UseCase
	noteUC *noteuc.UseCase
}

func NewPlanHandler(uc *planuc.UseCase, pinUC *pinuc.UseCase, noteUC *noteuc.UseCase) *PlanHandler {
	return &PlanHandler{uc: uc, pinUC: pinUC, noteUC: noteUC}
}

type createPlanRequest struct {
	Title string `json:"title"`
}

type planResponse struct {
	ID         string        `json:"id"`
	ShareToken string        `json:"share_token"`
	Title      string        `json:"title"`
	Pins       []pinWithNotes `json:"pins"`
	CreatedAt  string        `json:"created_at"`
	UpdatedAt  string        `json:"updated_at"`
}

type pinWithNotes struct {
	pinResponse
	Notes []noteResponse `json:"notes"`
}

func toPlanResponse(p *domain.Plan, pins []pinWithNotes) planResponse {
	return planResponse{
		ID:         p.ID.String(),
		ShareToken: p.ShareToken,
		Title:      p.Title,
		Pins:       pins,
		CreatedAt:  p.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:  p.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func (h *PlanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p, err := h.uc.CreatePlan(r.Context(), req.Title)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, toPlanResponse(p, []pinWithNotes{}))
}

func (h *PlanHandler) GetByShareToken(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "share_token")
	ctx := r.Context()

	p, err := h.uc.GetPlanByShareToken(ctx, token)
	if errors.Is(err, domain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "plan not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	rawPins, err := h.pinUC.ListPins(ctx, p.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	pins := make([]pinWithNotes, 0, len(rawPins))
	for _, pin := range rawPins {
		rawNotes, err := h.noteUC.ListNotes(ctx, pin.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		notes := make([]noteResponse, 0, len(rawNotes))
		for _, n := range rawNotes {
			notes = append(notes, toNoteResponse(n))
		}
		pins = append(pins, pinWithNotes{
			pinResponse: toPinResponse(pin),
			Notes:       notes,
		})
	}

	writeJSON(w, http.StatusOK, toPlanResponse(p, pins))
}
