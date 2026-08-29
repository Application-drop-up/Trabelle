package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	plandomain "github.com/Application-drop-up/Travellle/internal/domain/plan"
	"github.com/Application-drop-up/Travellle/internal/domain/spot"
	spotuc "github.com/Application-drop-up/Travellle/internal/usecase/spot"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SpotHandler struct {
	useCase *spotuc.UseCase
}

func NewSpotHandler(useCase *spotuc.UseCase) *SpotHandler {
	return &SpotHandler{useCase: useCase}
}

type saveSpotRequest struct {
	PlaceID   string  `json:"place_id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	PlanID    string  `json:"plan_id"`
}

type spotResponse struct {
	PlaceID   string  `json:"place_id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func toSpotResponse(spotItem *spot.Spot) spotResponse {
	return spotResponse{
		PlaceID:   spotItem.PlaceID.String(),
		Name:      spotItem.Name,
		Address:   spotItem.Address,
		Latitude:  spotItem.Location.Latitude,
		Longitude: spotItem.Location.Longitude,
	}
}

func (spotHandler *SpotHandler) Search(rw http.ResponseWriter, req *http.Request) {
	query := req.URL.Query().Get("query")
	if query == "" {
		writeError(rw, http.StatusBadRequest, "query parameter is required")
		return
	}

	spots, err := spotHandler.useCase.SearchSpots(req.Context(), query)
	if errors.Is(err, spot.ErrInvalidQuery) {
		writeError(rw, http.StatusBadRequest, "query parameter is required")
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := make([]spotResponse, 0, len(spots))
	for _, spotItem := range spots {
		resp = append(resp, toSpotResponse(spotItem))
	}
	writeJSON(rw, http.StatusOK, resp)
}

func (spotHandler *SpotHandler) Save(rw http.ResponseWriter, req *http.Request) {
	if _, err := uuid.Parse(chi.URLParam(req, "id")); err != nil {
		writeError(rw, http.StatusBadRequest, "invalid id")
		return
	}

	var body saveSpotRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Name == "" || body.Address == "" {
		writeError(rw, http.StatusBadRequest, "invalid request body")
		return
	}

	planID, err := uuid.Parse(body.PlanID)
	if err != nil {
		writeError(rw, http.StatusBadRequest, "invalid plan_id")
		return
	}

	placeID, err := spot.NewPlaceID(body.PlaceID)
	if err != nil {
		writeError(rw, http.StatusBadRequest, err.Error())
		return
	}

	location, err := spot.NewLocation(body.Latitude, body.Longitude)
	if err != nil {
		writeError(rw, http.StatusBadRequest, err.Error())
		return
	}

	newSpot := &spot.Spot{
		ID:          uuid.New(),
		PlaceID:     placeID,
		Name:        body.Name,
		Address:     body.Address,
		Location:    location,
		FirstPlanID: planID,
	}

	if err := spotHandler.useCase.SaveSpot(req.Context(), newSpot); err != nil {
		if errors.Is(err, plandomain.ErrNotFound) {
			writeError(rw, http.StatusNotFound, "plan not found")
			return
		}
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(rw, http.StatusCreated, toSpotResponse(newSpot))
}
