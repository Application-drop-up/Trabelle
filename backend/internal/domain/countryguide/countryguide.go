package countryguide

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNotFound            = errors.New("country guide not found")
	ErrInvalidCountryCode  = errors.New("country code must be a 2-letter ISO 3166-1 alpha-2 code")
	ErrEmptyCountryName    = errors.New("country name must not be empty")
	ErrEmptyItemTitle      = errors.New("item title must not be empty")
	ErrInvalidItemCategory = errors.New("invalid item category")
)

type ItemCategory string

const (
	CategoryEntryCard         ItemCategory = "entry_card"
	CategorySIMRecommendation ItemCategory = "sim_recommendation"
	CategoryPackingTip        ItemCategory = "packing_tip"
)

func (category ItemCategory) isValid() bool {
	switch category {
	case CategoryEntryCard, CategorySIMRecommendation, CategoryPackingTip:
		return true
	default:
		return false
	}
}

// Item is a single piece of guidance within a CountryGuide (e.g. one entry
// card link, one SIM recommendation). URL is optional -- not every item
// (e.g. a packing tip) has one.
type Item struct {
	ID          uuid.UUID
	Category    ItemCategory
	Title       string
	Description string
	URL         string
}

// CountryGuide is a read-only, seeded-by-us collection of practical
// pre-travel guidance for a single destination. It has no business
// invariants tied to end-user actions; NewCountryGuide exists to catch
// mistakes in our own seed data early, not to guard against untrusted input.
type CountryGuide struct {
	ID          uuid.UUID
	CountryCode string
	CountryName string
	Items       []Item
}

func NewCountryGuide(countryCode, countryName string, items []Item) (*CountryGuide, error) {
	if len(countryCode) != 2 {
		return nil, ErrInvalidCountryCode
	}
	if countryName == "" {
		return nil, ErrEmptyCountryName
	}
	for _, item := range items {
		if item.Title == "" {
			return nil, ErrEmptyItemTitle
		}
		if !item.Category.isValid() {
			return nil, ErrInvalidItemCategory
		}
	}

	return &CountryGuide{
		ID:          uuid.New(),
		CountryCode: countryCode,
		CountryName: countryName,
		Items:       items,
	}, nil
}
