package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"unicode"
)

// GeoLocation represents a geocoded concert location with coordinates and metadata.
type GeoLocation struct {
	Lat     float64  `json:"lat"`
	Lng     float64  `json:"lng"`
	City    string   `json:"city"`
	Country string   `json:"country"`
	Dates   []string `json:"dates"`
	RawKey  string   `json:"rawKey"`
}

// Geocode cache to avoid redundant Nominatim lookups
var (
	geocodeCacheMu sync.RWMutex
	geocodeCache   = make(map[string][2]float64)
)

// Rate limiter: Nominatim free tier allows 1 request per second
var rateLimiter = make(chan struct{}, 1)

func init() {
	// Feed the rate limiter at 1 token per 1.1 seconds
	go func() {
		for {
			rateLimiter <- struct{}{}
			time.Sleep(1100 * time.Millisecond)
		}
	}()
}

// TitleCase capitalizes the first letter of each word (avoids deprecated strings.Title).
func TitleCase(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			runes := []rune(w)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}

// FormatLocationName splits a raw API location key like "north_carolina-usa" into city and country.
func FormatLocationName(raw string) (city, country string) {
	// Find the last "-" to split city from country
	lastDash := strings.LastIndex(raw, "-")
	if lastDash == -1 {
		// No dash: treat the whole string as the city
		city = TitleCase(strings.ReplaceAll(raw, "_", " "))
		return city, ""
	}
	rawCity := raw[:lastDash]
	rawCountry := raw[lastDash+1:]

	city = TitleCase(strings.ReplaceAll(rawCity, "_", " "))
	country = TitleCase(strings.ReplaceAll(rawCountry, "_", " "))
	return city, country
}

// nominatimResult represents a single result from the Nominatim search API.
type nominatimResult struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

// geocodeAddress looks up coordinates for a location query using Nominatim.
func geocodeAddress(query string) (lat, lng float64, err error) {
	// Check cache first
	geocodeCacheMu.RLock()
	if coords, ok := geocodeCache[query]; ok {
		geocodeCacheMu.RUnlock()
		return coords[0], coords[1], nil
	}
	geocodeCacheMu.RUnlock()

	// Wait for rate limiter
	<-rateLimiter

	reqURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1",
		url.QueryEscape(query))

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("User-Agent", "GroupieTracker/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("geocode request failed: %v", err)
	}
	defer resp.Body.Close()

	var results []nominatimResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return 0, 0, fmt.Errorf("failed to parse geocode response: %v", err)
	}

	if len(results) == 0 {
		return 0, 0, fmt.Errorf("no geocode results for %q", query)
	}

	var latF, lngF float64
	fmt.Sscanf(results[0].Lat, "%f", &latF)
	fmt.Sscanf(results[0].Lon, "%f", &lngF)

	// Cache the result
	geocodeCacheMu.Lock()
	geocodeCache[query] = [2]float64{latF, lngF}
	geocodeCacheMu.Unlock()

	return latF, lngF, nil
}

// GeocodeLocations geocodes all locations from an artist's DatesLocations map.
func GeocodeLocations(datesLocations map[string][]string) []GeoLocation {
	var locations []GeoLocation

	for rawKey, dates := range datesLocations {
		city, country := FormatLocationName(rawKey)
		query := city
		if country != "" {
			query = city + ", " + country
		}

		lat, lng, err := geocodeAddress(query)
		if err != nil {
			// Skip locations that fail to geocode
			continue
		}

		locations = append(locations, GeoLocation{
			Lat:     lat,
			Lng:     lng,
			City:    city,
			Country: country,
			Dates:   dates,
			RawKey:  rawKey,
		})
	}

	return locations
}
