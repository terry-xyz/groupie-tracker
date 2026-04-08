package handlers

import (
	"encoding/json"
	"groupie-tracker/models"
	"groupie-tracker/services"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Suggestion represents an autocomplete suggestion with its source category.
type Suggestion struct {
	Text     string `json:"text"`
	Category string `json:"category"`
}

// SuggestionsHandler handles GET /api/suggestions?q=... for autocomplete.
func SuggestionsHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode([]Suggestion{}); err != nil {
			log.Printf("Error encoding empty suggestions: %v", err)
		}
		return
	}

	artists, err := services.GetArtists()
	if err != nil {
		log.Printf("Error fetching artists for suggestions: %v", err)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode([]Suggestion{}); err != nil {
			log.Printf("Error encoding empty suggestions on fetch error: %v", err)
		}
		return
	}

	relations, relErr := services.GetAllRelations()
	if relErr != nil {
		log.Printf("Error fetching relations for suggestions: %v", relErr)
	}

	relationMap := make(map[int]map[string][]string)
	for _, rel := range relations {
		relationMap[rel.ID] = rel.DatesLocations
	}

	suggestions := buildSuggestions(artists, relationMap, query, 15)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(suggestions); err != nil {
		log.Printf("Error encoding suggestions: %v", err)
	}
}

// buildSuggestions collects matches by category first, then interleaves categories so broad
// queries still include artist/band names and locations within the capped dropdown.
func buildSuggestions(artists []models.Artist, relationMap map[int]map[string][]string, query string, limit int) []Suggestion {
	categoryMatches := map[string][]Suggestion{
		"artist/band":      {},
		"location":         {},
		"member":           {},
		"creation date":    {},
		"first album date": {},
	}
	categoryOrder := []string{"artist/band", "location", "member", "creation date", "first album date"}
	seen := make(map[string]bool)

	addToCategory := func(text, category string) {
		key := text + "|" + category
		if seen[key] {
			return
		}
		seen[key] = true
		categoryMatches[category] = append(categoryMatches[category], Suggestion{Text: text, Category: category})
	}

	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), query) {
			addToCategory(artist.Name, "artist/band")
		}

		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				addToCategory(member, "member")
			}
		}

		yearStr := strconv.Itoa(artist.CreationDate)
		if strings.Contains(yearStr, query) {
			addToCategory(yearStr, "creation date")
		}

		if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
			addToCategory(artist.FirstAlbum, "first album date")
		}

		if locs, ok := relationMap[artist.ID]; ok {
			for location := range locs {
				if strings.Contains(strings.ToLower(location), query) {
					city, country := services.FormatLocationName(location)
					formatted := city
					if country != "" {
						formatted = city + ", " + country
					}
					addToCategory(formatted, "location")
				}
			}
		}
	}

	suggestions := make([]Suggestion, 0, limit)
	for len(suggestions) < limit {
		added := false
		for _, category := range categoryOrder {
			if len(categoryMatches[category]) == 0 {
				continue
			}
			suggestions = append(suggestions, categoryMatches[category][0])
			categoryMatches[category] = categoryMatches[category][1:]
			added = true
			if len(suggestions) >= limit {
				break
			}
		}
		if !added {
			break
		}
	}

	return suggestions
}
