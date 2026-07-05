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
	uc     *planuc.UseCase
	pinUC  *pinuc.UseCase
	noteUC *noteuc.UseCase
}

func NewPlanHandler(uc *planuc.UseCase, pinUC *pinuc.UseCase, noteUC *noteuc.UseCase) *PlanHandler {
	return &PlanHandler{uc: uc, pinUC: pinUC, noteUC: noteUC}
}

type createPlanRequest struct {
	Title string `json:"title"`
}

type planResponse struct {
	ID         string         `json:"id"`
	ShareToken string         `json:"share_token"`
	Title      string         `json:"title"`
	Pins       []pinWithNotes `json:"pins"`
	CreatedAt  string         `json:"created_at"`
	UpdatedAt  string         `json:"updated_at"`
}

type pinWithNotes struct {
	pinResponse
	Notes []noteResponse `json:"notes"`
}

func toPlanResponse(plan *domain.Plan, pins []pinWithNotes) planResponse {
	return planResponse{
		ID:         plan.ID.String(),
		ShareToken: plan.ShareToken,
		Title:      plan.Title,
		Pins:       pins,
		CreatedAt:  plan.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:  plan.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func (ph *PlanHandler) Create(rw http.ResponseWriter, req *http.Request) {
	var body createPlanRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil || body.Title == "" {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	plan, err := ph.uc.CreatePlan(req.Context(), body.Title)
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(rw, http.StatusCreated, toPlanResponse(plan, []pinWithNotes{}))
}

func (ph *PlanHandler) GetByShareToken(rw http.ResponseWriter, req *http.Request) {
	token := chi.URLParam(req, "share_token")
	ctx := req.Context()

	plan, err := ph.uc.GetPlanByShareToken(ctx, token)
	if errors.Is(err, domain.ErrNotFound) {
		writeError(rw, http.StatusNotFound, "plan not found")
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	rawPins, err := ph.pinUC.ListPins(ctx, plan.ID)
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	pins := make([]pinWithNotes, 0, len(rawPins))
	for _, pin := range rawPins {
		rawNotes, err := ph.noteUC.ListNotes(ctx, pin.ID)
		if err != nil {
			writeError(rw, http.StatusInternalServerError, "internal server error")
			return
		}
		notes := make([]noteResponse, 0, len(rawNotes))
		for _, note := range rawNotes {
			notes = append(notes, toNoteResponse(note))
		}
		pins = append(pins, pinWithNotes{
			pinResponse: toPinResponse(pin),
			Notes:       notes,
		})
	}

	writeJSON(rw, http.StatusOK, toPlanResponse(plan, pins))
}
