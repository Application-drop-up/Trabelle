package handler

import (
	"errors"
	"net/http"

	domain "github.com/Application-drop-up/Travellle/internal/domain/countryguide"
	countryguideuc "github.com/Application-drop-up/Travellle/internal/usecase/countryguide"
	"github.com/go-chi/chi/v5"
)

type CountryGuideHandler struct {
	useCase *countryguideuc.UseCase
}

func NewCountryGuideHandler(useCase *countryguideuc.UseCase) *CountryGuideHandler {
	return &CountryGuideHandler{useCase: useCase}
}

type countryGuideItemResponse struct {
	ID          string `json:"id"`
	Category    string `json:"category"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	IsMandatory bool   `json:"is_mandatory"`
}

type countryGuideResponse struct {
	ID          string                     `json:"id"`
	CountryCode string                     `json:"country_code"`
	CountryName string                     `json:"country_name"`
	Items       []countryGuideItemResponse `json:"items"`
}

func toCountryGuideResponse(dto *countryguideuc.CountryGuideDto) countryGuideResponse {
	items := make([]countryGuideItemResponse, len(dto.Items))
	for i, item := range dto.Items {
		items[i] = countryGuideItemResponse{
			ID:          item.ID.String(),
			Category:    item.Category,
			Title:       item.Title,
			Description: item.Description,
			URL:         item.URL,
			IsMandatory: item.IsMandatory,
		}
	}

	return countryGuideResponse{
		ID:          dto.ID.String(),
		CountryCode: dto.CountryCode,
		CountryName: dto.CountryName,
		Items:       items,
	}
}

func (countryGuideHandler *CountryGuideHandler) List(rw http.ResponseWriter, req *http.Request) {
	guides, err := countryGuideHandler.useCase.ListCountryGuides(req.Context())
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := make([]countryGuideResponse, len(guides))
	for i, dto := range guides {
		resp[i] = toCountryGuideResponse(dto)
	}
	writeJSON(rw, http.StatusOK, resp)
}

func (countryGuideHandler *CountryGuideHandler) GetByCode(rw http.ResponseWriter, req *http.Request) {
	code := chi.URLParam(req, "code")

	dto, err := countryGuideHandler.useCase.GetCountryGuide(req.Context(), code)
	if errors.Is(err, domain.ErrNotFound) {
		writeError(rw, http.StatusNotFound, "country guide not found")
		return
	}
	if err != nil {
		writeError(rw, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(rw, http.StatusOK, toCountryGuideResponse(dto))
}
