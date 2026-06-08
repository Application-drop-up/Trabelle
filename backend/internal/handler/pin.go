package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	plandomain "github.com/Application-drop-up/Travellle/internal/domain/plan"
	domain "github.com/Application-drop-up/Travellle/internal/domain/pin"
	pinuc "github.com/Application-drop-up/Travellle/internal/usecase/pin"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PinHandler struct {
	uc *pinuc.UseCase
}

func NewPinHandler(uc *pinuc.UseCase) *PinHandler {
	return &PinHandler{uc: uc}
}

type createPinRequest struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Category  string  `json:"category"`
	Colour    string  `json:"colour"`
}

type updatePinRequest struct {
	Category *string `json:"category"`
	Colour   *string `json:"colour"`
}

type pinResponse struct {
	ID        string  `json:"id"`
	PlanID    string  `json:"plan_id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Category  string  `json:"category"`
	Colour    string  `json:"colour"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

func toPinResponse(p *domain.Pin) pinResponse {
	return pinResponse{
		ID:        p.ID.String(),
		PlanID:    p.PlanID.String(),
		Name:      p.Name,
		Latitude:  p.Latitude,
		Longitude: p.Longitude,
		Category:  string(p.Category),
		Colour:    p.Colour,
		CreatedAt: p.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: p.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func (h *PinHandler) List(w http.ResponseWriter, r *http.Request) {
	planID, err := uuid.Parse(chi.URLParam(r, "plan_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid plan_id")
		return
	}

	pins, err := h.uc.ListPins(r.Context(), planID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := make([]pinResponse, 0, len(pins))
	for _, p := range pins {
		resp = append(resp, toPinResponse(p))
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *PinHandler) Create(w http.ResponseWriter, r *http.Request) {
	planID, err := uuid.Parse(chi.URLParam(r, "plan_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid plan_id")
		return
	}

	var req createPinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	cat := domain.Category(req.Category)
	if req.Name == "" || req.Colour == "" || !cat.IsValid() {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p, err := h.uc.CreatePin(r.Context(), pinuc.CreateInput{
		PlanID:    planID,
		Name:      req.Name,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Category:  cat,
		Colour:    req.Colour,
	})
	if errors.Is(err, plandomain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "plan not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, toPinResponse(p))
}

func (h *PinHandler) Update(w http.ResponseWriter, r *http.Request) {
	planID, err := uuid.Parse(chi.URLParam(r, "plan_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid plan_id")
		return
	}
	pinID, err := uuid.Parse(chi.URLParam(r, "pin_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pin_id")
		return
	}

	var req updatePinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input := pinuc.UpdateInput{}
	if req.Category != nil {
		cat := domain.Category(*req.Category)
		if !cat.IsValid() {
			writeError(w, http.StatusBadRequest, "invalid category")
			return
		}
		input.Category = &cat
	}
	if req.Colour != nil {
		input.Colour = req.Colour
	}

	p, err := h.uc.UpdatePin(r.Context(), planID, pinID, input)
	if errors.Is(err, domain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "pin not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, toPinResponse(p))
}

func (h *PinHandler) Delete(w http.ResponseWriter, r *http.Request) {
	planID, err := uuid.Parse(chi.URLParam(r, "plan_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid plan_id")
		return
	}
	pinID, err := uuid.Parse(chi.URLParam(r, "pin_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pin_id")
		return
	}

	if err := h.uc.DeletePin(r.Context(), planID, pinID); errors.Is(err, domain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "pin not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
