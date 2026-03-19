package main

import (
	"groupie-tracker/handlers"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs)) // Strip "/static/" so file paths resolve relative to the "static" directory

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/artist/", handlers.ArtistHandler)
	http.HandleFunc("/api/search", handlers.SearchHandler)
	http.HandleFunc("/api/suggestions", handlers.SuggestionsHandler)
	http.HandleFunc("/api/locations", handlers.LocationsHandler)
	http.HandleFunc("/api/artist/", handlers.ArtistLocationsHandler)

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
