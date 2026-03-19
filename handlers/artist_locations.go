package handlers

import (
	"encoding/json"
	"groupie-tracker/services"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// ArtistLocationsHandler geocodes concert locations for a single artist and returns them as JSON.
// This is called asynchronously by the browser after the artist page has loaded, so geocoding
// never blocks the initial page render.
func ArtistLocationsHandler(w http.ResponseWriter, r *http.Request) {
	// URL pattern: /api/artist/{id}/locations
	path := strings.TrimPrefix(r.URL.Path, "/api/artist/")
	path = strings.TrimSuffix(path, "/locations")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	relation, err := services.FetchRelation(id)
	if err != nil {
		log.Printf("Error fetching relation for artist %d: %v", id, err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]services.GeoLocation{})
		return
	}

	locations := services.GeocodeLocations(relation.DatesLocations)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}
