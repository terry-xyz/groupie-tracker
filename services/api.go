package services

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/models"
	"net/http"
	"sync"
	"time"
)

const baseURL = "https://groupietrackers.herokuapp.com/api"

// In-memory cache to avoid hammering the external API on every request
var (
	cacheMu        sync.Mutex
	cachedArtists  []models.Artist
	artistsCacheAt time.Time
	cachedRelations []models.Relation
	relationsCacheAt time.Time
	cacheTTL       = 5 * time.Minute
)

// FetchArtists retrieves all artists from the Groupie Trackers API (/api/artists).
func FetchArtists() ([]models.Artist, error) {
	client := &http.Client{Timeout: 10 * time.Second} // Timeout prevents hanging if the external API is slow or unreachable
	resp, err := client.Get(baseURL + "/artists")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var artists []models.Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, fmt.Errorf("failed to parse artist data: %v", err)
	}
	return artists, nil
}

// FetchRelation fetches all relations then filters by ID because the API has no single-relation endpoint
func FetchRelation(id int) (*models.Relation, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(baseURL + "/relation")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	// API wraps the relation array in {"index": [...]}, so we decode into a matching struct
	var relations struct {
		Index []models.Relation `json:"index"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&relations); err != nil {
		return nil, fmt.Errorf("failed to parse relation data: %v", err)
	}

	// Linear search through all relations to find the one matching our artist ID
	for _, rel := range relations.Index {
		if rel.ID == id {
			return &rel, nil
		}
	}
	return nil, fmt.Errorf("relation not found for artist ID %d", id)
}

// FetchAllRelations retrieves all relations from the API.
func FetchAllRelations() ([]models.Relation, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(baseURL + "/relation")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var relations struct {
		Index []models.Relation `json:"index"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&relations); err != nil {
		return nil, fmt.Errorf("failed to parse relation data: %v", err)
	}
	return relations.Index, nil
}

// GetArtists returns cached artists, fetching fresh data if cache is expired or empty.
func GetArtists() ([]models.Artist, error) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	if cachedArtists != nil && time.Since(artistsCacheAt) < cacheTTL {
		return cachedArtists, nil
	}

	artists, err := FetchArtists()
	if err != nil {
		return nil, err
	}
	cachedArtists = artists
	artistsCacheAt = time.Now()
	return artists, nil
}

// GetAllRelations returns cached relations, fetching fresh data if cache is expired or empty.
func GetAllRelations() ([]models.Relation, error) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	if cachedRelations != nil && time.Since(relationsCacheAt) < cacheTTL {
		return cachedRelations, nil
	}

	relations, err := FetchAllRelations()
	if err != nil {
		return nil, err
	}
	cachedRelations = relations
	relationsCacheAt = time.Now()
	return relations, nil
}
