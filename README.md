# Groupie Tracker

A web application built with Go that displays information about music bands and artists by fetching data from an external API. Features include artist profiles, concert information, search functionality, and filtering options.

![Go Version](https://img.shields.io/badge/Go-1.24-blue)
![License](https://img.shields.io/badge/license-MIT-green)

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Technologies](#technologies)
- [Testing](#testing)
- [Screenshots](#screenshots)

## Features

### Core Functionality
- **Artist Listing** - Browse all artists with their basic information
- **Artist Profiles** - Detailed view of each artist including:
  - Band members
  - Formation year and years active
  - First album release date
  - Band type classification (Solo, Duo, Trio, etc.)
  - Concert history with locations and dates
  - Statistics (total concerts, countries visited)

### Interactive Features
- **Search** - Search artists by name, member, location, creation date, or first album date
- **Autocomplete Suggestions** - Typing suggestions with category labels (artist/band, member, location, creation date, first album date)
- **Filter by Creation Year** - Filter artists by their formation year (min/max range)
- **Filter by First Album Year** - Filter artists by first album release year (min/max range)
- **Filter by Member Count** - Checkbox filters for 1-7 members or 8+
- **Filter by Location** - Searchable, country-grouped location checkboxes
- **Sorting** - Sort artists by:
  - Name (A-Z)
  - Newest first (most recent formation)
  - Oldest first (earliest formation)
  - Default (by ID)
- **Geolocalization** - Interactive world tour map on artist pages with geocoded concert markers, popups with dates, and route polylines (powered by Leaflet + OpenStreetMap)
- **Result Count** - Displays number of matching artists after filtering

### User Experience
- **Responsive Design** - Works on desktop, tablet, and mobile devices
- **Loading Indicators** - Visual feedback during data fetching
- **Error Handling** - User-friendly error messages
- **Glassmorphism UI** - Modern, elegant design with blur effects
- **Dark/Light Theme** - Toggle between themes

## Project Structure

```
groupie-tracker/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── handlers/               # HTTP request handlers
│   ├── home.go            # Home page handler
│   ├── artist.go          # Artist detail page handler
│   ├── search.go          # Search/filter handler
│   ├── suggestions.go     # Autocomplete suggestions endpoint
│   ├── locations.go       # Location list endpoint
│   ├── error.go           # Error page handler
│   ├── search_test.go     # Search/filter tests
│   └── suggestions_test.go # Suggestions tests
├── models/                 # Data structures
│   └── artist.go          # Artist and Relation models
├── services/               # Business logic
│   ├── api.go             # API client with caching
│   ├── geocode.go         # Geocoding service (Nominatim)
│   ├── api_test.go        # API service tests
│   └── geocode_test.go    # Geocoding tests
├── templates/              # HTML templates
│   ├── home.html          # Home page template
│   ├── artist.html        # Artist detail template
│   └── error.html         # Error page template
├── static/                 # Static assets
│   ├── style.css          # Stylesheet
│   └── script.js          # Client-side JavaScript
```

## Installation

### Prerequisites
- Go 1.24 or higher
- Internet connection (to fetch data from API)

### Steps

1. **Clone the repository**
```bash
git clone https://platform.zone01.gr/git/ckotsalas/groupie-tracker.git
cd groupie-tracker
```

2. **Build the application**
```bash
go build -o groupie-tracker
```

3. **Run the application**
```bash
./groupie-tracker
```

Or run directly without building:
```bash
go run main.go
```

4. **Open your browser**
```
http://localhost:8080
```

## Usage

### Home Page
- View all artists in a grid layout
- Use the search bar with autocomplete to find artists by name, member, location, or date
- Apply filters: creation year range, first album year range, member count, locations
- Sort artists using the dropdown menu
- Filters auto-apply on change, or click "Apply Filters"
- Click "Reset" to clear all filters

### Artist Profile
- Click on any artist card to view their detailed profile
- See comprehensive statistics and information
- View interactive world tour map with geocoded concert markers
- Browse concert history organized by location
- View all concert dates for each location

### Search & Filter Examples

**Search by artist name:**
```
Search: "Queen"
```

**Search by member name:**
```
Search: "Freddie Mercury"
```

**Filter by year range:**
```
Min Year: 1970
Max Year: 1980
```

**Search by location:**
```
Search: "london"
```

**Filter by member count:**
```
Members: 4, 5 (checkboxes)
```

**Combine search and filters:**
```
Search: "Rock"
Min Year: 1960
Max Year: 1990
Sort By: Newest First
Members: 4
```

## API Endpoints

### External API (Data Source)
Base URL: `https://groupietrackers.herokuapp.com/api`

- `GET /artists` - List of all artists
- `GET /locations` - Concert locations
- `GET /dates` - Concert dates
- `GET /relation` - Relations between artists, dates, and locations

### Application Routes

- `GET /` - Home page (list all artists)
- `GET /artist/{id}` - Artist detail page
- `GET /api/search` - Search and filter API
  - Query parameters:
    - `q` - Search query (artist name, member, location, date)
    - `minYear` - Minimum formation year
    - `maxYear` - Maximum formation year
    - `minAlbumYear` - Minimum first album year
    - `maxAlbumYear` - Maximum first album year
    - `members` - Comma-separated member counts (e.g., "1,4,8")
    - `locations` - Comma-separated location keys
    - `sort` - Sort order (name, newest, oldest)
- `GET /api/suggestions` - Autocomplete suggestions
  - Query parameters:
    - `q` - Search query (returns max 10 categorized suggestions)
- `GET /api/locations` - All unique locations grouped by country
- `GET /static/*` - Static files (CSS, JS)

## Technologies

### Backend
- **Go 1.24** - Programming language
- **net/http** - HTTP server and client
- **html/template** - HTML templating
- **encoding/json** - JSON parsing

### Frontend
- **HTML5** - Markup
- **CSS3** - Styling (Glassmorphism design)
- **JavaScript (ES6+)** - Client-side interactivity
- **Fetch API** - AJAX requests
- **Leaflet.js** - Interactive maps (OpenStreetMap tiles)

### Standards
- **RESTful API** - API design pattern
- **Client-Server Architecture** - Application architecture
- **Responsive Design** - Mobile-first approach

## Testing

### Run All Tests
```bash
go test ./... -v
```

### Run Specific Package Tests
```bash
# Test handlers
go test ./handlers -v

# Test services
go test ./services -v
```

### Test Coverage
```bash
go test ./... -cover
```

### Current Test Suite
- ✅ API fetching (artists, relations, all relations)
- ✅ Search functionality (name, member, location, date, creation year)
- ✅ Year filtering (creation year and first album year)
- ✅ Member count filtering (including 8+ logic)
- ✅ Location filtering (with parent-region matching)
- ✅ Sorting algorithms
- ✅ Autocomplete suggestions (dedup, categories, cap)
- ✅ Geocoding (location parsing, title case)
- ✅ Error handling

## Configuration

### Port Configuration
Default port: `8080`

To change the port, modify `main.go`:
```go
http.ListenAndServe(":YOUR_PORT", nil)
```

### API Timeout
Default timeout: `10 seconds`

To change the timeout, modify `services/api.go`:
```go
client := &http.Client{Timeout: YOUR_DURATION * time.Second}
```

## Error Handling

The application handles various error scenarios:
- Network failures
- API unavailability
- Invalid user input
- Missing data
- Template parsing errors
- 404 Not Found
- 500 Internal Server Error

All errors display user-friendly messages.

---