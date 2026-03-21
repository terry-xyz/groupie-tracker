package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

// Per-artist geocode cache to avoid re-geocoding on repeated page visits
var (
	artistGeoMu    sync.RWMutex
	artistGeoCache = make(map[int][]GeoLocation)
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

	latF, err := strconv.ParseFloat(results[0].Lat, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid lat %q from Nominatim: %v", results[0].Lat, err)
	}
	lngF, err := strconv.ParseFloat(results[0].Lon, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid lng %q from Nominatim: %v", results[0].Lon, err)
	}

	// Cache the result and persist to disk
	geocodeCacheMu.Lock()
	geocodeCache[query] = [2]float64{latF, lngF}
	geocodeCacheMu.Unlock()

	saveGeocodeCache("data/geocode_cache.json")

	return latF, lngF, nil
}

// GetGeoLocations returns cached geocode results for an artist, geocoding on the first call.
func GetGeoLocations(artistID int, datesLocations map[string][]string) []GeoLocation {
	artistGeoMu.RLock()
	if locs, ok := artistGeoCache[artistID]; ok {
		artistGeoMu.RUnlock()
		return locs
	}
	artistGeoMu.RUnlock()

	locs := GeocodeLocations(datesLocations)

	artistGeoMu.Lock()
	artistGeoCache[artistID] = locs
	artistGeoMu.Unlock()

	return locs
}

// LoadGeocodeCache reads a previously saved geocode cache from disk into memory.
func LoadGeocodeCache(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No cache file yet — not an error
		}
		return fmt.Errorf("failed to read geocode cache: %v", err)
	}

	var cache map[string][2]float64
	if err := json.Unmarshal(data, &cache); err != nil {
		return fmt.Errorf("failed to parse geocode cache: %v", err)
	}

	geocodeCacheMu.Lock()
	for k, v := range cache {
		geocodeCache[k] = v
	}
	geocodeCacheMu.Unlock()

	log.Printf("Loaded %d geocode entries from %s", len(cache), path)
	return nil
}

// saveGeocodeCache writes the current in-memory geocode cache to disk.
func saveGeocodeCache(path string) {
	geocodeCacheMu.RLock()
	snapshot := make(map[string][2]float64, len(geocodeCache))
	for k, v := range geocodeCache {
		snapshot[k] = v
	}
	geocodeCacheMu.RUnlock()

	if err := os.MkdirAll("data", 0755); err != nil {
		log.Printf("Failed to create data dir: %v", err)
		return
	}

	data, err := json.Marshal(snapshot)
	if err != nil {
		log.Printf("Failed to marshal geocode cache: %v", err)
		return
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		log.Printf("Failed to write geocode cache: %v", err)
	}
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
