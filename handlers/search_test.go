package handlers

import (
	"groupie-tracker/models"
	"testing"
)

// Verifies text search matches against artist name, member names, creation date, first album, and locations
func TestMatchesSearch(t *testing.T) {
	artist := models.Artist{
		ID:           1,
		Name:         "Queen",
		Members:      []string{"Freddie Mercury", "Brian May"},
		CreationDate: 1970,
		FirstAlbum:   "13-07-1973",
	}
	relation := &models.Relation{
		ID: 1,
		DatesLocations: map[string][]string{
			"london-uk": {"01-01-1975"},
		},
	}

	tests := []struct {
		query    string
		expected bool
	}{
		{"queen", true},
		{"freddie", true},
		{"brian", true},
		{"beatles", false},
		{"", true},
		{"1970", true},          // creation date search
		{"13-07-1973", true},    // first album search
		{"london", true},        // location search
		{"tokyo", false},        // location not present
	}

	for _, test := range tests {
		result := matchesSearch(artist, relation, test.query)
		if result != test.expected {
			t.Errorf("matchesSearch(%q) = %v, expected %v", test.query, result, test.expected)
		}
	}
}

// Verifies search works with nil relation (no location data)
func TestMatchesSearchNilRelation(t *testing.T) {
	artist := models.Artist{
		ID:           1,
		Name:         "Queen",
		Members:      []string{"Freddie Mercury"},
		CreationDate: 1970,
		FirstAlbum:   "13-07-1973",
	}

	result := matchesSearch(artist, nil, "queen")
	if !result {
		t.Error("Expected true for name match with nil relation")
	}

	result = matchesSearch(artist, nil, "london")
	if result {
		t.Error("Expected false for location search with nil relation")
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

// Verifies first album year parsing from "dd-mm-yyyy" format
func TestParseFirstAlbumYear(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"14-02-1992", 1992},
		{"05-08-1967", 1967},
		{"01-01-2000", 2000},
		{"bad-data", 0},
		{"", 0},
		{"only-two", 0},
	}

	for _, test := range tests {
		result := parseFirstAlbumYear(test.input)
		if result != test.expected {
			t.Errorf("parseFirstAlbumYear(%q) = %d, expected %d", test.input, result, test.expected)
		}
	}
}

// Verifies first album year range filtering
func TestMatchesAlbumYearFilter(t *testing.T) {
	artist := models.Artist{
		FirstAlbum: "14-02-1992",
	}

	tests := []struct {
		minAlbumYear string
		maxAlbumYear string
		expected     bool
	}{
		{"1990", "1995", true},
		{"1993", "2000", false},
		{"", "1995", true},
		{"1990", "", true},
		{"", "", true},
	}

	for _, test := range tests {
		result := matchesAlbumYearFilter(artist, test.minAlbumYear, test.maxAlbumYear)
		if result != test.expected {
			t.Errorf("matchesAlbumYearFilter(min=%s, max=%s) = %v, expected %v",
				test.minAlbumYear, test.maxAlbumYear, result, test.expected)
		}
	}

	// Test malformed date with filter active
	malformed := models.Artist{FirstAlbum: "bad-data"}
	if matchesAlbumYearFilter(malformed, "1990", "") {
		t.Error("Expected false for malformed date with active filter")
	}
}

// Verifies member count filter including 8+ logic
func TestMatchesMembersFilter(t *testing.T) {
	tests := []struct {
		memberCount  int
		memberCounts []int
		expected     bool
	}{
		{4, []int{4}, true},
		{4, []int{1, 2, 3}, false},
		{4, []int{}, true},         // Empty = no filter
		{8, []int{8}, true},        // 8+ matches exactly 8
		{10, []int{8}, true},       // 8+ matches >=8
		{7, []int{8}, false},       // 7 doesn't match 8+
		{3, []int{1, 3, 5}, true},
	}

	for _, test := range tests {
		members := make([]string, test.memberCount)
		artist := models.Artist{Members: members}
		result := matchesMembersFilter(artist, test.memberCounts)
		if result != test.expected {
			t.Errorf("matchesMembersFilter(count=%d, filter=%v) = %v, expected %v",
				test.memberCount, test.memberCounts, result, test.expected)
		}
	}
}

// Verifies location filter with parent-region matching
func TestMatchesLocationFilter(t *testing.T) {
	relation := &models.Relation{
		DatesLocations: map[string][]string{
			"north_carolina-usa": {"01-01-2020"},
			"london-uk":          {"02-02-2020"},
		},
	}

	tests := []struct {
		selected []string
		expected bool
	}{
		{[]string{}, true},                       // Empty = no filter
		{[]string{"north_carolina-usa"}, true},   // Exact match
		{[]string{"usa"}, true},                  // Parent-region match
		{[]string{"uk"}, true},                   // Parent-region match
		{[]string{"japan"}, false},               // No match
		{[]string{"london-uk"}, true},            // Exact match
	}

	for _, test := range tests {
		result := matchesLocationFilter(relation, test.selected)
		if result != test.expected {
			t.Errorf("matchesLocationFilter(selected=%v) = %v, expected %v",
				test.selected, result, test.expected)
		}
	}

	// Nil relation with active filter
	if matchesLocationFilter(nil, []string{"usa"}) {
		t.Error("Expected false for nil relation with active filter")
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
