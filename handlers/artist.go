package handlers

import (
	"groupie-tracker/models"
	"groupie-tracker/services"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/artist/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, "Invalid artist ID")
		return
	}

	artists, err := services.FetchArtists()
	if err != nil {
		log.Printf("Error fetching artists: %v", err)
		ErrorHandler(w, http.StatusInternalServerError, "Unable to load artist data")
		return
	}

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

	relation, err := services.FetchRelation(id)
	if err != nil {
		log.Printf("Error fetching relation for artist %d: %v", id, err)
	}

	pageData := ArtistPageData{
		Artist:         *foundArtist,
		Relation:       relation,
		TotalConcerts:  calculateTotalConcerts(relation),
		TotalCountries: calculateTotalCountries(relation),
		YearsActive:    time.Now().Year() - foundArtist.CreationDate,
		BandType:       getBandType(len(foundArtist.Members)),
	}

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		ErrorHandler(w, http.StatusInternalServerError, "Unable to load page")
		return
	}

	tmpl.Execute(w, pageData)
}

func calculateTotalConcerts(relation *models.Relation) int {
	if relation == nil {
		return 0
	}
	total := 0
	for _, dates := range relation.DatesLocations {
		total += len(dates)
	}
	return total
}

func calculateTotalCountries(relation *models.Relation) int {
	if relation == nil {
		return 0
	}
	countries := make(map[string]bool)
	for location := range relation.DatesLocations {
		parts := strings.Split(location, "-")
		if len(parts) > 0 {
			country := strings.TrimSpace(parts[len(parts)-1])
			countries[country] = true
		}
	}
	return len(countries)
}

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
