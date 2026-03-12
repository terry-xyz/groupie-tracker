package handlers

import (
	"encoding/json"
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

// TestSuggestionsHandlerCap verifies max 10 suggestions returned
func TestSuggestionsHandlerCap(t *testing.T) {
	// "a" should match many artists/members
	req := httptest.NewRequest("GET", "/api/suggestions?q=a", nil)
	w := httptest.NewRecorder()

	SuggestionsHandler(w, req)

	var suggestions []Suggestion
	if err := json.NewDecoder(w.Body).Decode(&suggestions); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(suggestions) > 10 {
		t.Errorf("Expected max 10 suggestions, got %d", len(suggestions))
	}
}
