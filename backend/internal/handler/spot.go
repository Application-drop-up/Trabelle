package handler

import (
	"errors"
	"net/http"

	"github.com/Application-drop-up/Travellle/internal/domain/spot"
	spotuc "github.com/Application-drop-up/Travellle/internal/usecase/spot"
)

type SpotHandler struct {
	uc *spotuc.UseCase
}

func NewSpotHandler(uc *spotuc.UseCase) *SpotHandler {
	return &SpotHandler{uc: uc}
}

type spotResponse struct {
	PlaceID   string  `json:"place_id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func toSpotResponse(s *spot.Spot) spotResponse {
	return spotResponse{
		PlaceID:   s.PlaceID.String(),
		Name:      s.Name,
		Address:   s.Address,
		Latitude:  s.Location.Latitude,
		Longitude: s.Location.Longitude,
	}
}

func (h *SpotHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		writeError(w, http.StatusBadRequest, "query parameter is required")
		return
	}

	spots, err := h.uc.SearchSpots(r.Context(), query)
	if errors.Is(err, spot.ErrInvalidQuery) {
		writeError(w, http.StatusBadRequest, "query parameter is required")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := make([]spotResponse, 0, len(spots))
	for _, s := range spots {
		resp = append(resp, toSpotResponse(s))
	}
	writeJSON(w, http.StatusOK, resp)
}
