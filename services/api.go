package services

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/models"
	"net/http"
	"time"
)

const baseURL = "https://groupietrackers.herokuapp.com/api"

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
