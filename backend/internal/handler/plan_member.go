package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	domain "github.com/Application-drop-up/Travellle/internal/domain/planmember"
	plandomain "github.com/Application-drop-up/Travellle/internal/domain/plan"
	userdomain "github.com/Application-drop-up/Travellle/internal/domain/user"
	planmemberuc "github.com/Application-drop-up/Travellle/internal/usecase/planmember"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PlanMemberHandler struct {
	useCase *planmemberuc.UseCase
}

func NewPlanMemberHandler(useCase *planmemberuc.UseCase) *PlanMemberHandler {
	return &PlanMemberHandler{useCase: useCase}
}

type addMemberRequest struct {
	UserID string `json:"user_id"`
}

type planMemberResponse struct {
	ID        string `json:"id"`
	PlanID    string `json:"plan_id"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

func toPlanMemberResponse(member *domain.PlanMember) planMemberResponse {
	return planMemberResponse{
		ID:        member.ID.String(),
		PlanID:    member.PlanID.String(),
		UserID:    member.UserID.String(),
		CreatedAt: member.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func (h *PlanMemberHandler) Add(rw http.ResponseWriter, req *http.Request) {
	planID, err := uuid.Parse(chi.URLParam(req, "plan_id"))
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid plan_id")
		return
	}

	var body addMemberRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, err := uuid.Parse(body.UserID)
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid user_id")
		return
	}

	member, err := h.useCase.AddMember(req.Context(), planID, userID)
	if errors.Is(err, plandomain.ErrNotFound) {
		writeError(rw, http.StatusNotFound, "plan not found")
		return
	}
	if errors.Is(err, userdomain.ErrNotFound) {
		writeError(rw, http.StatusNotFound, "user not found")
		return
	}
	if errors.Is(err, domain.ErrAlreadyMember) {
		writeError(rw, http.StatusConflict, "user is already a member of this plan")
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(rw, http.StatusCreated, toPlanMemberResponse(member))
}

func (h *PlanMemberHandler) List(rw http.ResponseWriter, req *http.Request) {
	planID, err := uuid.Parse(chi.URLParam(req, "plan_id"))
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid plan_id")
		return
	}

	members, err := h.useCase.ListMembers(req.Context(), planID)
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := make([]planMemberResponse, 0, len(members))
	for _, member := range members {
		resp = append(resp, toPlanMemberResponse(member))
	}
	writeJSON(rw, http.StatusOK, resp)
}

func (h *PlanMemberHandler) Remove(rw http.ResponseWriter, req *http.Request) {
	planID, err := uuid.Parse(chi.URLParam(req, "plan_id"))
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid plan_id")
		return
	}
	userID, err := uuid.Parse(chi.URLParam(req, "user_id"))
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid user_id")
		return
	}

	err = h.useCase.RemoveMember(req.Context(), planID, userID)
	if errors.Is(err, domain.ErrNotFound) {
		writeError(rw, http.StatusNotFound, "plan member not found")
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
