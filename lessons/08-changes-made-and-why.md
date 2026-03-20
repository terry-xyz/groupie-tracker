# Lesson 08: Changes Made and Why

## What Happened Here?

The original codebase was a solid foundation — it fetched artists, displayed them in a grid, and had basic search and year filtering. But the project specification required **four extension modules** that were missing or incomplete: advanced filters, a full search bar, geolocalization (maps), and visualization improvements.

This lesson walks through every change that was made, **why** it was needed, and **where** to find it. Think of it as the "director's commentary" for the codebase.

---

## The Big Picture: Before and After

```
BEFORE                                    AFTER
──────                                    ─────
2 filter types (search + year range)  →   6 filter types (search, 2 ranges, 2 checkboxes, sort)
Search by name + members only         →   Search by name, members, locations, dates
No autocomplete                       →   Live suggestions with categories
No map                                →   Interactive world tour map per artist
No caching                            →   5-minute in-memory cache for API data
2 test files (9 tests)                →   4 test files (16 tests)
3 handler files                       →   5 handler files
1 service file                        →   2 service files
```

---

## Change 1: API Caching Layer

### What Changed
`services/api.go` — Added in-memory caching with a 5-minute TTL, plus a new `FetchAllRelations()` function.

### Why
The original code called the external API on **every single request**. With autocomplete firing on every keystroke and filters updating live, that's potentially dozens of API calls per second. The external API would slow down or rate-limit us.

**Real-world analogy:** Instead of driving to the grocery store every time you need an egg, you keep a carton in the fridge and restock every few days.

### How It Works

```go
// Before: fresh API call every time
artists, err := services.FetchArtists()

// After: checks cache first, fetches only if stale
artists, err := services.GetArtists()
```

The cache uses a `sync.Mutex` (a lock) because Go's web server handles requests concurrently — two requests arriving at the same time could both try to update the cache. The lock ensures only one goes at a time.

```go
var (
    cacheMu        sync.Mutex          // The lock
    cachedArtists  []models.Artist     // The stored data
    artistsCacheAt time.Time           // When we last fetched
    cacheTTL       = 5 * time.Minute   // How long before data is "stale"
)

func GetArtists() ([]models.Artist, error) {
    cacheMu.Lock()                     // Acquire the lock
    defer cacheMu.Unlock()             // Release when done

    if cachedArtists != nil && time.Since(artistsCacheAt) < cacheTTL {
        return cachedArtists, nil      // Cache hit — return instantly
    }

    artists, err := FetchArtists()     // Cache miss — fetch fresh data
    if err != nil {
        return nil, err
    }
    cachedArtists = artists
    artistsCacheAt = time.Now()
    return artists, nil
}
```

### Where
- `services/api.go:14-22` — Cache variables
- `services/api.go:63-83` — `FetchAllRelations()` (new)
- `services/api.go:85-119` — `GetArtists()` and `GetAllRelations()` (cached wrappers)

---

## Change 2: Three New Filter Functions

### What Changed
`handlers/search.go` — Added filters for first album year, member count, and concert locations.

### Why
The project spec requires **at least 4 filters** including both **range** and **checkbox** types. The original code only had text search and creation year range. Missing:

| Filter | Type | Requirement |
|--------|------|-------------|
| First album year range | Range | "Filter by first album date" |
| Number of members | Checkbox | "Filter by number of members" |
| Concert locations | Checkbox | "Filter by concert locations" |

### How They Work

**First Album Year Filter:**

The tricky part is that `FirstAlbum` is a string like `"14-02-1992"`, not a number. We need to extract the year first.

```go
// Parse "14-02-1992" → 1992
func parseFirstAlbumYear(firstAlbum string) int {
    parts := strings.Split(firstAlbum, "-")
    if len(parts) < 3 {
        return 0  // Malformed date — can't extract year
    }
    year, err := strconv.Atoi(parts[len(parts)-1])  // Last part is the year
    if err != nil {
        return 0
    }
    return year
}
```

**Member Count Filter:**

Users check boxes like "1", "2", "3", ..., "8+". The value `8` is special — it means "8 or more."

```go
func matchesMembersFilter(artist models.Artist, memberCounts []int) bool {
    if len(memberCounts) == 0 {
        return true  // No checkboxes checked = no filter
    }
    count := len(artist.Members)
    for _, mc := range memberCounts {
        if mc == 8 && count >= 8 {
            return true  // "8+" means 8, 9, 10, etc.
        }
        if count == mc {
            return true  // Exact match
        }
    }
    return false
}
```

**Location Filter with Parent-Region Matching:**

This is the most interesting filter. The spec says: "Seattle, Washington, USA can be found through Seattle, Washington, USA, or USA." We use **suffix matching** — selecting "usa" matches any location ending in "-usa".

```go
func matchesLocationFilter(relation *models.Relation, selectedLocations []string) bool {
    // ...
    for location := range relation.DatesLocations {
        for _, selected := range selectedLocations {
            if loc == sel {
                return true  // Exact match: "london-uk" == "london-uk"
            }
            if strings.HasSuffix(loc, "-"+sel) {
                return true  // Parent match: "north_carolina-usa" ends with "-usa"
            }
        }
    }
    return false
}
```

**Why `HasSuffix` instead of `Contains`?** Because `Contains` would cause "uk" to match "fuk**uk**a-japan". Suffix matching is precise — "uk" only matches locations that **end** with "-uk".

### Where
- `handlers/search.go:158-169` — `parseFirstAlbumYear`
- `handlers/search.go:171-197` — `matchesAlbumYearFilter`
- `handlers/search.go:199-216` — `matchesMembersFilter`
- `handlers/search.go:218-242` — `matchesLocationFilter`

---

## Change 3: Expanded Search Handler

### What Changed
`handlers/search.go` — The `SearchHandler` now reads 4 new query parameters and applies all filters with AND logic. The `matchesSearch` function now also searches locations, creation dates, and first album dates.

### Why
The spec says the search bar must find artists by **name, members, locations, first album date, AND creation date**. The original only searched name + members.

### How It Works

The filter pipeline is the same pattern as before (see [Lesson 03](03-patterns.md)), just with more stages:

```go
for _, artist := range artists {
    rel := relationMap[artist.ID]      // Get this artist's concert data

    if !matchesSearch(artist, rel, query) { continue }
    if !matchesYearFilter(artist, minYear, maxYear) { continue }
    if !matchesAlbumYearFilter(artist, minAlbumYear, maxAlbumYear) { continue }
    if !matchesMembersFilter(artist, memberCounts) { continue }
    if !matchesLocationFilter(rel, selectedLocations) { continue }

    results = append(results, artist)
}
```

Each `continue` skips to the next artist if the current one doesn't pass a filter. Think of it as a series of checkpoints — an artist must pass ALL of them to appear in results.

**Expanded search:** `matchesSearch` now takes a `relation` parameter and checks more fields:

```go
func matchesSearch(artist models.Artist, relation *models.Relation, query string) bool {
    // Check artist name          → "queen" matches Queen
    // Check member names         → "freddie" matches Queen (via Freddie Mercury)
    // Check creation date        → "1973" matches ACDC (formed in 1973)
    // Check first album date     → "05-08-1967" matches Pink Floyd
    // Check concert locations    → "london" matches any artist with London concerts
}
```

**Graceful degradation:** If the relations API fails, the search still works for everything except locations. The relation is just `nil`, and location matching returns `false` (no crash).

### Where
- `handlers/search.go:14-96` — Expanded `SearchHandler`
- `handlers/search.go:98-136` — Expanded `matchesSearch`

---

## Change 4: Autocomplete Suggestions

### What Changed
New file `handlers/suggestions.go` — A `/api/suggestions` endpoint that returns categorized search suggestions as users type.

### Why
The spec requires: "Typing suggestions must identify the search type." When you type "phil", you should see suggestions like:

```
Phil Collins    — artist/band
Phil Collins    — member
```

And when you type "japan":

```
saitama-japan   — location
osaka-japan     — location
nagoya-japan    — location
```

### How It Works

```go
type Suggestion struct {
    Text     string `json:"text"`
    Category string `json:"category"`  // "artist/band", "member", "location", etc.
}
```

The handler scans across all data sources for substring matches:

```
User types "phil"
    │
    ├─ Scan artist names        → "Phil Collins" matches → add {Phil Collins, "artist/band"}
    ├─ Scan member names        → "Phil Collins" in Genesis → add {Phil Collins, "member"}
    ├─ Scan creation dates      → no match
    ├─ Scan first album dates   → no match
    └─ Scan location keys       → "philadelphia-usa" matches → add {philadelphia-usa, "location"}
    │
    ▼
   Return max 10 suggestions, deduplicated
```

**Deduplication** prevents "Phil Collins — member" from appearing 10 times if he's listed as a member in multiple contexts. A `map[string]bool` keyed on `"text|category"` handles this.

**Cap at 10** keeps the dropdown manageable and stops the loop early once 10 suggestions are found (performance).

### Where
- `handlers/suggestions.go` — The full handler
- `static/script.js:62-116` — Frontend autocomplete (debounced fetch, dropdown rendering)
- `templates/home.html:17-20` — The dropdown container in HTML

---

## Change 5: Location Endpoint

### What Changed
New file `handlers/locations.go` — A `/api/locations` endpoint that returns all concert locations grouped by country.

### Why
The location filter uses checkboxes, and those checkboxes need to be populated with actual data. We can't hardcode locations because they come from the API and could change. So the frontend fetches available locations on page load and builds checkboxes dynamically.

### How It Works

```go
// API returns locations like:
//   "north_carolina-usa", "london-uk", "saitama-japan", "paris-france"
//
// We group them by country (last segment after "-"):
//   "usa"    → ["north_carolina-usa", "texas-usa", ...]
//   "uk"     → ["london-uk", "manchester-uk", ...]
//   "japan"  → ["saitama-japan", "osaka-japan", ...]
```

The frontend renders these as collapsible groups with a search box to filter the (potentially long) list:

```
🇺🇸 USA
  ☐ north_carolina-usa
  ☐ texas-usa
  ☐ new_york-usa
🇬🇧 UK
  ☐ london-uk
  ☐ manchester-uk
```

### Where
- `handlers/locations.go` — The endpoint
- `static/script.js:142-196` — Frontend: `loadLocations()`, location search filter
- `templates/home.html:59-66` — The checkbox container

---

## Change 6: Geocoding Service (Maps)

### What Changed
New file `services/geocode.go` — A geocoding service that converts location names to latitude/longitude coordinates using OpenStreetMap's Nominatim API.

### Why
The spec requires mapping concert locations on an actual map. To place markers, we need coordinates. The API only gives us location names like `"paris-france"`, not coordinates like `(48.8566, 2.3522)`.

**Real-world analogy:** You have an address ("123 Main Street"), but GPS needs coordinates. Geocoding is the conversion from one to the other.

### How It Works

```
"north_carolina-usa"
        │
        ▼ FormatLocationName()
"North Carolina, Usa"
        │
        ▼ geocodeAddress() → Nominatim API
{lat: 35.78, lng: -80.79}
        │
        ▼ Stored in GeoLocation struct
{Lat: 35.78, Lng: -80.79, City: "North Carolina", Country: "Usa", Dates: [...]}
```

**Three important safeguards:**

1. **Rate limiter** — Nominatim's free tier allows 1 request per second. A channel-based rate limiter feeds one token every 1.1 seconds. Without this, Nominatim would block us.

```go
var rateLimiter = make(chan struct{}, 1)

func init() {
    go func() {
        for {
            rateLimiter <- struct{}{}           // Feed one token
            time.Sleep(1100 * time.Millisecond) // Wait 1.1 seconds
        }
    }()
}

// In geocodeAddress():
<-rateLimiter  // Wait for a token before making the request
```

2. **Cache** — Once we geocode "paris-france", we store the result forever (coordinates don't change). A `sync.RWMutex` allows multiple readers OR one writer — so concurrent requests can read the cache simultaneously.

3. **Graceful failure** — If a location can't be geocoded (Nominatim doesn't recognize it), we skip it instead of crashing. The map just shows fewer markers.

### Where
- `services/geocode.go` — The full geocoding service (rate limiter, coordinate cache, disk persistence)
- `handlers/artist_geo.go` — The `/api/artist-geo?id=X` endpoint (geocoding moved here, out of the page handler)
- `templates/artist.html` — Map initialization script (fetches from `/api/artist-geo` asynchronously)

---

## Change 7: Interactive Map on Artist Pages

### What Changed
`templates/artist.html` — Added a Leaflet.js map between the info cards and concert history.

### Why
The geolocalization spec requires that concert locations appear on an actual map with markers. We chose Leaflet + OpenStreetMap because they're free, require no API key, and load from a CDN (no Go packages needed).

### How It Works

```html
<!-- Load Leaflet from CDN -->
<link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" />
<script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>

<!-- Map container -->
<div id="concert-map"></div>
```

The inline script at the bottom:
1. On `window.load`, calls `fetch('/api/artist-geo?id={{.Artist.ID}}')` — does NOT block page render
2. When the response arrives, calls `initMap(locations)`
3. `initMap` creates a Leaflet map with OpenStreetMap tiles
4. Places a marker at each concert location with a popup showing city, country, and dates
5. Draws a dashed purple polyline connecting all markers (the "tour route")
6. Auto-zooms to fit all markers with `fitBounds`

**Why async?** Geocoding via Nominatim takes 1+ second per location. An artist with 8 locations would stall the page for 8+ seconds if done synchronously. By fetching from `/api/artist-geo` after load, the page always renders instantly. The map just fills in a moment later (or instantly, if coordinates were cached from a previous visit).

### Where
- `templates/artist.html` — Leaflet CDN imports, map section HTML, async fetch + map init
- `handlers/artist_geo.go` — The JSON endpoint that runs geocoding
- `services/geocode.go` — `GetGeoLocations()` with per-artist + disk caching

---

## Change 8: Expanded Filter UI

### What Changed
`templates/home.html` — Added first album year range inputs, member count checkboxes, location checkboxes with search, and a result count display.

### Why
The spec requires 4 filter types with both **range** and **checkbox** interfaces. The original only had text search + year range + sort.

### What Was Added

| UI Element | Filter Type | Where in HTML |
|-----------|-------------|---------------|
| First Album Year (min/max) | Range | Lines 28-32 |
| Members checkboxes (1-8+) | Checkbox | Lines 44-57 |
| Location checkboxes (dynamic) | Checkbox | Lines 59-66 |
| Result count ("Showing X artists") | Feedback | Line 74 |

The member checkboxes are hardcoded (the dataset has artists with 1-8+ members), but the location checkboxes are populated dynamically by JavaScript calling `/api/locations` on page load.

### Where
- `templates/home.html` — The full updated template

---

## Change 9: JavaScript Expansion

### What Changed
`static/script.js` — Major expansion from 125 lines to 356 lines, adding autocomplete, location loading, expanded filter collection, and async filter triggers.

### Why
Three new frontend behaviors were needed:
1. **Autocomplete** — Show live suggestions as users type
2. **Dynamic checkboxes** — Load and render location checkboxes
3. **Live filtering** — Update results as filters change (no need to click "Apply")

### Key Additions

**Debounce utility** — Prevents firing a request on every single keystroke. Waits until the user pauses typing (250ms for search, 400ms for number inputs):

```javascript
function debounce(fn, delay) {
    let timer;
    return function() {
        clearTimeout(timer);
        timer = setTimeout(() => fn.apply(this, arguments), delay);
    };
}
```

**XSS protection** — All dynamic HTML now uses `escapeHTML()` and `escapeAttr()` instead of raw template literals. This prevents malicious data from executing as code:

```javascript
// Before (vulnerable):
`<h2>${artist.name}</h2>`

// After (safe):
'<h2>' + escapeHTML(artist.name) + '</h2>'
```

**Auto-trigger on change** — Checkboxes and dropdowns call `applyFilters()` immediately on change. Number inputs are debounced. This means results update live without clicking "Apply":

```javascript
sortBySelect.addEventListener('change', applyFilters);
document.querySelectorAll('input[name="members"]').forEach(cb => {
    cb.addEventListener('change', applyFilters);
});
```

### Where
- `static/script.js:48-59` — Debounce utility
- `static/script.js:62-116` — Autocomplete system
- `static/script.js:142-196` — Location loading and filtering
- `static/script.js:207-223` — Auto-trigger listeners
- `static/script.js:237-294` — Expanded `applyFilters()`

---

## Change 10: New CSS Components

### What Changed
`static/style.css` — Added styles for autocomplete dropdown, checkbox filters, location search, map section, and result count. All with dark theme variants.

### Why
Every new UI component needs styling that matches the existing glassmorphism design. Without CSS, the checkboxes, dropdown, and map would look like unstyled browser defaults — breaking visual consistency (Shneiderman's Rule 1).

### What Was Added

| Component | Key Styles |
|-----------|-----------|
| `.search-container` | `position: relative` (anchors the dropdown) |
| `.suggestions-dropdown` | Absolute positioning, glassmorphism, max-height scroll |
| `.suggestion-item` | Hover highlight, flex layout with category label |
| `.checkbox-group` | Flex-wrap for inline checkboxes |
| `.checkbox-list.scrollable` | Max-height 200px with overflow scroll |
| `.location-search` | Small filter input above location checkboxes |
| `.map-section` / `#concert-map` | 450px height (300px mobile), border-radius |
| `.result-count` | Subtle text above the grid |

Every component has a `body.dark-theme` variant so the design stays consistent in both themes.

### Where
- `static/style.css` — Search the file for `.suggestions-`, `.checkbox-`, `.map-section`, `.location-`, and `.result-count`

---

## Change 11: Expanded Test Suite

### What Changed
Added 2 new test files and expanded the existing search tests from 3 to 8 test functions.

### Why
Every new filter function needs tests to prove it works correctly. The audit checks specific scenarios like "6 members should return Pink Floyd, Arctic Monkeys, Linkin Park, Foo Fighters." Tests verify these work without manually clicking through the UI.

### Test Summary

| File | Tests | What They Verify |
|------|-------|-----------------|
| `handlers/search_test.go` | 8 | Text search (name, members, dates, locations), year filter, album year filter and parsing, member count filter (including 8+), location filter (including parent-region), nil relation handling, sort modes |
| `handlers/suggestions_test.go` | 3 | Empty query returns `[]`, real query returns categorized results, max 10 cap enforced |
| `services/api_test.go` | 3 | API connectivity, relation fetch, all-relations fetch |
| `services/geocode_test.go` | 2 | Location name parsing, custom TitleCase function |

### Where
- `handlers/search_test.go` — All filter function tests
- `handlers/suggestions_test.go` — Suggestion handler tests
- `services/geocode_test.go` — Geocoding utility tests
- `services/api_test.go` — API integration tests (hits real API)

---

## Change 12: Route Registration

### What Changed
`main.go` — Two new routes registered.

### Why
New endpoints need to be wired to their handlers. Without registration, the server doesn't know these URLs exist.

```go
// Before: 3 routes
http.HandleFunc("/", handlers.HomeHandler)
http.HandleFunc("/artist/", handlers.ArtistHandler)
http.HandleFunc("/api/search", handlers.SearchHandler)

// After: 6 routes
http.HandleFunc("/", handlers.HomeHandler)
http.HandleFunc("/artist/", handlers.ArtistHandler)
http.HandleFunc("/api/search", handlers.SearchHandler)
http.HandleFunc("/api/suggestions", handlers.SuggestionsHandler)  // autocomplete
http.HandleFunc("/api/locations", handlers.LocationsHandler)      // location filter data
http.HandleFunc("/api/artist-geo", handlers.ArtistGeoHandler)     // async geocoding
```

### Where
- `main.go` — route registration + startup pre-warming goroutine

---

## Files Changed vs. Created

| Action | File | What |
|--------|------|------|
| Modified | `services/api.go` | Added caching + `FetchAllRelations` |
| Modified | `handlers/search.go` | 3 new filters + expanded search |
| Modified | `handlers/artist.go` | Uses cached APIs; async map (no geocoding here) |
| Modified | `main.go` | 3 new routes + geocode cache load + startup pre-warming |
| Modified | `templates/home.html` | New filter UI |
| Modified | `templates/artist.html` | Async map via `/api/artist-geo` fetch |
| Modified | `static/script.js` | Autocomplete + expanded filters |
| Modified | `static/style.css` | New component styles |
| Modified | `handlers/search_test.go` | 5 new test functions |
| Modified | `services/api_test.go` | 1 new test function |
| **Created** | `handlers/suggestions.go` | Autocomplete endpoint |
| **Created** | `handlers/locations.go` | Location grouping endpoint |
| **Created** | `handlers/artist_geo.go` | Async geocoding JSON endpoint |
| **Created** | `services/geocode.go` | Geocoding service with caching + disk persistence |
| **Created** | `handlers/suggestions_test.go` | Suggestion tests |
| **Created** | `services/geocode_test.go` | Geocoding tests |

---

## Change 13: Performance — Async Geocoding & Startup Pre-Warming

### What Changed
`handlers/artist.go`, `handlers/artist_geo.go` (new), `services/geocode.go`, `main.go`, `templates/artist.html`.

### Why
"View Details" was taking multiple seconds because `ArtistHandler` was synchronously geocoding each concert location before rendering — each uncached location waited 1.1s for the Nominatim rate limiter. An artist with 8 locations = ~9 seconds of blocking.

### How It Works

**Four layered improvements:**

1. **Async geocoding** — Geocoding moved out of the page handler into `handlers/artist_geo.go`. The browser fetches `/api/artist-geo?id=X` after page load. The page is always instant; the map fills in when ready.

2. **Per-artist cache** — `GetGeoLocations(artistID, ...)` in `services/geocode.go` stores the `[]GeoLocation` result per artist. Revisiting the same artist page: instant map.

3. **Disk persistence** — `saveGeocodeCache("data/geocode_cache.json")` writes every new coordinate to disk after it's geocoded. `LoadGeocodeCache` reads it back at startup. After the first full pre-warm, subsequent server restarts skip all Nominatim calls for known locations.

4. **Startup pre-warming** — A goroutine in `main.go` runs after the server starts, warming the artist/relation caches concurrently, then geocoding all artists in the background. If a user clicks "View Details" before their artist is pre-warmed, it geocodes on-demand (same as before) but caches the result immediately.

**Also:**
- Templates (`home.html`, `artist.html`) are now parsed once at startup using `sync.Once`, not on every request.
- `ArtistHandler` now uses `GetArtists()` and `GetAllRelations()` (cached) instead of raw `FetchArtists()` / `FetchRelation()`.

### Where
- `handlers/artist_geo.go` — new async geocoding endpoint
- `services/geocode.go` — `GetGeoLocations()`, `LoadGeocodeCache()`, `saveGeocodeCache()`
- `main.go` — `LoadGeocodeCache` at startup + pre-warming goroutine
- `handlers/home.go` and `handlers/artist.go` — `sync.Once` lazy template init, cached API calls
- `templates/artist.html` — `fetch('/api/artist-geo?id=...')` after page load

---

## How Each Spec Module Maps to Changes

| Spec Module | Changes That Address It |
|-------------|------------------------|
| **Filters** | Changes 2, 3, 8, 9 (new filter functions, expanded handler, checkbox UI, JS collection) |
| **Search Bar** | Changes 3, 4, 9 (expanded matchesSearch, suggestions endpoint, autocomplete JS) |
| **Geolocalization** | Changes 6, 7 (geocoding service, Leaflet map on artist page) |
| **Visualizations** | Changes 8, 9, 10 (result count, live filter updates, glassmorphism components) |
| **Good Practices** | Changes 1, 11 (caching for performance, comprehensive tests) |

---

## What's Next?

You now know every change, why it was made, and where to find it. If you want to go deeper:

- **Understand the patterns** behind these changes → [Lesson 03](03-patterns.md)
- **Read the functions line-by-line** → [Lesson 04](04-line-by-line.md)
- **Try modifying something yourself** → [Lesson 05](05-exercises.md)
- **Watch out for common mistakes** → [Lesson 06](06-gotchas.md)
