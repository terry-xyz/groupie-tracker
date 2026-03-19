package main

import (
	"groupie-tracker/handlers"
	"groupie-tracker/services"
	"log"
	"net/http"
	"sync"
)

func main() {
	if err := services.LoadGeocodeCache("data/geocode_cache.json"); err != nil {
		log.Printf("Warning: could not load geocode cache: %v", err)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs)) // Strip "/static/" so file paths resolve relative to the "static" directory

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/artist/", handlers.ArtistHandler)
	http.HandleFunc("/api/search", handlers.SearchHandler)
	http.HandleFunc("/api/suggestions", handlers.SuggestionsHandler)
	http.HandleFunc("/api/locations", handlers.LocationsHandler)

	// Pre-warm caches in the background so the server starts immediately
	go func() {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() { defer wg.Done(); services.GetArtists() }()
		go func() { defer wg.Done(); services.GetAllRelations() }()
		wg.Wait()

		// Geocode all artists in background (respects rate limiter)
		relations, err := services.GetAllRelations()
		if err != nil {
			log.Printf("Pre-warm: failed to fetch relations: %v", err)
			return
		}
		for _, rel := range relations {
			services.GetGeoLocations(rel.ID, rel.DatesLocations)
		}
		log.Println("Pre-warm: geocoding complete")
	}()

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
