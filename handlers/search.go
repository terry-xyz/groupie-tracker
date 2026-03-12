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

// SearchHandler handles GET /api/search requests, filtering and sorting artists by query params (q, sort, minYear, maxYear).
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q")) // Lowercased for case-insensitive matching later
	sortBy := r.URL.Query().Get("sort")
	minYear := r.URL.Query().Get("minYear")
	maxYear := r.URL.Query().Get("maxYear")

	artists, err := services.FetchArtists()
	if err != nil {
		log.Printf("Error fetching artists for search: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch artists"})
		return
	}

	// Filter first, then sort—avoids sorting entries that will be discarded
	var results []models.Artist
	for _, artist := range artists {
		if matchesSearch(artist, query) && matchesYearFilter(artist, minYear, maxYear) {
			results = append(results, artist)
		}
	}

	sortArtists(results, sortBy)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// matchesSearch returns true if the query matches the artist's name or any member name (case-insensitive).
func matchesSearch(artist models.Artist, query string) bool {
	if query == "" {
		return true // Empty query means no text filter, so all artists pass
	}

	if strings.Contains(strings.ToLower(artist.Name), query) {
		return true
	}

	// Also search member names so users can find bands by member (e.g., "Freddie" finds Queen)
	for _, member := range artist.Members {
		if strings.Contains(strings.ToLower(member), query) {
			return true
		}
	}

	return false
}

// matchesYearFilter returns true if the artist's creation year falls within the given range. Empty strings mean no bound.
func matchesYearFilter(artist models.Artist, minYear, maxYear string) bool {
	// Empty string means no bound on that side, so the filter is open-ended
	if minYear != "" {
		min, _ := strconv.Atoi(minYear)
		if artist.CreationDate < min { // Artist formed before the minimum year
			return false
		}
	}

	if maxYear != "" {
		max, _ := strconv.Atoi(maxYear)
		if artist.CreationDate > max { // Artist formed after the maximum year
			return false
		}
	}

	return true
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
