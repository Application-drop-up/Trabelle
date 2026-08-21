package countryguide

import (
	domain "github.com/Application-drop-up/Travellle/internal/domain/countryguide"
)

type CountryGuideDtoCollection []*CountryGuideDto

func NewCountryGuideDtoCollection(guides []*domain.CountryGuide) CountryGuideDtoCollection {
	collection := make(CountryGuideDtoCollection, len(guides))
	for i, guide := range guides {
		collection[i] = NewCountryGuideDto(guide)
	}
	return collection
}
