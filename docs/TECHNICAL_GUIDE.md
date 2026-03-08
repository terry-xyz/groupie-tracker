# Technical Implementation Guide

## Understanding the API

### API Structure

**Base URL**: `https://groupietrackers.herokuapp.com/api`

#### 1. Artists Endpoint
```
GET /api/artists
```
Returns array of artists with:
- id, name, image, members[], creationDate, firstAlbum
- URLs to related locations, dates, relations

#### 2. Locations Endpoint
```
GET /api/locations
```
Returns concert locations for all artists

#### 3. Dates Endpoint
```
GET /api/dates
```
Returns concert dates for all artists

#### 4. Relation Endpoint
```
GET /api/relation
```
Links locations with dates (which concert happened where and when)

---

## Go Implementation Examples

### 1. Data Models (models/artist.go)

```go
package models

type Artist struct {
    ID           int      `json:"id"`
    Name         string   `json:"name"`
    Image        string   `json:"image"`
    Members      []string `json:"members"`
    CreationDate int      `json:"creationDate"`
    FirstAlbum   string   `json:"firstAlbum"`
    Locations    string   `json:"locations"`
    ConcertDates string   `json:"concertDates"`
    Relations    string   `json:"relations"`
}

type Location struct {
    ID        int      `json:"id"`
    Locations []string `json:"locations"`
    Dates     string   `json:"dates"`
}

type Relation struct {
    ID             int                 `json:"id"`
    DatesLocations map[string][]string `json:"datesLocations"`
}
```

### 2. API Service (services/api.go)

```go
package services

import (
    "encoding/json"
    "net/http"
    "time"
)

func FetchArtists() ([]Artist, error) {
    client := &http.Client{Timeout: 10 * time.Second}
    
    resp, err := client.Get("https://groupietrackers.herokuapp.com/api/artists")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var artists []Artist
    if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
        return nil, err
    }
    
    return artists, nil
}
```

### 3. HTTP Handlers (handlers/home.go)

```go
package handlers

import (
    "html/template"
    "net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    
    artists, err := services.FetchArtists()
    if err != nil {
        http.Error(w, "Unable to fetch data", http.StatusInternalServerError)
        return
    }
    
    tmpl := template.Must(template.ParseFiles("templates/home.html"))
    tmpl.Execute(w, artists)
}
```

### 4. Main Server (main.go)

```go
package main

import (
    "log"
    "net/http"
)

func main() {
    // Serve static files
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    
    // Routes
    http.HandleFunc("/", handlers.HomeHandler)
    http.HandleFunc("/artist/", handlers.ArtistHandler)
    
    // Start server
    log.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
```

---

## HTML Template Example

### templates/home.html

```html
<!DOCTYPE html>
<html>
<head>
    <title>Groupie Tracker</title>
    <link rel="stylesheet" href="/static/style.css">
</head>
<body>
    <h1>Artists</h1>
    <div class="artist-grid">
        {{range .}}
        <div class="artist-card">
            <img src="{{.Image}}" alt="{{.Name}}">
            <h2>{{.Name}}</h2>
            <p>Formed: {{.CreationDate}}</p>
            <a href="/artist/{{.ID}}">View Details</a>
        </div>
        {{end}}
    </div>
</body>
</html>
```

---

## Client-Server Feature Example

### JavaScript (Search Feature)

```javascript
// static/script.js
async function searchArtists(query) {
    const response = await fetch(`/api/search?q=${query}`);
    const results = await response.json();
    displayResults(results);
}

document.getElementById('search').addEventListener('input', (e) => {
    searchArtists(e.target.value);
});
```

### Go Handler (Search)

```go
func SearchHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    
    artists, _ := services.FetchArtists()
    var results []Artist
    
    for _, artist := range artists {
        if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
            results = append(results, artist)
        }
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}
```

---

## Error Handling Best Practices

### 1. Check Every Error
```go
resp, err := http.Get(url)
if err != nil {
    // Handle error
    return
}
defer resp.Body.Close()
```

### 2. Validate Status Codes
```go
if resp.StatusCode != http.StatusOK {
    return fmt.Errorf("API returned status: %d", resp.StatusCode)
}
```

### 3. User-Friendly Error Pages
```go
func ErrorHandler(w http.ResponseWriter, status int) {
    w.WriteHeader(status)
    tmpl := template.Must(template.ParseFiles("templates/error.html"))
    tmpl.Execute(w, map[string]int{"Status": status})
}
```

---

## Testing Example

### tests/api_test.go

```go
package tests

import "testing"

func TestFetchArtists(t *testing.T) {
    artists, err := services.FetchArtists()
    
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if len(artists) == 0 {
        t.Error("Expected artists, got empty array")
    }
}
```

Run tests:
```bash
go test ./tests/...
```

---

## Running the Project

```bash
# Navigate to project directory
cd groupie-tracker

# Run the server
go run main.go

# Visit in browser
http://localhost:8080
```
