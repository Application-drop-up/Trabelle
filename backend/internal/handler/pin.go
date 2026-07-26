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

func toPinResponse(pin *domain.Pin) pinResponse {
	return pinResponse{
		ID:        pin.ID.String(),
		PlanID:    pin.PlanID.String(),
		Name:      pin.Name,
		Latitude:  pin.Latitude,
		Longitude: pin.Longitude,
		Category:  string(pin.Category),
		Colour:    pin.Colour,
		CreatedAt: pin.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: pin.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func (ph *PinHandler) List(rw http.ResponseWriter, req *http.Request) {
	planID, err := uuid.Parse(chi.URLParam(req, "plan_id"))
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid plan_id")
		return
	}

	pins, err := ph.uc.ListPins(req.Context(), planID)
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := make([]pinResponse, 0, len(pins))
	for _, pin := range pins {
		resp = append(resp, toPinResponse(pin))
	}
	writeJSON(rw, http.StatusOK, resp)
}

func (ph *PinHandler) Create(rw http.ResponseWriter, req *http.Request) {
	planID, err := uuid.Parse(chi.URLParam(req, "plan_id"))
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid plan_id")
		return
	}

	var body createPinRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	cat := domain.Category(body.Category)
	if body.Name == "" || body.Colour == "" || !cat.IsValid() {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	pin, err := ph.uc.CreatePin(req.Context(), pinuc.CreateInput{
		PlanID:    planID,
		Name:      body.Name,
		Latitude:  body.Latitude,
		Longitude: body.Longitude,
		Category:  cat,
		Colour:    body.Colour,
	})
	if errors.Is(err, plandomain.ErrNotFound) {
		writeError(rw, http.StatusNotFound, "plan not found")
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(rw, http.StatusCreated, toPinResponse(pin))
}

func (ph *PinHandler) Update(rw http.ResponseWriter, req *http.Request) {
	planID, err := uuid.Parse(chi.URLParam(req, "plan_id"))
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid plan_id")
		return
	}
	pinID, err := uuid.Parse(chi.URLParam(req, "pin_id"))
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid pin_id")
		return
	}

	var body updatePinRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	input := pinuc.UpdateInput{}
	if body.Category != nil {
		cat := domain.Category(*body.Category)
		if !cat.IsValid() {
			writeError(rw, http.StatusBadRequest, "invalid category")
			return
		}
		input.Category = &cat
	}
	if body.Colour != nil {
		input.Colour = body.Colour
	}

	pin, err := ph.uc.UpdatePin(req.Context(), planID, pinID, input)
	if errors.Is(err, domain.ErrNotFound) {
		writeError(rw, http.StatusNotFound, "pin not found")
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(rw, http.StatusOK, toPinResponse(pin))
}

func (ph *PinHandler) Delete(rw http.ResponseWriter, req *http.Request) {
	planID, err := uuid.Parse(chi.URLParam(req, "plan_id"))
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid plan_id")
		return
	}
	pinID, err := uuid.Parse(chi.URLParam(req, "pin_id"))
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid pin_id")
		return
	}

	if err := ph.uc.DeletePin(req.Context(), planID, pinID); errors.Is(err, domain.ErrNotFound) {
		writeError(rw, http.StatusNotFound, "pin not found")
		return
	} else if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
