package handlers

import (
	"groupie-tracker/services"
	"html/template"
	"log"
	"net/http"
)

var homeTmpl *template.Template

func init() {
	var err error
	homeTmpl, err = template.ParseFiles("templates/home.html")
	if err != nil {
		log.Fatalf("Error parsing home template: %v", err)
	}
}

// HomeHandler serves the main page by fetching all artists from the API and rendering the home template.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { // Exact match only; prevents "/" from being a catch-all for unregistered routes
		ErrorHandler(w, http.StatusNotFound, "Page not found")
		return
	}

	artists, err := services.GetArtists() // All artist data comes from the external Groupie Trackers API
	if err != nil {
		log.Printf("Error fetching artists: %v", err)
		ErrorHandler(w, http.StatusInternalServerError, "Unable to load artists. Please try again later.")
		return
	}

	homeTmpl.Execute(w, artists) // Pass the full artist slice as template data so home.html can render the grid
}
