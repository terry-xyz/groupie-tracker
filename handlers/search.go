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

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))
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

func matchesSearch(artist models.Artist, query string) bool {
	if query == "" {
		return true
	}

	if strings.Contains(strings.ToLower(artist.Name), query) {
		return true
	}

	for _, member := range artist.Members {
		if strings.Contains(strings.ToLower(member), query) {
			return true
		}
	}

	return false
}

func matchesYearFilter(artist models.Artist, minYear, maxYear string) bool {
	if minYear != "" {
		min, _ := strconv.Atoi(minYear)
		if artist.CreationDate < min {
			return false
		}
	}

	if maxYear != "" {
		max, _ := strconv.Atoi(maxYear)
		if artist.CreationDate > max {
			return false
		}
	}

	return true
}

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
			return artists[i].ID < artists[j].ID
		}
	})
}
