# Lesson 02: Data Flow

## The Factory Assembly Line

Think of Groupie Tracker like a **factory that builds artist profile pages**:

1. **Raw materials arrive** — JSON data from the external API
2. **Workers inspect and sort them** — Go parses, filters, and calculates
3. **Assembly line shapes them** — Templates fill in the HTML
4. **Finished product ships out** — The browser receives a complete web page

Every request follows this pattern. Let's trace each one.

---

## How to Read These Diagrams

- **Boxes** `[ ]` = things (data, components, files)
- **Arrows** `→` or `▼` = movement or transformation
- **Labels** on arrows = what's happening at that step
- **Dashed lines** `---` = network boundary (data travels over the internet)

---

## Flow 1: Loading the Home Page

**Trigger:** User opens `http://localhost:8080/`

```
┌──────────┐   GET /    ┌──────────────┐
│  Browser  │ ────────→ │   main.go     │
└──────────┘            │  (router)     │
                        └──────┬───────┘
                               │ routes to
                               ▼
                        ┌──────────────┐
                        │ HomeHandler   │
                        │ (home.go)     │
                        └──────┬───────┘
                               │ calls
                               ▼
                        ┌──────────────┐       HTTP GET        ┌─────────────────┐
                        │ FetchArtists │ ─────────────────────→│ External API     │
                        │ (api.go)     │                       │ /api/artists     │
                        └──────┬───────┘  ←─── JSON response   └─────────────────┘
                               │
                               │ []Artist (Go structs)
                               ▼
                        ┌──────────────┐
                        │  home.html   │
                        │  (template)  │
                        └──────┬───────┘
                               │ rendered HTML
                               ▼
                        ┌──────────────┐
                        │   Browser    │
                        │ (shows grid) │
                        └──────────────┘
```

### The Story in Plain English

1. **User opens the page** — browser sends `GET /` to our server
2. **Router checks the URL** — `"/"` matches `HomeHandler`
3. **HomeHandler asks for data** — calls `services.FetchArtists()`
4. **FetchArtists calls the external API** — HTTP GET to `groupietrackers.herokuapp.com/api/artists`
5. **API returns JSON** — a big list of artists in text format
6. **Go decodes the JSON** — text becomes `[]Artist` structs
7. **Template renders HTML** — `home.html` loops over artists, creating cards
8. **Browser displays the page** — user sees the artist grid

### Data Transformations (Step by Step)

**STEP 1: Raw JSON from API**
```json
[
  {
    "id": 1,
    "name": "Queen",
    "image": "https://groupietrackers.herokuapp.com/api/images/queen.jpeg",
    "members": ["Freddie Mercury", "Brian May", "Roger Taylor", "John Deacon"],
    "creationDate": 1970,
    "firstAlbum": "14-12-1973"
  },
  ...
]
```

**STEP 2: Go structs (after JSON decoding)**
```go
[]models.Artist{
    {
        ID:           1,
        Name:         "Queen",
        Image:        "https://...queen.jpeg",
        Members:      []string{"Freddie Mercury", "Brian May", "Roger Taylor", "John Deacon"},
        CreationDate: 1970,
        FirstAlbum:   "14-12-1973",
    },
    // ... more artists
}
```

**STEP 3: HTML output (after template rendering)**
```html
<div class="artist-card">
    <img src="https://...queen.jpeg" alt="Queen">
    <h2>Queen</h2>
    <p>Freddie Mercury, Brian May, Roger Taylor, John Deacon</p>
    <a href="/artist/1">View Details</a>
</div>
<!-- ... more cards -->
```

---

## Flow 2: Viewing an Artist Profile

**Trigger:** User clicks on an artist card (e.g., `/artist/3`)

```
┌──────────┐  GET /artist/3  ┌──────────────┐
│  Browser  │ ──────────────→ │   main.go     │
└──────────┘                  │  (router)     │
                              └──────┬───────┘
                                     │ routes to
                                     ▼
                              ┌──────────────┐
                              │ ArtistHandler │
                              │ (artist.go)   │
                              └──────┬───────┘
                                     │
                    ┌────────────────┼────────────────┐
                    │ (1) Extract ID │                │
                    │ from URL path  │                │
                    ▼                ▼                │
             ┌─────────────┐  ┌─────────────┐       │
             │FetchArtists │  │FetchRelation│       │
             │ (all artists)│  │ (concerts)  │       │
             └──────┬──────┘  └──────┬──────┘       │
                    │                │               │
                    │ Find artist    │ Concert data  │
                    │ with ID=3      │               │
                    ▼                ▼               │
             ┌───────────────────────────────┐      │
             │    Calculate Statistics       │      │
             │  - Total concerts             │      │
             │  - Total countries            │      │
             │  - Years active               │      │
             │  - Band type (Solo/Duo/etc)   │      │
             └──────────────┬────────────────┘      │
                            │                       │
                            │ ArtistPageData        │
                            ▼                       │
                     ┌──────────────┐               │
                     │ artist.html  │               │
                     │ (template)   │               │
                     └──────┬───────┘               │
                            │ rendered HTML          │
                            ▼                       │
                     ┌──────────────┐               │
                     │   Browser    │               │
                     │ (profile)    │               │
                     └──────────────┘               │
```

### The Story in Plain English

1. **User clicks "View Details"** — browser navigates to `/artist/3`
2. **Router matches `/artist/`** — sends request to `ArtistHandler`
3. **Handler extracts "3" from the URL** — `strings.TrimPrefix(r.URL.Path, "/artist/")` → `"3"`
4. **Converts "3" to integer** — `strconv.Atoi("3")` → `3`
5. **Fetches ALL artists from cache** — `GetArtists()` returns the in-memory cached list; scans to find `ID == 3`
6. **Fetches ALL relations from cache** — `GetAllRelations()` returns cached relations; scans for `ID == 3`
7. **Calculates stats** — total concerts, countries, years active, band type
8. **Packs everything into `ArtistPageData`** — one struct with all the data the template needs
9. **Page renders instantly** — HTML with artist info, stats, and concert history
10. **Browser displays the page, then fetches map data** — after page load, JavaScript calls `GET /api/artist-geo?id=3` and renders the interactive map when coordinates arrive

### Data Transformations (Step by Step)

**STEP 1: URL path**
```
/artist/3
```

**STEP 2: Extracted ID**
```go
id := 3
```

**STEP 3: Found artist (from the full list)**
```go
Artist{ID: 3, Name: "Mamonas Assassinas", Members: [...], CreationDate: 1993, ...}
```

**STEP 4: Relation data (from separate API call)**
```go
Relation{
    ID: 3,
    DatesLocations: map[string][]string{
        "sao_paulo-brazil":    ["15-04-1995", "23-06-1995"],
        "rio_de_janeiro-brazil": ["01-05-1995"],
    },
}
```

**STEP 5: Calculated statistics**
```go
ArtistPageData{
    Artist:         Artist{...},
    Relation:       &Relation{...},
    TotalConcerts:  3,           // 2 + 1 = 3 concerts
    TotalCountries: 1,           // only "brazil"
    YearsActive:    33,          // 2026 - 1993
    BandType:       "Quintet",   // 5 members
}
```

**STEP 6: Rendered HTML (simplified)**
```html
<h1>Mamonas Assassinas</h1>
<span class="stat">3 Concerts</span>
<span class="stat">1 Countries</span>
<span class="stat">33 Years Active</span>
<span class="badge">Quintet</span>
<div class="concert-item">
    <strong>sao_paulo-brazil</strong>
    <span>15-04-1995</span>
    <span>23-06-1995</span>
</div>
```

---

## Flow 3: Searching and Filtering

**Trigger:** User types in the search box or changes filters

```
┌──────────────────────────────────────────────────────┐
│                    BROWSER                            │
│                                                      │
│  User types "queen" in search box                    │
│       │                                              │
│       ▼                                              │
│  script.js: applyFilters()                           │
│  Builds URL: /api/search?q=queen&minYear=&maxYear=   │
│       │                                              │
│       │  fetch() (AJAX request — no page reload)     │
└───────┼──────────────────────────────────────────────┘
        │
        ▼ GET /api/search?q=queen
┌──────────────────────────────────────────────────────┐
│                    GO SERVER                          │
│                                                      │
│  SearchHandler (search.go)                           │
│       │                                              │
│       ├─ 1. Parse query params: q="queen"            │
│       │                                              │
│       ├─ 2. FetchArtists() → get all 52 artists      │
│       │                                              │
│       ├─ 3. FILTER: loop through each artist          │
│       │    ├─ matchesSearch("queen")?                 │
│       │    │   ├─ "Queen" contains "queen"? YES ✓    │
│       │    │   ├─ "Bee Gees" contains "queen"? NO    │
│       │    │   └─ Check members too...               │
│       │    └─ matchesYearFilter(min, max)?            │
│       │        └─ No filter applied → YES ✓           │
│       │                                              │
│       ├─ 4. SORT: by default (ID order)              │
│       │                                              │
│       └─ 5. ENCODE: structs → JSON                   │
│                                                      │
└───────┬──────────────────────────────────────────────┘
        │
        ▼ JSON response
┌──────────────────────────────────────────────────────┐
│                    BROWSER                            │
│                                                      │
│  script.js receives JSON                             │
│       │                                              │
│       ▼                                              │
│  displayArtists(data)                                │
│  - Clears current grid                               │
│  - Creates new HTML cards from JSON                   │
│  - Inserts cards into the page                       │
│  - No page reload needed!                            │
│                                                      │
└──────────────────────────────────────────────────────┘
```

### The Story in Plain English

1. **User types "queen"** in the search box
2. **JavaScript intercepts** — `applyFilters()` runs on button click or Enter key
3. **JavaScript builds a URL** — `/api/search?q=queen&minYear=&maxYear=&sort=default`
4. **JavaScript sends AJAX request** — `fetch(url)` (no page reload)
5. **Go server receives the request** — `SearchHandler` parses the query parameters
6. **Server fetches ALL artists** — then filters them in memory
7. **Filter: search match** — checks if "queen" appears in artist name or member names (case-insensitive)
8. **Filter: year range** — if min/max years provided, exclude artists outside range
9. **Sort results** — by name, newest, oldest, or default (ID)
10. **Encode and send** — converts filtered list to JSON, sends to browser
11. **JavaScript updates the page** — `displayArtists()` rebuilds the card grid with matching results

### Data Transformations (Step by Step)

**STEP 1: User input**
```
Search box: "queen"
Min year: (empty)
Max year: (empty)
Sort: "default"
```

**STEP 2: URL built by JavaScript**
```
/api/search?q=queen&minYear=&maxYear=&sort=default
```

**STEP 3: All 52 artists fetched from API**
```go
[Artist{ID:1, Name:"Queen"}, Artist{ID:2, Name:"SOJA"}, ... 50 more]
```

**STEP 4: After filtering (only matches)**
```go
[Artist{ID:1, Name:"Queen"}]
```

**STEP 5: JSON sent to browser**
```json
[{"id":1,"name":"Queen","image":"...","members":["Freddie Mercury",...],"creationDate":1970,"firstAlbum":"14-12-1973"}]
```

**STEP 6: HTML generated by JavaScript**
```html
<div class="artist-card">
    <img src="..." alt="Queen">
    <h2>Queen</h2>
    <p>Freddie Mercury, Brian May, Roger Taylor, John Deacon</p>
    <a href="/artist/1">View Details</a>
</div>
```

---

## Flow 4: Error Handling

**Trigger:** Something goes wrong (bad URL, API down, etc.)

```
ERROR SCENARIOS AND THEIR FLOWS:

1. Bad URL: /artist/banana
   Browser → ArtistHandler → strconv.Atoi("banana") FAILS
                            → ErrorHandler(400, "Invalid artist ID")
                            → error.html rendered
                            → Browser shows "400: Invalid artist ID"

2. Artist not found: /artist/9999
   Browser → ArtistHandler → FetchArtists() OK
                            → Loop through artists... ID 9999 not found
                            → ErrorHandler(404, "Artist not found")
                            → error.html rendered
                            → Browser shows "404: Artist not found"

3. API is down: /
   Browser → HomeHandler → FetchArtists() FAILS (timeout)
                          → log.Printf("Error: connection timeout")
                          → ErrorHandler(500, "Unable to load artists...")
                          → error.html rendered
                          → Browser shows "500: Unable to load artists..."

4. Unknown page: /random
   Browser → HomeHandler → r.URL.Path != "/" → TRUE
                          → ErrorHandler(404)
                          → error.html rendered
                          → Browser shows "404: Not Found"
```

---

## The Complete Picture

Here's the entire application data flow in one diagram:

```
                         ┌─────────────────────────────┐
                         │      EXTERNAL API            │
                         │  groupietrackers.herokuapp   │
                         │                             │
                         │  /api/artists → Artist JSON  │
                         │  /api/relation → Concert JSON│
                         └─────────┬───────────────────┘
                                   │
                              HTTP GET + JSON
                                   │
┌──────────────────────────────────┼──────────────────────────────────┐
│                          GO SERVER                                  │
│                                  │                                  │
│  services/api.go ◄───────────────┘                                  │
│  (FetchArtists, FetchRelation)                                     │
│         │                                                           │
│         │ []Artist, *Relation                                       │
│         ▼                                                           │
│  handlers/*.go                                                     │
│  ├─ HomeHandler    → template → HTML (artist grid)                  │
│  ├─ ArtistHandler  → template → HTML (artist profile)               │
│  ├─ SearchHandler  → filter   → JSON (search results)              │
│  └─ ErrorHandler   → template → HTML (error page)                  │
│                                                                     │
│  models/artist.go                                                  │
│  (Artist, Relation — the data shapes)                              │
│                                                                     │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                        HTML or JSON
                             │
┌────────────────────────────┼────────────────────────────────────────┐
│                       BROWSER                                       │
│                            │                                        │
│  templates/ (rendered HTML)│                                        │
│  ├─ home.html   → Artist grid page                                  │
│  ├─ artist.html → Artist profile page                               │
│  └─ error.html  → Error message page                               │
│                                                                     │
│  static/style.css → Glassmorphism styling, themes                  │
│  static/script.js → Search, filter, theme toggle                   │
│                            │                                        │
│                    ┌───────┘                                        │
│                    │ AJAX (for search only)                         │
│                    │ GET /api/search?q=...                          │
│                    │ → receives JSON                                │
│                    │ → updates page without reload                  │
│                    ▼                                                │
│              User sees the result                                  │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Key Takeaways

1. **Two rendering strategies**: Server-side (Go templates for full pages) and client-side (JavaScript for search results and map data)
2. **Caching at multiple levels**: External API responses are cached 5 minutes in memory; geocode coordinates are cached in memory and on disk (`data/geocode_cache.json`)
3. **Three data formats**: JSON (API ↔ server ↔ browser), Go structs (server processing), HTML (display)
4. **Error handling at every step**: Network failure, bad input, missing data — all covered
5. **Async geocoding**: Map coordinates are NOT fetched during the page render — the browser fetches them from `/api/artist-geo?id=X` after the page loads, keeping "View Details" fast

---

## What's Next?

Now that you can trace data through the app, [Lesson 03](03-patterns.md) shows the recurring design patterns — the reusable strategies that make this code organized and maintainable.
