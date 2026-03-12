package handlers

import (
	"encoding/json"
	"groupie-tracker/services"
	"log"
	"net/http"
	"sort"
	"strings"
)

// CountryLocations groups locations by country for the filter UI.
type CountryLocations struct {
	Name      string   `json:"name"`
	Locations []string `json:"locations"`
}

// LocationsResponse is the JSON response for GET /api/locations.
type LocationsResponse struct {
	Countries []CountryLocations `json:"countries"`
}

// LocationsHandler handles GET /api/locations, returning all unique locations grouped by country.
func LocationsHandler(w http.ResponseWriter, r *http.Request) {
	relations, err := services.GetAllRelations()
	if err != nil {
		log.Printf("Error fetching relations for locations: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch locations"})
		return
	}

	// Extract all unique location keys and group by country
	countryMap := make(map[string]map[string]bool) // country -> set of locations
	for _, rel := range relations {
		for location := range rel.DatesLocations {
			parts := strings.Split(location, "-")
			country := location // fallback: use the whole string
			if len(parts) > 1 {
				country = parts[len(parts)-1]
			}
			if countryMap[country] == nil {
				countryMap[country] = make(map[string]bool)
			}
			countryMap[country][location] = true
		}
	}

	// Convert to sorted response
	var countries []CountryLocations
	for name, locSet := range countryMap {
		var locs []string
		for loc := range locSet {
			locs = append(locs, loc)
		}
		sort.Strings(locs)
		countries = append(countries, CountryLocations{Name: name, Locations: locs})
	}
	sort.Slice(countries, func(i, j int) bool {
		return countries[i].Name < countries[j].Name
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LocationsResponse{Countries: countries})
}
