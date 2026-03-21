package handlers

import (
	"encoding/json"
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
	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q"))) // TrimSpace prevents whitespace-only input from returning suggestions
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

	// Build relation map for location lookups
	relationMap := make(map[int]map[string][]string)
	for _, rel := range relations {
		relationMap[rel.ID] = rel.DatesLocations
	}

	// Deduplicate suggestions by "text|category" key
	seen := make(map[string]bool)
	var suggestions []Suggestion

	// addSuggestion appends a suggestion only if the 10-item cap hasn't been reached and the
	// text|category pair hasn't been seen before, so the same value from two artists isn't listed twice.
	addSuggestion := func(text, category string) {
		if len(suggestions) >= 10 {
			return
		}
		key := text + "|" + category // composite key so identical text from different categories stays distinct
		if seen[key] {
			return
		}
		seen[key] = true
		suggestions = append(suggestions, Suggestion{Text: text, Category: category})
	}

	for _, artist := range artists {
		if len(suggestions) >= 10 {
			break
		}

		// Artist/band names
		if strings.Contains(strings.ToLower(artist.Name), query) {
			addSuggestion(artist.Name, "artist/band")
		}

		// Member names
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				addSuggestion(member, "member")
			}
		}

		// Creation date
		yearStr := strconv.Itoa(artist.CreationDate)
		if strings.Contains(yearStr, query) {
			addSuggestion(yearStr, "creation date")
		}

		// First album date
		if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
			addSuggestion(artist.FirstAlbum, "first album date")
		}

		// Locations from relations
		if locs, ok := relationMap[artist.ID]; ok {
			for location := range locs {
				if strings.Contains(strings.ToLower(location), query) {
					addSuggestion(location, "location")
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(suggestions); err != nil {
		log.Printf("Error encoding suggestions: %v", err)
	}
}
