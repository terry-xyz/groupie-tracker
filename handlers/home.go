package handlers

import (
	"groupie-tracker/services"
	"html/template"
	"log"
	"net/http"
	"sync"
)

var (
	homeTmpl     *template.Template
	homeTmplOnce sync.Once
)

// getHomeTmpl parses and caches the home template once so repeated requests don't re-parse the file.
func getHomeTmpl() *template.Template {
	homeTmplOnce.Do(func() {
		var err error
		homeTmpl, err = template.ParseFiles("templates/home.html")
		if err != nil {
			log.Fatalf("Error parsing home template: %v", err)
		}
	})
	return homeTmpl
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

	if err := getHomeTmpl().Execute(w, artists); err != nil { // artists slice is the template data; home.html iterates it to render the initial artist grid
		log.Printf("Error rendering home template: %v", err)
	}
}
