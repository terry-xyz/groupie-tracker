package handlers

import (
	"encoding/json"
	"groupie-tracker/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestSuggestionsHandlerEmpty verifies empty query returns empty array
func TestSuggestionsHandlerEmpty(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/suggestions?q=", nil)
	w := httptest.NewRecorder()

	SuggestionsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var suggestions []Suggestion
	if err := json.NewDecoder(w.Body).Decode(&suggestions); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(suggestions) != 0 {
		t.Errorf("Expected 0 suggestions for empty query, got %d", len(suggestions))
	}
}

// TestSuggestionsHandlerResults verifies suggestions are returned with categories
func TestSuggestionsHandlerResults(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/suggestions?q=queen", nil)
	w := httptest.NewRecorder()

	SuggestionsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var suggestions []Suggestion
	if err := json.NewDecoder(w.Body).Decode(&suggestions); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(suggestions) == 0 {
		t.Error("Expected at least 1 suggestion for 'queen'")
	}

	// Check that all suggestions have non-empty text and category
	for _, s := range suggestions {
		if s.Text == "" {
			t.Error("Suggestion text should not be empty")
		}
		if s.Category == "" {
			t.Error("Suggestion category should not be empty")
		}
	}
}

// TestSuggestionsHandlerCap verifies max 15 suggestions returned
func TestSuggestionsHandlerCap(t *testing.T) {
	// "a" should match many artists/members
	req := httptest.NewRequest("GET", "/api/suggestions?q=a", nil)
	w := httptest.NewRecorder()

	SuggestionsHandler(w, req)

	var suggestions []Suggestion
	if err := json.NewDecoder(w.Body).Decode(&suggestions); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(suggestions) > 15 {
		t.Errorf("Expected max 15 suggestions, got %d", len(suggestions))
	}
}

func TestBuildSuggestionsBalancesCategories(t *testing.T) {
	artists := []models.Artist{
		{ID: 1, Name: "Genesis", Members: []string{"Gabriel"}, CreationDate: 1967, FirstAlbum: "01-01-1969"},
		{ID: 2, Name: "Ghost", Members: []string{"Ghoulette"}, CreationDate: 2006, FirstAlbum: "01-01-2010"},
	}
	relationMap := map[int]map[string][]string{
		1: {"glasgow-scotland": {"2026-01-01"}},
		2: {"geneva-switzerland": {"2026-01-02"}},
	}

	suggestions := buildSuggestions(artists, relationMap, "g", 10)

	hasCategory := func(category string) bool {
		for _, suggestion := range suggestions {
			if suggestion.Category == category {
				return true
			}
		}
		return false
	}

	if !hasCategory("artist/band") {
		t.Fatal("Expected artist/band suggestions for broad query")
	}
	if !hasCategory("location") {
		t.Fatal("Expected location suggestions for broad query")
	}
	if !hasCategory("member") {
		t.Fatal("Expected member suggestions for broad query")
	}
}
