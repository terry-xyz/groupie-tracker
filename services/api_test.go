package services

import (
	"testing"
)

// Integration test: hits the real Groupie Trackers API to verify connectivity and response shape
func TestFetchArtists(t *testing.T) {
	artists, err := FetchArtists()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(artists) == 0 {
		t.Error("Expected artists, got empty array")
	}

	for _, artist := range artists {
		if artist.ID == 0 {
			t.Error("Artist ID should not be 0")
		}
		if artist.Name == "" {
			t.Error("Artist name should not be empty")
		}
	}
}

// Integration test: verifies relation fetch and ID filtering for a known artist (ID 1)
func TestFetchRelation(t *testing.T) {
	relation, err := FetchRelation(1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if relation == nil {
		t.Error("Expected relation data, got nil")
	}

	if relation != nil && relation.ID != 1 {
		t.Errorf("Expected relation ID 1, got %d", relation.ID)
	}
}
