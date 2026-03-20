# Lesson 00: Project Overview

## Your Goal

You've built (or inherited) **Groupie Tracker** — a web app that shows music artists and their concert data. Maybe you used AI to write parts of it ("vibecoding"), or maybe you're just getting started with Go. Either way, this guide will take you from **"it works, I think"** to **"I know exactly what every line does and why."**

By the end of these lessons, you'll be able to:
- Explain the entire project in one sentence
- Draw the data flow from memory
- Find the right file for any change
- Debug problems without asking an AI
- Read Go code fluently

---

## What is Code? (If You're New)

Think of code like a **recipe**:
- **Ingredients** = Data (artist names, concert dates, images)
- **Instructions** = Code (fetch data, filter it, display it)
- **Kitchen** = Your computer (runs the instructions)
- **Served Dish** = The web page the user sees

Every programming language uses the same 5 building blocks:

| Building Block | Plain English | Go Example |
|----------------|---------------|------------|
| **Variable** | A labeled box that holds data | `name := "Queen"` |
| **Function** | A reusable recipe | `func FetchArtists() {...}` |
| **Condition** | A yes/no question | `if err != nil {...}` |
| **Loop** | Repeat until done | `for _, artist := range artists {...}` |
| **Data Structure** | An organized container | `type Artist struct {...}` |

That's it. Every program — from this one to Google — is built from these five pieces.

---

## What Does This Project Do?

**One sentence:** Groupie Tracker is a Go web app that fetches music artist data from an external API and displays it with search, filtering, and detailed artist profiles.

Think of it like a **mini Spotify profile browser**:
1. You open the home page → see a grid of artist cards
2. You can search by name or band member
3. You can filter by formation year
4. You click an artist → see their full profile with concert history

---

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    USER'S BROWSER                       │
│                                                         │
│  ┌──────────┐  ┌──────────────┐  ┌───────────────────┐  │
│  │ home.html│  │ artist.html  │  │ script.js         │  │
│  │ (grid)   │  │ (profile)    │  │ (search/filter)   │  │
│  └────┬─────┘  └──────┬───────┘  └────────┬──────────┘  │
│       │               │                   │              │
└───────┼───────────────┼───────────────────┼──────────────┘
        │ GET /         │ GET /artist/3     │ GET /api/search
        ▼               ▼                   ▼
┌─────────────────────────────────────────────────────────┐
│                    GO SERVER (:8080)                     │
│                                                         │
│  main.go ─── Routes requests to handlers                │
│       │                                                 │
│       ├── handlers/home.go    → HomeHandler              │
│       ├── handlers/artist.go  → ArtistHandler            │
│       ├── handlers/search.go  → SearchHandler            │
│       └── handlers/error.go   → ErrorHandler             │
│                    │                                    │
│              services/api.go                            │
│              (fetches external data)                    │
│                    │                                    │
└────────────────────┼────────────────────────────────────┘
                     │ HTTP GET
                     ▼
┌─────────────────────────────────────────────────────────┐
│         EXTERNAL API (herokuapp.com)                    │
│                                                         │
│  /api/artists   → All artist data (JSON)                │
│  /api/relation  → Concert locations & dates (JSON)      │
└─────────────────────────────────────────────────────────┘
```

---

## Directory Map

```
groupie-tracker/
│
├── main.go                 ← THE FRONT DOOR: starts the server, defines routes
│
├── models/
│   └── artist.go           ← THE BLUEPRINT: defines what an Artist looks like
│
├── services/
│   ├── api.go              ← THE DELIVERY TRUCK: fetches data from external API (with 5-min cache)
│   ├── api_test.go         ← Tests for the delivery truck
│   ├── geocode.go          ← THE CARTOGRAPHER: converts location names to coordinates
│   └── geocode_test.go     ← Tests for geocoding utilities
│
├── handlers/
│   ├── home.go             ← THE RECEPTIONIST: serves the home page
│   ├── artist.go           ← THE TOUR GUIDE: serves artist detail pages
│   ├── artist_geo.go       ← THE MAP API: returns geocoded locations as JSON (async)
│   ├── search.go           ← THE LIBRARIAN: handles search and filtering
│   ├── locations.go        ← THE DIRECTORY: returns all locations grouped by country
│   ├── suggestions.go      ← THE AUTOCOMPLETE: returns live search suggestions
│   ├── error.go            ← THE APOLOGY NOTE: shows error pages
│   └── search_test.go      ← Tests for search logic
│
├── templates/
│   ├── home.html           ← THE STOREFRONT: what the home page looks like
│   ├── artist.html         ← THE PROFILE PAGE: what artist details look like
│   └── error.html          ← THE ERROR SIGN: what errors look like
│
├── static/
│   ├── style.css           ← THE PAINT & WALLPAPER: colors, layout, animations
│   └── script.js           ← THE INTERACTIVE BUTTONS: search, theme toggle, filters
│
├── data/
│   └── geocode_cache.json  ← Persisted coordinate cache (auto-created on first geocode)
│
├── docs/                   ← Project documentation
├── go.mod                  ← Go module definition (like a package.json)
├── .gitignore              ← Files Git should ignore
├── LICENSE                 ← MIT License
└── README.md               ← Project readme
```

### Directory Responsibilities in Plain English

| Directory | Role | Analogy |
|-----------|------|---------|
| `main.go` | Starts the server, connects everything | The power switch |
| `models/` | Defines data shapes | The cookie cutter |
| `services/` | Talks to external APIs | The delivery driver |
| `handlers/` | Responds to web requests | The waitstaff |
| `templates/` | HTML page layouts | The plate presentation |
| `static/` | CSS and JavaScript | The decoration and interactivity |

---

## Entry and Exit Points

### Entry Points (Where things START)

1. **Server startup:** `main.go:main()` — everything begins here
2. **Home page request:** `GET /` → `handlers.HomeHandler`
3. **Artist page request:** `GET /artist/{id}` → `handlers.ArtistHandler`
4. **Search API request:** `GET /api/search?q=...` → `handlers.SearchHandler`
5. **Autocomplete API:** `GET /api/suggestions?q=...` → `handlers.SuggestionsHandler`
6. **Locations API:** `GET /api/locations` → `handlers.LocationsHandler`
7. **Artist geocoding API:** `GET /api/artist-geo?id=X` → `handlers.ArtistGeoHandler`

### Exit Points (Where things END)

1. **HTML response** — rendered template sent to browser (home, artist, error pages)
2. **JSON response** — search results returned to JavaScript
3. **Error page** — when something goes wrong (404, 500)

---

## How to Run This Project

```bash
# Option 1: Run directly
go run main.go

# Option 2: Build then run
go build -o groupie-tracker
./groupie-tracker

# Then open http://localhost:8080 in your browser
```

### Running Tests

```bash
go test ./...          # Run all tests
go test ./... -v       # Verbose (see each test name)
go test ./... -cover   # See code coverage percentage
```

---

## Dependencies

**Zero external dependencies.** This project uses only Go's standard library:

| Package | What It Does |
|---------|-------------|
| `net/http` | Web server and HTTP client |
| `html/template` | HTML templating |
| `encoding/json` | JSON parsing |
| `log` | Logging |
| `fmt` | String formatting |
| `strconv` | String-to-number conversions |
| `strings` | String manipulation |
| `sort` | Sorting |
| `time` | Timeouts |

This means: no package manager drama, no version conflicts, no supply chain worries.

---

## What's Next?

| Lesson | What You'll Learn |
|--------|-------------------|
| [01 - Core Concepts](01-core-concepts.md) | Go syntax basics and the key concepts in this codebase |
| [02 - Data Flow](02-data-flow.md) | How data moves from API → server → browser |
| [03 - Patterns](03-patterns.md) | Reusable design patterns used here |
| [04 - Line by Line](04-line-by-line.md) | Detailed walkthrough of every important function |
| [05 - Exercises](05-exercises.md) | Hands-on practice from beginner to advanced |
| [06 - Gotchas](06-gotchas.md) | Common mistakes and how to avoid them |
| [07 - Glossary](07-glossary.md) | Every term, abbreviation, and concept decoded |
