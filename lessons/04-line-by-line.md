# Lesson 04: Line-by-Line Walkthrough

## Reading Go Code: Quick Syntax Guide

```go
package handlers                        // This file belongs to the "handlers" group

import (                                // Load libraries we need
    "encoding/json"                     // For converting data to/from JSON
    "net/http"                          // For web server functionality
)

func Name(param Type) ReturnType {      // Function definition
    // param is input, ReturnType is output
}

func Name(a, b Type) (Type, error) {    // Multiple returns
    // Returns BOTH a result AND an error
}

if condition {                          // If statement
    // runs if condition is true
}

for i := 0; i < 10; i++ {             // Loop 0 to 9
    // i starts at 0, runs while < 10, adds 1 each time
}

for _, item := range slice {           // Loop over a list
    // _ = "ignore the index", item = current element
}

defer cleanup()                        // "Do this LATER, when function ends"

switch value {                         // Multi-way branch (like if/else if/else)
case "a":                              // If value == "a"
case "b":                              // If value == "b"
default:                               // Otherwise
}
```

---

## Section 1: The Entry Point

**Big Picture:** `main.go` is the power switch. It sets up the web server, tells it which handler to use for each URL, and starts listening for requests. It's the shortest file in the project — 8 meaningful lines — but everything starts here.

**File:** `main.go:1-18`

```go
// INPUT:  Nothing (this is where the program starts)
// OUTPUT: A running web server on port 8080

package main  // "main" package = this is an executable program (not a library)

import (
    "groupie-tracker/handlers"  // Our request handlers (from the handlers/ folder)
    "log"                       // For printing messages to the terminal
    "net/http"                  // Go's built-in web server
)

func main() {  // main() is THE function Go runs when the program starts

    // STEP 1: Serve static files (CSS, JS, images)
    // When browser requests /static/style.css, serve the file from static/style.css
    fs := http.FileServer(http.Dir("static"))               // Create a file server pointing to "static/" folder
    http.Handle("/static/", http.StripPrefix("/static/", fs)) // Remove "/static/" prefix before looking up the file

    // STEP 2: Register route handlers
    // "When someone visits this URL, run this function"
    http.HandleFunc("/", handlers.HomeHandler)              // Home page
    http.HandleFunc("/artist/", handlers.ArtistHandler)     // Artist detail pages
    http.HandleFunc("/api/search", handlers.SearchHandler)  // Search API endpoint

    // STEP 3: Start the server
    log.Println("Server starting on http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)  // If server can't start (port busy?), crash with error message
    }
}
```

**Key Insight:** `http.StripPrefix("/static/", fs)` is necessary because without it, a request for `/static/style.css` would look for a file at `static/static/style.css` (doubling the path). `StripPrefix` removes the URL prefix before passing the path to the file server.

---

## Section 2: Data Models

**Big Picture:** `models/artist.go` defines the shape of our data — like a cookie cutter that determines what an artist "looks like" in code. Every other file references these definitions.

**File:** `models/artist.go:1-16`

```go
// INPUT:  Nothing (these are just definitions, not functions)
// OUTPUT: Reusable data types for the entire project

package models  // "models" = this package defines data shapes

// Artist is the main data type. Everything revolves around this.
type Artist struct {
    ID           int      `json:"id"`           // Unique number from the API (1, 2, 3...)
    Name         string   `json:"name"`         // "Queen", "SOJA", "Mamonas Assassinas"
    Image        string   `json:"image"`        // Full URL to artist photo
    Members      []string `json:"members"`      // ["Freddie Mercury", "Brian May", ...]
    CreationDate int      `json:"creationDate"` // Year formed: 1970 (NOT a date, just a year)
    FirstAlbum   string   `json:"firstAlbum"`   // "14-12-1973" (DD-MM-YYYY format as STRING)
}

// Relation connects an artist to their concert history.
// Think of it as an artist's "tour schedule."
type Relation struct {
    ID             int                 `json:"id"`             // Must match an Artist's ID
    DatesLocations map[string][]string `json:"datesLocations"` // Location → list of concert dates
    // Example:
    // "london-england": ["23-08-2019", "15-03-2020"]
    // "new_york-usa":   ["01-12-2019"]
}
```

**Key Insight:** `FirstAlbum` is a `string`, not a date type. The API sends dates as `"DD-MM-YYYY"` text. This is a common real-world pattern — APIs often send dates as strings, and it's the receiver's job to parse them if needed.

---

## Section 3: API Service

**Big Picture:** `services/api.go` is the "delivery driver." It makes HTTP calls to the external API, waits for the response, converts the JSON into Go structs, and brings the data back. Every handler depends on this file.

**File:** `services/api.go:1-58`

```go
package services

import (
    "encoding/json"
    "fmt"
    "groupie-tracker/models"
    "net/http"
    "time"
)

// Base URL for all API calls — like a phone number's area code
const baseURL = "https://groupietrackers.herokuapp.com/api"

// FetchArtists gets ALL artists from the external API.
// INPUT:  Nothing
// OUTPUT: A list of artists, or an error explaining what went wrong
func FetchArtists() ([]models.Artist, error) {
    // STEP 1: Create an HTTP client with a safety timeout
    // Why 10 seconds? Long enough for slow networks, short enough to not freeze the app
    client := &http.Client{Timeout: 10 * time.Second}

    // STEP 2: Make the actual HTTP request
    // This is like typing a URL into your browser, but in code
    resp, err := client.Get(baseURL + "/artists")
    if err != nil {
        // Network error: API is down, no internet, DNS failure, etc.
        return nil, fmt.Errorf("failed to connect to API: %v", err)
    }
    defer resp.Body.Close()  // "When this function ends, close the connection"
                              // Like hanging up the phone when the call is over

    // STEP 3: Check the HTTP status code
    // 200 = "OK, here's your data"
    // 404 = "Not found", 500 = "Server error", etc.
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
    }

    // STEP 4: Decode the JSON response into Go structs
    // resp.Body is raw text: [{"id":1,"name":"Queen",...}]
    // After decoding: []Artist{{ID:1, Name:"Queen",...}}
    var artists []models.Artist
    if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
        return nil, fmt.Errorf("failed to parse artist data: %v", err)
    }

    return artists, nil  // nil error = "everything went fine"
}
```

```go
// FetchRelation gets concert data for ONE specific artist.
// INPUT:  Artist ID (integer)
// OUTPUT: Pointer to Relation data, or an error
func FetchRelation(id int) (*models.Relation, error) {
    // Same pattern as FetchArtists, but for the relation endpoint
    client := &http.Client{Timeout: 10 * time.Second}

    resp, err := client.Get(baseURL + "/relation")
    if err != nil {
        return nil, fmt.Errorf("failed to connect to API: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
    }

    // The relation endpoint returns ALL relations wrapped in an "index" array
    // Structure: {"index": [{"id":1,"datesLocations":{...}}, {"id":2,...}]}
    var result struct {
        Index []models.Relation `json:"index"`
    }
    // ↑ Anonymous struct: a one-time-use struct defined right here
    //   We only need it to unwrap the "index" layer

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to parse relation data: %v", err)
    }

    // STEP 5: Find the relation matching our artist ID
    for _, rel := range result.Index {
        if rel.ID == id {
            return &rel, nil  // Found it! Return a pointer to it
            // "&rel" means "the memory address of rel"
            // We return a pointer (*Relation) instead of a copy for efficiency
        }
    }

    // If we get here, no relation was found for this artist
    return nil, fmt.Errorf("relation not found for artist %d", id)
}
```

**Key Insight:** `FetchRelation` fetches ALL relations and then searches for the right one. This is because the API doesn't have a "get relation by ID" endpoint — it only has "get all relations." This means every artist detail page downloads ALL relations just to find one. It works fine for 52 artists but would be inefficient for thousands.

---

## Section 4: Home Handler

**Big Picture:** The HomeHandler is the simplest handler — it fetches all artists and displays them in a grid. It's a straightforward fetch-and-render with validation and error handling.

**File:** `handlers/home.go:1-31`

```go
package handlers

import (
    "groupie-tracker/services"
    "html/template"
    "log"
    "net/http"
)

// INPUT:  HTTP request for "/"
// OUTPUT: Rendered home page with all artists, or error page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    // GUARD: Make sure they're visiting exactly "/"
    // Without this, "/anything" would also show the home page
    // because Go's "/" pattern matches ALL paths as a prefix
    if r.URL.Path != "/" {
        ErrorHandler(w, http.StatusNotFound)
        return
    }

    // STEP 1: Get all artists from external API
    artists, err := services.FetchArtists()
    if err != nil {
        log.Printf("Error fetching artists: %v", err)
        ErrorHandler(w, http.StatusInternalServerError, "Unable to load artists. Please try again later.")
        return
    }

    // STEP 2: Load and parse the HTML template
    tmpl, err := template.ParseFiles("templates/home.html")
    if err != nil {
        log.Printf("Error parsing template: %v", err)
        ErrorHandler(w, http.StatusInternalServerError)
        return
    }

    // STEP 3: Execute template — fills in {{.Name}}, {{.Image}}, etc.
    // "w" is the response writer — template output goes directly to the browser
    // "artists" is the data — the template's "." (dot) becomes this slice
    tmpl.Execute(w, artists)
}
```

**Key Insight:** `tmpl.Execute(w, artists)` — the second argument becomes `{{.}}` (dot) inside the template. Since `artists` is a slice, `{{range .}}` loops over each artist.

---

## Section 5: Artist Handler

**Big Picture:** The most complex handler. It extracts an ID from the URL, fetches the artist and their concert data, calculates statistics, and renders the profile page. Think of it as a factory that assembles a complete artist profile from multiple data sources.

**File:** `handlers/artist.go:1-117`

```go
// INPUT:  HTTP request for "/artist/{id}"
// OUTPUT: Rendered artist profile page, or error page

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
    // STEP 1: Extract the artist ID from the URL
    // URL: "/artist/3" → TrimPrefix removes "/artist/" → "3"
    idStr := strings.TrimPrefix(r.URL.Path, "/artist/")

    // STEP 2: Convert string "3" to integer 3
    // If someone visits "/artist/banana", Atoi fails → err is not nil
    id, err := strconv.Atoi(idStr)
    if err != nil {
        ErrorHandler(w, http.StatusBadRequest, "Invalid artist ID")
        return  // "banana" is not a number → 400 Bad Request
    }

    // STEP 3: Fetch ALL artists (there's no "get one artist" endpoint)
    artists, err := services.FetchArtists()
    if err != nil {
        log.Printf("Error fetching artists: %v", err)
        ErrorHandler(w, http.StatusInternalServerError, "Unable to load artist data")
        return
    }

    // STEP 4: Search through artists to find the one with our ID
    var artist *models.Artist  // Pointer — nil until we find a match
    for i, a := range artists {
        if a.ID == id {
            artist = &artists[i]  // Found it! Point to it
            break                  // Stop searching
        }
    }

    // STEP 5: If no artist matched, it doesn't exist
    if artist == nil {
        ErrorHandler(w, http.StatusNotFound, "Artist not found")
        return
    }

    // STEP 6: Fetch concert relation data
    relation, err := services.FetchRelation(id)
    if err != nil {
        log.Printf("Error fetching relation: %v", err)
        ErrorHandler(w, http.StatusInternalServerError, "Unable to load concert data")
        return
    }

    // STEP 7: Calculate statistics for the profile page
    // This builds the "stats bar" — concerts, countries, years active, band type
    currentYear := time.Now().Year()
    data := ArtistPageData{
        Artist:         *artist,                                 // The artist's basic info
        Relation:       relation,                                // Concert data
        TotalConcerts:  calculateTotalConcerts(relation),        // Sum of all concerts
        TotalCountries: calculateTotalCountries(relation),       // Unique countries
        YearsActive:    currentYear - artist.CreationDate,       // 2026 - 1970 = 56
        BandType:       getBandType(len(artist.Members)),        // "Quartet" for 4 members
    }

    // STEP 8: Parse and render the template
    tmpl, err := template.ParseFiles("templates/artist.html")
    if err != nil {
        log.Printf("Error parsing template: %v", err)
        ErrorHandler(w, http.StatusInternalServerError)
        return
    }

    tmpl.Execute(w, data)
}
```

### Helper Functions

```go
// calculateTotalConcerts counts every concert date across all locations.
// INPUT:  Relation with DatesLocations map
// OUTPUT: Total number of concerts (integer)
//
// Example: {"london": ["date1", "date2"], "paris": ["date3"]} → 3
func calculateTotalConcerts(relation *models.Relation) int {
    total := 0
    for _, dates := range relation.DatesLocations {
        // Each location has a list of dates
        // len(dates) = number of concerts at that location
        total += len(dates)
    }
    return total
}

// calculateTotalCountries counts UNIQUE countries from location strings.
// INPUT:  Relation with location keys like "london-england", "paris-france"
// OUTPUT: Number of unique countries
//
// The location format is "city-country" or "city_name-country"
// We extract the part after the LAST "-" as the country
func calculateTotalCountries(relation *models.Relation) int {
    countries := make(map[string]bool)  // map acts as a "set" (no duplicates)

    for location := range relation.DatesLocations {
        // "new_york-usa" → split by "-" → ["new_york", "usa"] → last = "usa"
        parts := strings.Split(location, "-")
        country := parts[len(parts)-1]        // Last element = country
        countries[strings.ToLower(country)] = true  // Add to set (lowercase for consistency)
    }

    return len(countries)  // Number of unique keys = number of unique countries
}

// getBandType classifies a band by its member count.
// INPUT:  Number of members (integer)
// OUTPUT: Human-readable label
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
        return "Band"  // 6+ members
    }
}
```

**Key Insight:** `calculateTotalCountries` uses `map[string]bool` as a **set** — a collection that automatically prevents duplicates. Adding `"usa": true` twice doesn't create two entries. This is the standard Go idiom for "unique items" since Go doesn't have a built-in Set type.

---

## Section 6: Search Handler

**Big Picture:** SearchHandler is the API endpoint that powers the search feature. Unlike other handlers that return HTML, this one returns JSON — because it's called by JavaScript in the background, not by a browser page load.

**File:** `handlers/search.go:1-100`

```go
// INPUT:  HTTP request with query params: ?q=queen&minYear=1970&maxYear=2000&sort=name
// OUTPUT: JSON array of matching artists

func SearchHandler(w http.ResponseWriter, r *http.Request) {
    // STEP 1: Extract search parameters from the URL
    query := r.URL.Query().Get("q")          // Search text (e.g., "queen")
    minYear := r.URL.Query().Get("minYear")  // Minimum formation year
    maxYear := r.URL.Query().Get("maxYear")  // Maximum formation year
    sortBy := r.URL.Query().Get("sort")      // Sort order

    // STEP 2: Fetch ALL artists (we filter in memory)
    artists, err := services.FetchArtists()
    if err != nil {
        log.Printf("Error fetching artists: %v", err)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch artists"})
        return
    }

    // STEP 3: FILTER — keep only artists that match ALL criteria
    var filtered []models.Artist
    for _, artist := range artists {
        if matchesSearch(artist, query) && matchesYearFilter(artist, minYear, maxYear) {
            filtered = append(filtered, artist)
            // append() adds to the list — like pushing onto a stack
        }
    }

    // STEP 4: SORT — order the filtered results
    sortArtists(filtered, sortBy)

    // STEP 5: RESPOND — send filtered results as JSON
    w.Header().Set("Content-Type", "application/json")  // Tell browser "this is JSON, not HTML"
    json.NewEncoder(w).Encode(filtered)                  // Convert structs → JSON text → send
}
```

### Filter Functions

```go
// matchesSearch checks if an artist matches the search query.
// INPUT:  An artist, and the search query string
// OUTPUT: true if the artist matches, false if not
//
// Case-insensitive: "queen" matches "Queen"
// Searches in: artist name AND member names
func matchesSearch(artist models.Artist, query string) bool {
    if query == "" {
        return true  // No search query = match everything
    }

    q := strings.ToLower(query)  // Normalize to lowercase for comparison

    // Check if artist name contains the query
    if strings.Contains(strings.ToLower(artist.Name), q) {
        return true  // "Queen" contains "que" → match!
    }

    // Check if ANY member name contains the query
    for _, member := range artist.Members {
        if strings.Contains(strings.ToLower(member), q) {
            return true  // "Freddie Mercury" contains "fred" → match!
        }
    }

    return false  // No match found anywhere
}
```

```go
// matchesYearFilter checks if an artist's creation year falls within a range.
// INPUT:  An artist, min year string, max year string (may be empty)
// OUTPUT: true if within range (or range not specified)
func matchesYearFilter(artist models.Artist, minYear, maxYear string) bool {
    // If minYear is set and valid, check lower bound
    if minYear != "" {
        min, err := strconv.Atoi(minYear)   // "1970" → 1970
        if err == nil && artist.CreationDate < min {
            return false  // Artist formed before the minimum year → exclude
        }
    }

    // If maxYear is set and valid, check upper bound
    if maxYear != "" {
        max, err := strconv.Atoi(maxYear)
        if err == nil && artist.CreationDate > max {
            return false  // Artist formed after the maximum year → exclude
        }
    }

    return true  // Within range (or no range specified)
}
```

```go
// sortArtists sorts the artist slice in-place.
// INPUT:  Slice of artists, sort criteria string
// OUTPUT: Nothing (modifies the slice directly)
func sortArtists(artists []models.Artist, sortBy string) {
    switch sortBy {
    case "name":
        // Alphabetical by name (A → Z)
        sort.Slice(artists, func(i, j int) bool {
            return strings.ToLower(artists[i].Name) < strings.ToLower(artists[j].Name)
        })
    case "newest":
        // By creation date, newest first (2020 before 1970)
        sort.Slice(artists, func(i, j int) bool {
            return artists[i].CreationDate > artists[j].CreationDate
        })
    case "oldest":
        // By creation date, oldest first (1970 before 2020)
        sort.Slice(artists, func(i, j int) bool {
            return artists[i].CreationDate < artists[j].CreationDate
        })
    default:
        // By ID (the API's original order)
        sort.Slice(artists, func(i, j int) bool {
            return artists[i].ID < artists[j].ID
        })
    }
}
```

**Key Insight:** `sort.Slice` takes a "less" function — it answers "should item `i` come before item `j`?" Go calls this function repeatedly to sort the entire slice. The `func(i, j int) bool` is an **anonymous function** (a function without a name, defined inline).

---

## Section 7: Error Handler

**Big Picture:** ErrorHandler is the "apology department." When anything goes wrong, every handler delegates to this function. It renders a consistent error page with the right HTTP status code.

**File:** `handlers/error.go:1-27`

```go
// INPUT:  Response writer, HTTP status code, optional custom message
// OUTPUT: Rendered error page

func ErrorHandler(w http.ResponseWriter, status int, customMsg ...string) {
    // STEP 1: Set the HTTP status code (404, 500, etc.)
    // This tells the browser "something went wrong" (not just displaying error text)
    w.WriteHeader(status)

    // STEP 2: Choose the error message
    message := http.StatusText(status)  // Default: "Not Found", "Internal Server Error"
    if len(customMsg) > 0 {
        message = customMsg[0]  // Override with custom message if provided
    }

    // STEP 3: Parse and render the error template
    tmpl, err := template.ParseFiles("templates/error.html")
    if err != nil {
        // If even the error template is broken, fall back to plain text
        // This is the "error handler's error handler" — the last resort
        http.Error(w, message, status)
        return
    }

    // STEP 4: Execute template with status code and message
    tmpl.Execute(w, map[string]interface{}{
        "Status":  status,    // e.g., 404
        "Message": message,   // e.g., "Artist not found"
    })
}
```

**Key Insight:** `map[string]interface{}` is a map where keys are strings and values can be **anything** (number, string, struct, etc.). It's Go's equivalent of a "generic dictionary." Here, `"Status"` holds an `int` and `"Message"` holds a `string` in the same map.

---

## Section 8: JavaScript (Client-Side)

**Big Picture:** `script.js` handles everything that happens in the browser after the page loads — theme toggling, search, filtering, and dynamically rebuilding the artist grid without page reloads.

**File:** `static/script.js:1-124`

```javascript
// ==========================================
// THEME MANAGEMENT
// ==========================================

// When the page loads, check if the user previously chose a theme
document.addEventListener('DOMContentLoaded', () => {
    const savedTheme = localStorage.getItem('theme');
    // localStorage = browser's persistent storage (survives page refresh)
    if (savedTheme === 'light') {
        document.body.classList.add('light-theme');
        // Adds a CSS class that swaps all the colors
    }
});

// Toggle between dark and light themes
function toggleTheme() {
    document.body.classList.toggle('light-theme');
    // .toggle() = if class exists, remove it; if not, add it

    const isLight = document.body.classList.contains('light-theme');
    localStorage.setItem('theme', isLight ? 'light' : 'dark');
    // Save preference so it persists across page loads
}

// ==========================================
// SEARCH & FILTER
// ==========================================

async function applyFilters() {
    // STEP 1: Collect all filter values from the form
    const query = document.getElementById('search').value;
    const minYear = document.getElementById('min-year').value;
    const maxYear = document.getElementById('max-year').value;
    const sort = document.getElementById('sort-select').value;

    // STEP 2: Build the search URL
    const url = `/api/search?q=${encodeURIComponent(query)}&minYear=${minYear}&maxYear=${maxYear}&sort=${sort}`;
    // encodeURIComponent: makes special characters URL-safe
    // "Queen & Kings" → "Queen%20%26%20Kings"

    // STEP 3: Show loading spinner
    showLoading();

    try {
        // STEP 4: Send the request (async — doesn't freeze the page)
        const response = await fetch(url);
        // "await" = pause HERE until the server responds
        // But the rest of the page stays responsive (not frozen)

        if (!response.ok) {
            throw new Error('Search failed');
            // "throw" = stop and jump to the "catch" block below
        }

        // STEP 5: Parse the JSON response
        const data = await response.json();
        // Server sent: [{"id":1,"name":"Queen",...}]
        // After .json(): [{id:1, name:"Queen",...}]  (JavaScript objects)

        // STEP 6: Rebuild the page with results
        displayArtists(data);
    } catch (error) {
        // STEP 7: If anything failed, show error message
        showError('Failed to search artists. Please try again.');
    }
}

function displayArtists(artists) {
    const grid = document.getElementById('artists-grid');
    grid.innerHTML = '';  // Wipe the current grid clean

    if (artists.length === 0) {
        // No results found — show a "nothing here" message
        grid.innerHTML = '<p class="no-results">No artists found</p>';
        return;
    }

    // For each artist, create an HTML card and add it to the grid
    artists.forEach(artist => {
        const card = document.createElement('div');
        card.className = 'artist-card';
        card.innerHTML = `
            <img src="${artist.image}" alt="${artist.name}">
            <h2>${artist.name}</h2>
            <p>${artist.members.join(', ')}</p>
            <a href="/artist/${artist.id}">View Details</a>
        `;
        // Template literal (backticks): allows ${variable} insertion
        // artist.members.join(', '): ["A", "B", "C"] → "A, B, C"

        grid.appendChild(card);  // Add the card to the page
    });
}

// ==========================================
// UI HELPERS
// ==========================================

function showLoading() {
    const grid = document.getElementById('artists-grid');
    grid.innerHTML = '<div class="loading-spinner"></div>';
    // Replaces grid content with a CSS spinner animation
}

function showError(message) {
    const grid = document.getElementById('artists-grid');
    grid.innerHTML = `<p class="error-message">${message}</p>`;
}

function resetFilters() {
    // Clear all input fields
    document.getElementById('search').value = '';
    document.getElementById('min-year').value = '';
    document.getElementById('max-year').value = '';
    document.getElementById('sort-select').value = 'default';

    // Reload the page to show all artists again
    window.location.reload();
}

// ==========================================
// EVENT LISTENERS
// ==========================================

// When user presses Enter in the search box, trigger search
document.getElementById('search')?.addEventListener('keyup', (e) => {
    if (e.key === 'Enter') applyFilters();
    // ?. = "optional chaining" — if search element doesn't exist, don't crash
});
```

**Key Insight:** The `async/await` pattern makes asynchronous code (network requests) read like synchronous code. Without it, you'd need nested callbacks or `.then()` chains, which are harder to follow. `await fetch(url)` literally means "wait here until the server responds, then continue."

---

## Summary: The Code Tour

| File | Lines | Role | Complexity |
|------|-------|------|------------|
| `main.go` | ~18 | Wire everything together | Simple |
| `models/artist.go` | ~16 | Data definitions | Simple |
| `services/api.go` | ~58 | Fetch from external API | Medium |
| `handlers/home.go` | ~31 | Render home page | Simple |
| `handlers/artist.go` | ~117 | Render artist profile | Complex |
| `handlers/search.go` | ~100 | Search and filter API | Medium |
| `handlers/error.go` | ~27 | Render error pages | Simple |
| `static/script.js` | ~124 | Browser interactivity | Medium |

---

## What's Next?

Now that you understand every line, [Lesson 05](05-exercises.md) gives you hands-on practice — from finding files to adding new features.
