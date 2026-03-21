package handlers

import (
	"encoding/json"
	"groupie-tracker/models"
	"groupie-tracker/services"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// SearchHandler handles GET /api/search requests, filtering and sorting artists by query params.
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))
	sortBy := r.URL.Query().Get("sort")
	minYear := r.URL.Query().Get("minYear")
	maxYear := r.URL.Query().Get("maxYear")
	minAlbumYear := r.URL.Query().Get("minAlbumYear")
	maxAlbumYear := r.URL.Query().Get("maxAlbumYear")
	membersParam := r.URL.Query().Get("members")
	locationsParam := r.URL.Query().Get("locations")

	// Parse member count filter: comma-separated ints like "1,2,4"
	var memberCounts []int
	if membersParam != "" {
		for _, s := range strings.Split(membersParam, ",") {
			s = strings.TrimSpace(s)
			if n, err := strconv.Atoi(s); err == nil {
				memberCounts = append(memberCounts, n)
			}
		}
	}

	// Parse location filter: comma-separated location keys
	var selectedLocations []string
	if locationsParam != "" {
		for _, s := range strings.Split(locationsParam, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				selectedLocations = append(selectedLocations, s)
			}
		}
	}

	artists, err := services.GetArtists()
	if err != nil {
		log.Printf("Error fetching artists for search: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch artists"}) // JS client parses JSON even on error, so we respond with JSON not plain text
		return
	}

	// Build relation map for location filtering and search
	var relationMap map[int]*models.Relation
	relations, relErr := services.GetAllRelations()
	if relErr != nil {
		log.Printf("Error fetching relations for search: %v", relErr)
		// Graceful degradation: skip location filter if relations fail
	} else {
		relationMap = make(map[int]*models.Relation, len(relations))
		for i := range relations {
			relationMap[relations[i].ID] = &relations[i]
		}
	}

	// Filter first, then sort—avoids sorting entries that will be discarded
	results := make([]models.Artist, 0)
	for _, artist := range artists {
		rel := relationMap[artist.ID] // may be nil if relations failed

		if !matchesSearch(artist, rel, query) {
			continue
		}
		if !matchesYearFilter(artist, minYear, maxYear) {
			continue
		}
		if !matchesAlbumYearFilter(artist, minAlbumYear, maxAlbumYear) {
			continue
		}
		if !matchesMembersFilter(artist, memberCounts) {
			continue
		}
		if !matchesLocationFilter(rel, selectedLocations) {
			continue
		}
		results = append(results, artist)
	}

	sortArtists(results, sortBy)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("Error encoding search results: %v", err)
	}
}

// matchesSearch returns true if the query matches the artist's name, members, creation date,
// first album date, or any concert location (case-insensitive).
func matchesSearch(artist models.Artist, relation *models.Relation, query string) bool {
	if query == "" {
		return true // Empty query means no text filter, so all artists pass
	}

	if strings.Contains(strings.ToLower(artist.Name), query) {
		return true
	}

	// Search member names so users can find bands by member (e.g., "Freddie" finds Queen)
	for _, member := range artist.Members {
		if strings.Contains(strings.ToLower(member), query) {
			return true
		}
	}

	// Search by creation date (e.g., "1973" finds AC/DC)
	if strings.Contains(strconv.Itoa(artist.CreationDate), query) {
		return true
	}

	// Search by first album date (e.g., "05-08-1967" finds Pink Floyd)
	if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
		return true
	}

	// Search by concert locations — match both raw key and formatted name
	// so "playa del carmen, mexico" and "playa_del_carmen-mexico" both work
	if relation != nil {
		for location := range relation.DatesLocations {
			if strings.Contains(strings.ToLower(location), query) {
				return true
			}
			city, country := services.FormatLocationName(location)
			formatted := strings.ToLower(city)
			if country != "" {
				formatted += ", " + strings.ToLower(country)
			}
			if strings.Contains(formatted, query) {
				return true
			}
		}
	}

	return false
}

// matchesYearFilter returns true if the artist's creation year falls within the given range.
// Empty strings mean no bound.
func matchesYearFilter(artist models.Artist, minYear, maxYear string) bool {
	if minYear != "" {
		if min, err := strconv.Atoi(minYear); err == nil && artist.CreationDate < min {
			return false
		}
	}

	if maxYear != "" {
		if max, err := strconv.Atoi(maxYear); err == nil && artist.CreationDate > max {
			return false
		}
	}

	return true
}

// parseFirstAlbumYear extracts the year from a "dd-mm-yyyy" format string.
func parseFirstAlbumYear(firstAlbum string) int {
	parts := strings.Split(firstAlbum, "-")
	if len(parts) < 3 {
		return 0
	}
	year, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 0
	}
	return year
}

// matchesAlbumYearFilter returns true if the artist's first album year falls within the given range.
func matchesAlbumYearFilter(artist models.Artist, minAlbumYear, maxAlbumYear string) bool {
	if minAlbumYear == "" && maxAlbumYear == "" {
		return true
	}

	year := parseFirstAlbumYear(artist.FirstAlbum)
	if year == 0 {
		return false // Malformed date can't match a year range
	}

	if minAlbumYear != "" {
		if min, err := strconv.Atoi(minAlbumYear); err == nil && year < min {
			return false
		}
	}

	if maxAlbumYear != "" {
		if max, err := strconv.Atoi(maxAlbumYear); err == nil && year > max {
			return false
		}
	}

	return true
}

// matchesMembersFilter returns true if the artist's member count is in the selected list.
// Empty slice means no filter. Value 8 means "8 or more".
func matchesMembersFilter(artist models.Artist, memberCounts []int) bool {
	if len(memberCounts) == 0 {
		return true
	}

	count := len(artist.Members)
	for _, mc := range memberCounts {
		if mc == 8 && count >= 8 {
			return true
		}
		if count == mc {
			return true
		}
	}
	return false
}

// matchesLocationFilter returns true if any of the artist's locations match the selected locations.
// Uses suffix matching for parent-region filtering (e.g., selecting "usa" matches "north_carolina-usa").
func matchesLocationFilter(relation *models.Relation, selectedLocations []string) bool {
	if len(selectedLocations) == 0 {
		return true
	}
	if relation == nil {
		return false
	}

	for location := range relation.DatesLocations {
		loc := strings.ToLower(location)
		for _, selected := range selectedLocations {
			sel := strings.ToLower(selected)
			if loc == sel {
				return true
			}
			// Parent-region matching: "usa" matches "north_carolina-usa"
			if strings.HasSuffix(loc, "-"+sel) {
				return true
			}
		}
	}
	return false
}

// sortArtists sorts the slice in-place by the given criteria: "newest", "oldest", "name", or default (by ID).
func sortArtists(artists []models.Artist, sortBy string) {
	sort.Slice(artists, func(i, j int) bool {
		switch sortBy {
		case "newest":
			return artists[i].CreationDate > artists[j].CreationDate
		case "oldest":
			return artists[i].CreationDate < artists[j].CreationDate
		case "name":
			return artists[i].Name < artists[j].Name
		default:
			return artists[i].ID < artists[j].ID // Default: original API order (by ID)
		}
	})
}
