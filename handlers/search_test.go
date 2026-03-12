package handlers

import (
	"groupie-tracker/models"
	"testing"
)

// Verifies text search matches against artist name and individual member names
func TestMatchesSearch(t *testing.T) {
	artist := models.Artist{
		ID:      1,
		Name:    "Queen",
		Members: []string{"Freddie Mercury", "Brian May"},
	}

	tests := []struct {
		query    string
		expected bool
	}{
		{"queen", true},
		{"freddie", true},
		{"brian", true},
		{"beatles", false},
		{"", true}, // Empty query should match all artists (no filter)
	}

	for _, test := range tests {
		result := matchesSearch(artist, test.query)
		if result != test.expected {
			t.Errorf("matchesSearch(%q) = %v, expected %v", test.query, result, test.expected)
		}
	}
}

// Verifies year range filter with optional min/max bounds
func TestMatchesYearFilter(t *testing.T) {
	artist := models.Artist{
		ID:           1,
		Name:         "Queen",
		CreationDate: 1970,
	}

	tests := []struct {
		minYear  string
		maxYear  string
		expected bool
	}{
		{"1960", "1980", true},  // 1970 is within [1960, 1980]
		{"1980", "2000", false}, // 1970 is before the 1980 minimum
		{"", "1980", true},      // No lower bound, only upper
		{"1960", "", true},      // No upper bound, only lower
		{"", "", true},          // No bounds at all, everything passes
	}

	for _, test := range tests {
		result := matchesYearFilter(artist, test.minYear, test.maxYear)
		if result != test.expected {
			t.Errorf("matchesYearFilter(min=%s, max=%s) = %v, expected %v",
				test.minYear, test.maxYear, result, test.expected)
		}
	}
}

// Verifies sort order for name, newest, and oldest sort modes
func TestSortArtists(t *testing.T) {
	artists := []models.Artist{
		{ID: 1, Name: "Queen", CreationDate: 1970},
		{ID: 2, Name: "Beatles", CreationDate: 1960},
		{ID: 3, Name: "Zeppelin", CreationDate: 1968},
	}

	sortArtists(artists, "name")
	if artists[0].Name != "Beatles" {
		t.Errorf("Expected Beatles first when sorted by name, got %s", artists[0].Name)
	}

	sortArtists(artists, "newest")
	if artists[0].CreationDate != 1970 {
		t.Errorf("Expected 1970 first when sorted newest, got %d", artists[0].CreationDate)
	}

	sortArtists(artists, "oldest")
	if artists[0].CreationDate != 1960 {
		t.Errorf("Expected 1960 first when sorted oldest, got %d", artists[0].CreationDate)
	}
}
