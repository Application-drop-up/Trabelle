package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Application-drop-up/Travellle/internal/domain/spot"
)

const defaultTomTomSearchURL = "https://api.tomtom.com/search/2/search"

type TomTomClient struct {
	apiKey     string
	searchURL  string
	httpClient *http.Client
}

func NewTomTomClient(apiKey string) *TomTomClient {
	return NewTomTomClientWithURL(defaultTomTomSearchURL, apiKey)
}

// NewTomTomClientWithURL allows overriding the endpoint URL for testing.
func NewTomTomClientWithURL(searchURL, apiKey string) *TomTomClient {
	return &TomTomClient{
		apiKey:     apiKey,
		searchURL:  searchURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// --- response structs ---

type tomTomSearchResponse struct {
	Results []tomTomResult `json:"results"`
}

type tomTomResult struct {
	ID       string       `json:"id"`
	POI      tomTomPOI    `json:"poi"`
	Address  tomTomAddr   `json:"address"`
	Position tomTomLatLon `json:"position"`
}

type tomTomPOI struct {
	Name string `json:"name"`
}

type tomTomAddr struct {
	FreeformAddress string `json:"freeformAddress"`
}

type tomTomLatLon struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// --- Searcher implementation ---

func (c *TomTomClient) Search(ctx context.Context, query string) ([]*spot.Spot, error) {
	if query == "" {
		return nil, spot.ErrInvalidQuery
	}

	reqURL := fmt.Sprintf("%s/%s.json?key=%s", c.searchURL, url.PathEscape(query), url.QueryEscape(c.apiKey))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call tomtom search api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tomtom search api returned status %d", resp.StatusCode)
	}

	var result tomTomSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return toSpotsFromTomTom(result.Results), nil
}

func toSpotsFromTomTom(results []tomTomResult) []*spot.Spot {
	spots := make([]*spot.Spot, 0, len(results))
	for _, result := range results {
		placeID, err := spot.NewPlaceID(result.ID)
		if err != nil {
			continue
		}
		location, err := spot.NewLocation(result.Position.Lat, result.Position.Lon)
		if err != nil {
			continue
		}
		spots = append(spots, &spot.Spot{
			PlaceID:  placeID,
			Name:     result.POI.Name,
			Address:  result.Address.FreeformAddress,
			Location: location,
		})
	}
	return spots
}
