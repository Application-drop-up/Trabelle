package countryguide

import (
	domain "github.com/Application-drop-up/Travellle/internal/domain/countryguide"
	"github.com/google/uuid"
)

type ItemDto struct {
	ID          uuid.UUID
	Category    string
	Title       string
	Description string
	URL         string
	IsMandatory bool
}

type CountryGuideDto struct {
	ID          uuid.UUID
	CountryCode string
	CountryName string
	Items       []ItemDto
}

func NewCountryGuideDto(guide *domain.CountryGuide) *CountryGuideDto {
	items := make([]ItemDto, len(guide.Items))
	for i, item := range guide.Items {
		items[i] = ItemDto{
			ID:          item.ID,
			Category:    string(item.Category),
			Title:       item.Title,
			Description: item.Description,
			URL:         item.URL,
			IsMandatory: item.IsMandatory,
		}
	}

	return &CountryGuideDto{
		ID:          guide.ID,
		CountryCode: guide.CountryCode,
		CountryName: guide.CountryName,
		Items:       items,
	}
}
