package handlers

import (
	"groupie-tracker/models"
	"groupie-tracker/services"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ArtistPageData struct {
	Artist         models.Artist
	Relation       *models.Relation
	TotalConcerts  int
	TotalCountries int
	YearsActive    int
	BandType       string
}

var (
	artistTmpl     *template.Template
	artistTmplOnce sync.Once
)

func getArtistTmpl() *template.Template {
	artistTmplOnce.Do(func() {
		var err error
		artistTmpl, err = template.ParseFiles("templates/artist.html")
		if err != nil {
			log.Fatalf("Error parsing artist template: %v", err)
		}
	})
	return artistTmpl
}

// ArtistHandler serves the detail page for a single artist, identified by the numeric ID in the URL path (e.g., /artist/3).
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/artist/") // Extract the numeric ID from the URL path segment
	id, err := strconv.Atoi(idStr)                       // Convert path segment to int for artist lookup
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, "Invalid artist ID")
		return
	}

	artists, err := services.GetArtists()
	if err != nil {
		log.Printf("Error fetching artists: %v", err)
		ErrorHandler(w, http.StatusInternalServerError, "Unable to load artist data")
		return
	}

	// Linear search because the API returns all artists without a single-artist endpoint
	var foundArtist *models.Artist
	for _, a := range artists {
		if a.ID == id {
			foundArtist = &a
			break
		}
	}

	if foundArtist == nil {
		ErrorHandler(w, http.StatusNotFound, "Artist not found")
		return
	}

	// Use cached relations and find the one for this artist
	var relation *models.Relation
	relations, err := services.GetAllRelations()
	if err != nil {
		log.Printf("Error fetching relations for artist %d: %v", id, err)
	} else {
		for i := range relations {
			if relations[i].ID == id {
				relation = &relations[i]
				break
			}
		}
	}

	pageData := ArtistPageData{
		Artist:         *foundArtist,
		Relation:       relation,
		TotalConcerts:  calculateTotalConcerts(relation),
		TotalCountries: calculateTotalCountries(relation),
		YearsActive:    time.Now().Year() - foundArtist.CreationDate, // Years active = current year minus band formation year
		BandType:       getBandType(len(foundArtist.Members)),        // Classify group size (solo, duo, trio, etc.)
	}

	getArtistTmpl().Execute(w, pageData)
}

// calculateTotalConcerts returns the total number of concert dates across all locations for an artist.
func calculateTotalConcerts(relation *models.Relation) int {
	if relation == nil {
		return 0
	}
	total := 0
	for _, dates := range relation.DatesLocations {
		total += len(dates) // Each location maps to multiple concert dates, so we sum all date slices
	}
	return total
}

// calculateTotalCountries returns the number of unique countries extracted from the relation's location keys.
func calculateTotalCountries(relation *models.Relation) int {
	if relation == nil {
		return 0
	}
	countries := make(map[string]bool) // Map deduplicates country names so we count unique ones
	for location := range relation.DatesLocations {
		parts := strings.Split(location, "-") // Location format is "city-country", so we split on "-"
		if len(parts) > 0 {
			country := strings.TrimSpace(parts[len(parts)-1]) // Last segment is always the country name
			countries[country] = true
		}
	}
	return len(countries)
}

// getBandType returns a human-readable label for the group size (e.g., "Solo Artist", "Duo", "Trio").
func getBandType(memberCount int) string {
	switch memberCount {
	case 1:
		return "Solo Artist"
	case 2:
		return "Duo"
	case 3:
		return "Trio"
	case 4:
		return "Quartet"
	case 5:
		return "Quintet"
	default:
		return "Band"
	}
}
