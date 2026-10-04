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
	planmemberuc "github.com/Application-drop-up/Travellle/internal/usecase/planmember"
	useruc "github.com/Application-drop-up/Travellle/internal/usecase/user"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PlanHandler struct {
	planUseCase       *planuc.UseCase
	pinUseCase        *pinuc.UseCase
	noteUseCase       *noteuc.UseCase
	userUseCase       *useruc.UseCase
	planMemberUseCase *planmemberuc.UseCase
}

func NewPlanHandler(planUseCase *planuc.UseCase, pinUseCase *pinuc.UseCase, noteUseCase *noteuc.UseCase, userUseCase *useruc.UseCase, planMemberUseCase *planmemberuc.UseCase) *PlanHandler {
	return &PlanHandler{
		planUseCase:       planUseCase,
		pinUseCase:        pinUseCase,
		noteUseCase:       noteUseCase,
		userUseCase:       userUseCase,
		planMemberUseCase: planMemberUseCase,
	}
}

// currentUserID resolves the authenticated user from the session cookie, if
// any. ok is false for any failure (no cookie, invalid/expired session) --
// callers must treat that as "anonymous", not an error, since most Plan
// endpoints work without authentication by design.
func (planHandler *PlanHandler) currentUserID(req *http.Request) (uuid.UUID, bool) {
	cookie, err := req.Cookie(sessionCookieName)
	if err != nil {
		return uuid.UUID{}, false
	}
	dto, err := planHandler.userUseCase.CurrentUser(req.Context(), cookie.Value)
	if err != nil {
		return uuid.UUID{}, false
	}
	return dto.ID, true
}

type createPlanRequest struct {
	Title string `json:"title"`
}

type planResponse struct {
	ID         string         `json:"id"`
	ShareToken string         `json:"share_token"`
	Title      string         `json:"title"`
	IsPublic   bool           `json:"is_public"`
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
		IsPublic:   plan.IsPublic,
		Pins:       pins,
		CreatedAt:  plan.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:  plan.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

type planSummaryResponse struct {
	ID         string `json:"id"`
	ShareToken string `json:"share_token"`
	Title      string `json:"title"`
	IsPublic   bool   `json:"is_public"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func toPlanSummaryResponse(plan *domain.Plan) planSummaryResponse {
	return planSummaryResponse{
		ID:         plan.ID.String(),
		ShareToken: plan.ShareToken,
		Title:      plan.Title,
		IsPublic:   plan.IsPublic,
		CreatedAt:  plan.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:  plan.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func (planHandler *PlanHandler) Create(rw http.ResponseWriter, req *http.Request) {
	var body createPlanRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil || body.Title == "" {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	plan, err := planHandler.planUseCase.CreatePlan(req.Context(), body.Title)
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	// Best-effort: if the request is authenticated, add the creator as a
	// PlanMember so the plan shows up in their /plans list. This must never
	// fail Plan creation -- most callers are anonymous by design.
	if userID, ok := planHandler.currentUserID(req); ok {
		_, _ = planHandler.planMemberUseCase.AddMember(req.Context(), plan.ID, userID)
	}

	writeJSON(rw, http.StatusCreated, toPlanResponse(plan, []pinWithNotes{}))
}

func (planHandler *PlanHandler) GetByShareToken(rw http.ResponseWriter, req *http.Request) {
	token := chi.URLParam(req, "share_token")
	ctx := req.Context()

	plan, err := planHandler.planUseCase.GetPlanByShareToken(ctx, token)
	if errors.Is(err, domain.ErrNotFound) {
		writeError(rw, http.StatusNotFound, "plan not found")
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	rawPins, err := planHandler.pinUseCase.ListPins(ctx, plan.ID)
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	pins := make([]pinWithNotes, 0, len(rawPins))
	for _, pin := range rawPins {
		rawNotes, err := planHandler.noteUseCase.ListNotes(ctx, pin.ID)
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

func (planHandler *PlanHandler) ListForUser(rw http.ResponseWriter, req *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(req, "id"))
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid id")
		return
	}

	plans, err := planHandler.planUseCase.ListPlansForUser(req.Context(), userID)
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := make([]planSummaryResponse, 0, len(plans))
	for _, plan := range plans {
		resp = append(resp, toPlanSummaryResponse(plan))
	}
	writeJSON(rw, http.StatusOK, resp)
}

func (planHandler *PlanHandler) Publish(rw http.ResponseWriter, req *http.Request) {
	token := chi.URLParam(req, "share_token")

	plan, err := planHandler.planUseCase.PublishPlan(req.Context(), token)
	if errors.Is(err, domain.ErrNotFound) {
		writeError(rw, http.StatusNotFound, "plan not found")
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(rw, http.StatusOK, toPlanResponse(plan, []pinWithNotes{}))
}
