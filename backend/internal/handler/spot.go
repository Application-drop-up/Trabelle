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

func toSpotResponse(spotItem *spot.Spot) spotResponse {
	return spotResponse{
		PlaceID:   spotItem.PlaceID.String(),
		Name:      spotItem.Name,
		Address:   spotItem.Address,
		Latitude:  spotItem.Location.Latitude,
		Longitude: spotItem.Location.Longitude,
	}
}

func (sh *SpotHandler) Search(rw http.ResponseWriter, req *http.Request) {
	query := req.URL.Query().Get("query")
	if query == "" {
		writeError(rw, http.StatusBadRequest, "query parameter is required")
		return
	}

	spots, err := sh.uc.SearchSpots(req.Context(), query)
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
