package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	domain "github.com/Application-drop-up/Travellle/internal/domain/plan"
	planuc "github.com/Application-drop-up/Travellle/internal/usecase/plan"
	"github.com/go-chi/chi/v5"
)

type PlanHandler struct {
	uc *planuc.UseCase
}

func NewPlanHandler(uc *planuc.UseCase) *PlanHandler {
	return &PlanHandler{uc: uc}
}

type createPlanRequest struct {
	Title string `json:"title"`
}

type planResponse struct {
	ID         string `json:"id"`
	ShareToken string `json:"share_token"`
	Title      string `json:"title"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func toPlanResponse(p *domain.Plan) planResponse {
	return planResponse{
		ID:         p.ID.String(),
		ShareToken: p.ShareToken,
		Title:      p.Title,
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

	writeJSON(w, http.StatusCreated, toPlanResponse(p))
}

func (h *PlanHandler) GetByShareToken(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "share_token")

	p, err := h.uc.GetPlanByShareToken(r.Context(), token)
	if errors.Is(err, domain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "plan not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, toPlanResponse(p))
}
