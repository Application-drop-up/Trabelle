package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Application-drop-up/Travellle/internal/domain/spot"
)

const (
	defaultSearchURL = "https://places.googleapis.com/v1/places:searchText"
	fieldMask        = "places.id,places.displayName,places.formattedAddress,places.location"
)

type GooglePlacesClient struct {
	apiKey     string
	searchURL  string
	httpClient *http.Client
}

func NewGooglePlacesClient(apiKey string) *GooglePlacesClient {
	return NewGooglePlacesClientWithURL(defaultSearchURL, apiKey)
}

// NewGooglePlacesClientWithURL allows overriding the endpoint URL for testing.
func NewGooglePlacesClientWithURL(searchURL, apiKey string) *GooglePlacesClient {
	return &GooglePlacesClient{
		apiKey:     apiKey,
		searchURL:  searchURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// --- request / response structs ---

type searchTextRequest struct {
	TextQuery string `json:"textQuery"`
}

type searchTextResponse struct {
	Places []placeResult `json:"places"`
}

type placeResult struct {
	ID               string      `json:"id"`
	DisplayName      displayName `json:"displayName"`
	FormattedAddress string      `json:"formattedAddress"`
	Location         latLng      `json:"location"`
}

type displayName struct {
	Text string `json:"text"`
}

type latLng struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// --- Searcher implementation ---

func (c *GooglePlacesClient) Search(ctx context.Context, query string) ([]*spot.Spot, error) {
	if query == "" {
		return nil, spot.ErrInvalidQuery
	}

	body, err := json.Marshal(searchTextRequest{TextQuery: query})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.searchURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", c.apiKey)
	req.Header.Set("X-Goog-FieldMask", fieldMask)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call places api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("places api returned status %d", resp.StatusCode)
	}

	var result searchTextResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return toSpots(result.Places), nil
}

func toSpots(places []placeResult) []*spot.Spot {
	spots := make([]*spot.Spot, 0, len(places))
	for _, place := range places {
		placeID, err := spot.NewPlaceID(place.ID)
		if err != nil {
			continue
		}
		location, err := spot.NewLocation(place.Location.Latitude, place.Location.Longitude)
		if err != nil {
			continue
		}
		spots = append(spots, &spot.Spot{
			PlaceID:  placeID,
			Name:     place.DisplayName.Text,
			Address:  place.FormattedAddress,
			Location: location,
		})
	}
	return spots
}
