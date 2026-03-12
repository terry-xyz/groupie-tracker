package handlers

import (
	"groupie-tracker/services"
	"html/template"
	"log"
	"net/http"
)

// HomeHandler serves the main page by fetching all artists from the API and rendering the home template.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { // Exact match only; prevents "/" from being a catch-all for unregistered routes
		ErrorHandler(w, http.StatusNotFound, "Page not found")
		return
	}

	artists, err := services.FetchArtists() // All artist data comes from the external Groupie Trackers API
	if err != nil {
		log.Printf("Error fetching artists: %v", err)
		ErrorHandler(w, http.StatusInternalServerError, "Unable to load artists. Please try again later.")
		return
	}

	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		ErrorHandler(w, http.StatusInternalServerError, "Unable to load page")
		return
	}

	tmpl.Execute(w, artists) // Pass the full artist slice as template data so home.html can render the grid
}
