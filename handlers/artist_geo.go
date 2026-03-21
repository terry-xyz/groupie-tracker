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
	idStr := r.URL.Query().Get("id") // Artist ID comes from the ?id= query parameter set by the frontend
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid id"}`))
		return
	}

	relations, err := services.GetAllRelations()
	if err != nil {
		log.Printf("Error fetching relations for geo: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"failed to fetch relations"}`))
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
		// Artist not found in relations — return an empty array so the frontend map hides gracefully
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode([]services.GeoLocation{}); err != nil {
			log.Printf("Error encoding empty geo response: %v", err)
		}
		return
	}

	locations := services.GetGeoLocations(id, datesLocations)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(locations); err != nil {
		log.Printf("Error encoding geo locations for artist %d: %v", id, err)
	}
}
