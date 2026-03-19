package handlers

import (
	"encoding/json"
	"groupie-tracker/services"
	"log"
	"net/http"
	"strconv"
)

// ArtistGeoHandler returns geocoded locations for a single artist as JSON.
// Called asynchronously by the frontend so the artist page loads instantly.
func ArtistGeoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	relations, err := services.GetAllRelations()
	if err != nil {
		log.Printf("Error fetching relations for geo: %v", err)
		http.Error(w, `{"error":"failed to fetch relations"}`, http.StatusInternalServerError)
		return
	}

	// Find the relation for this artist
	var datesLocations map[string][]string
	for _, rel := range relations {
		if rel.ID == id {
			datesLocations = rel.DatesLocations
			break
		}
	}

	if datesLocations == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]services.GeoLocation{})
		return
	}

	locations := services.GetGeoLocations(id, datesLocations)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}
