# Lesson 03: Design Patterns

## Why Learn Patterns?

A pattern is a **proven solution to a common problem**. Like cooking techniques: you don't reinvent "sauteing" every time — you learn the technique once and apply it to different ingredients.

Learning patterns helps you:

1. **Read code faster** — "Oh, this is the handler pattern"
2. **Write better code** — Use battle-tested solutions
3. **Communicate clearly** — "Let's use the MVC pattern"
4. **Understand WHY** — Not just WHAT the code does, but why it's organized this way

---

## Pattern 1: MVC (Model-View-Controller)

### The Problem (Plain English)
Without organization, you end up with "spaghetti code" — one giant file that mixes data definitions, business logic, HTML generation, and API calls. Finding anything becomes a nightmare.

### The Solution (The Pattern)
**Separate concerns into three layers:**
- **Model** = What the data looks like (the blueprint)
- **View** = What the user sees (the presentation)
- **Controller** = What happens when a user does something (the logic)

### The Code

```
groupie-tracker/
├── models/          ← MODEL: Data definitions
│   └── artist.go    ← "An Artist has a Name, Members, etc."
│
├── templates/       ← VIEW: What the user sees
│   ├── home.html    ← "Show artists in a grid"
│   └── artist.html  ← "Show one artist's profile"
│
├── handlers/        ← CONTROLLER: Logic that connects them
│   ├── home.go      ← "Get artists, render home page"
│   └── artist.go    ← "Get one artist, calculate stats, render profile"
│
└── services/        ← (bonus layer: data access)
    └── api.go       ← "Fetch data from external API"
```

### Real-World Analogy
- **Bad:** A restaurant where the chef takes orders, cooks food, AND serves tables (everything in one place)
- **Good:** A restaurant with a **host** (controller), **kitchen** (model/service), and **waitstaff** (view) — each role is specialized

### Code Location
- Models: `models/artist.go`
- Views: `templates/*.html`
- Controllers: `handlers/*.go`
- Services: `services/api.go`

---

## Pattern 2: Error-First Returns

### The Problem (Plain English)
Operations can fail — APIs can be down, files can be missing, data can be corrupted. If you ignore failures, your app either crashes or shows garbage data.

### The Solution (The Pattern)
Every function that can fail returns **two values**: the result AND an error. The caller MUST check the error before using the result.

### The Code

```go
// WRONG way: ignoring errors (NEVER do this)
artists, _ := services.FetchArtists()  // "_" ignores the error
// If the API is down, artists is nil
// Next line CRASHES: nil pointer dereference
fmt.Println(artists[0].Name)

// RIGHT way: always check errors
artists, err := services.FetchArtists()
if err != nil {
    // Handle the failure gracefully
    log.Printf("Error: %v", err)
    ErrorHandler(w, http.StatusInternalServerError)
    return  // Stop here — don't use broken data
}
// Now safe to use artists
fmt.Println(artists[0].Name)
```

### Real-World Analogy
- **Bad:** Ordering food delivery and eating whatever arrives without checking (could be wrong order, spoiled, or missing)
- **Good:** Ordering food delivery and **checking the bag** before eating — right items? Right temperature? Nothing missing?

### Code Location
- `handlers/home.go:15-19` — error check after FetchArtists
- `handlers/artist.go:30-34` — error check after FetchArtists
- `handlers/artist.go:48-52` — error check after FetchRelation
- `services/api.go:14-28` — error wrapping with `fmt.Errorf`

---

## Pattern 3: Early Return (Guard Clauses)

### The Problem (Plain English)
Without early returns, you get deeply nested `if/else` blocks — code that drifts further and further to the right, making it hard to follow.

### The Solution (The Pattern)
Check for problems FIRST and `return` immediately. The "happy path" (success case) stays at the lowest indentation level, making it easy to read.

### The Code

```go
// WRONG way: deeply nested conditions
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/artist/")
    id, err := strconv.Atoi(idStr)
    if err == nil {
        artists, err := services.FetchArtists()
        if err == nil {
            for _, a := range artists {
                if a.ID == id {
                    relation, err := services.FetchRelation(id)
                    if err == nil {
                        // Finally! The actual logic, buried 4 levels deep
                    }
                }
            }
        }
    }
}

// RIGHT way: early returns (what this codebase uses)
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/artist/")

    id, err := strconv.Atoi(idStr)
    if err != nil {
        ErrorHandler(w, http.StatusBadRequest, "Invalid artist ID")
        return  // ← BAIL OUT early
    }

    artists, err := services.FetchArtists()
    if err != nil {
        ErrorHandler(w, http.StatusInternalServerError)
        return  // ← BAIL OUT early
    }

    // ... find artist, if not found, return early ...

    // Happy path: all checks passed, do the real work
    tmpl.Execute(w, data)
}
```

### Real-World Analogy
- **Bad:** A bouncer who lets everyone into the club, then checks IDs inside, then checks the dress code inside, then checks the guest list inside...
- **Good:** A bouncer who checks ID at the door → wrong? **Go home.** Dress code? **Go home.** Guest list? **Go home.** Only people who pass ALL checks get in.

### Code Location
- `handlers/home.go:11-14` — early return for wrong path
- `handlers/artist.go:25-28` — early return for invalid ID
- `handlers/artist.go:30-34` — early return for fetch failure

---

## Pattern 4: Centralized Error Handler

### The Problem (Plain English)
If every handler writes its own error HTML, you get inconsistent error pages and duplicated code. Change the error page design? You'd have to update it everywhere.

### The Solution (The Pattern)
One function handles ALL errors. Every handler delegates to it. Change the error page once → it updates everywhere.

### The Code

```go
// ONE function handles ALL error rendering
func ErrorHandler(w http.ResponseWriter, status int, customMsg ...string) {
    w.WriteHeader(status)

    // Use custom message or default HTTP status text
    message := http.StatusText(status)
    if len(customMsg) > 0 {
        message = customMsg[0]
    }

    // Render the error template
    tmpl, err := template.ParseFiles("templates/error.html")
    if err != nil {
        // Last resort: plain text error
        http.Error(w, message, status)
        return
    }

    tmpl.Execute(w, map[string]interface{}{
        "Status":  status,
        "Message": message,
    })
}

// Usage (from any handler):
ErrorHandler(w, 404)                              // "Not Found"
ErrorHandler(w, 404, "Artist not found")          // Custom message
ErrorHandler(w, 500, "Unable to load artists...")  // Server error
```

### Real-World Analogy
- **Bad:** Every employee writes their own apology letters when things go wrong (inconsistent, sometimes unprofessional)
- **Good:** The company has a **customer service department** — all complaints go there, and they handle responses consistently

### Code Location
- `handlers/error.go:10-25` — the centralized error handler
- Called from every other handler when something goes wrong

---

## Pattern 5: Separation of Data Fetching and Presentation

### The Problem (Plain English)
If your handler directly makes HTTP calls, parses JSON, AND renders HTML, it becomes a massive, untestable function. You can't test the API logic without also dealing with templates.

### The Solution (The Pattern)
Put data fetching in a **service layer** (`services/`). Handlers call services, then pass the data to templates. Each layer has one job.

### The Code

```go
// services/api.go — ONLY responsible for fetching data
func FetchArtists() ([]models.Artist, error) {
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get(baseURL + "/artists")
    // ... decode JSON, return structs
}

// handlers/home.go — ONLY responsible for handling requests
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    artists, err := services.FetchArtists()  // Delegate data fetching
    // ... render template with the data
}
```

### Real-World Analogy
- **Bad:** A doctor who also runs the blood lab, the X-ray machine, and the pharmacy (too many jobs)
- **Good:** A doctor who **orders tests** from the lab, **reads results**, and **prescribes treatment** — the lab is a separate department

### Code Location
- Service layer: `services/api.go`
- Handlers that use it: `handlers/home.go`, `handlers/artist.go`, `handlers/search.go`

---

## Pattern 6: Filter-Then-Sort Pipeline

### The Problem (Plain English)
When users search and filter, you need to apply multiple criteria. Doing it all in one giant loop is confusing and hard to modify.

### The Solution (The Pattern)
Process data in stages: **fetch → filter → sort → respond**. Each stage is a separate, focused function.

### The Code

```go
// handlers/search.go — the pipeline

// Stage 1: FETCH all data
artists, err := services.FetchArtists()

// Stage 2: FILTER — keep only matching artists
var filtered []models.Artist
for _, artist := range artists {
    if matchesSearch(artist, query) && matchesYearFilter(artist, minYear, maxYear) {
        filtered = append(filtered, artist)
    }
}

// Stage 3: SORT — order the results
sortArtists(filtered, sortBy)

// Stage 4: RESPOND — send results as JSON
json.NewEncoder(w).Encode(filtered)
```

Each filter function has a single responsibility:

```go
// "Does this artist match the search query?"
func matchesSearch(artist models.Artist, query string) bool { ... }

// "Is this artist within the year range?"
func matchesYearFilter(artist models.Artist, minYear, maxYear string) bool { ... }

// "Sort these artists by the chosen criteria"
func sortArtists(artists []models.Artist, sortBy string) { ... }
```

### Real-World Analogy
- **Bad:** Dumping all your clothes in a pile and trying to find a specific red shirt in size M
- **Good:** An assembly line: **Color filter** (keep only red) → **Size filter** (keep only M) → **Sort** (by price) → **Pick the top result**

### Code Location
- Pipeline: `handlers/search.go:13-56`
- Search matcher: `handlers/search.go:58-71`
- Year filter: `handlers/search.go:73-86`
- Sorter: `handlers/search.go:88-100`

---

## Pattern 7: Variadic Functions (Optional Parameters)

### The Problem (Plain English)
The `ErrorHandler` sometimes needs a custom message and sometimes doesn't. In many languages, you'd need two separate functions or use `null`.

### The Solution (The Pattern)
Go's **variadic parameters** (`...string`) let a function accept zero OR more extra arguments.

### The Code

```go
// The "..." means: zero or more strings
func ErrorHandler(w http.ResponseWriter, status int, customMsg ...string) {
    message := http.StatusText(status)  // Default: "Not Found", "Internal Server Error"
    if len(customMsg) > 0 {
        message = customMsg[0]  // Override with custom message
    }
    // ...
}

// Usage:
ErrorHandler(w, 404)                        // Uses default: "Not Found"
ErrorHandler(w, 404, "Artist not found")    // Uses custom message
```

### Real-World Analogy
- **Bad:** Having two phone numbers — one for "call with a message" and one for "call without a message"
- **Good:** One phone number where you CAN leave a voicemail, but you don't HAVE to

### Code Location
- `handlers/error.go:10` — function signature with `customMsg ...string`

---

## Pattern 8: Client-Side Dynamic Rendering

### The Problem (Plain English)
When searching, reloading the entire page for every keystroke is slow and jarring. Users expect instant feedback.

### The Solution (The Pattern)
Use **AJAX** (Asynchronous JavaScript) to fetch data in the background, then update the page without a reload.

### The Code

```javascript
// static/script.js

// 1. Build the search URL from form inputs
async function applyFilters() {
    const query = document.getElementById('search').value;
    const url = `/api/search?q=${query}&minYear=${minYear}&maxYear=${maxYear}&sort=${sort}`;

    // 2. Fetch results WITHOUT reloading the page
    const response = await fetch(url);
    const data = await response.json();

    // 3. Rebuild the artist grid with new data
    displayArtists(data);
}

// 4. Create HTML cards from JSON data
function displayArtists(artists) {
    const grid = document.getElementById('artists-grid');
    grid.innerHTML = '';  // Clear old cards

    artists.forEach(artist => {
        // Build card HTML from data
        const card = document.createElement('div');
        card.className = 'artist-card';
        card.innerHTML = `<h2>${artist.name}</h2>...`;
        grid.appendChild(card);
    });
}
```

### Real-World Analogy
- **Bad:** Leaving the restaurant, driving to another one, and sitting down again every time you want to change your order
- **Good:** Calling the waiter over and saying "actually, I'd like the salad instead" — you stay seated, only the food changes

### Code Location
- JavaScript: `static/script.js` — `applyFilters()`, `displayArtists()`
- Go endpoint: `handlers/search.go` — `SearchHandler` returns JSON

---

## Pattern Summary

| Pattern | What It Solves | Where Used |
|---------|---------------|------------|
| MVC | Code organization | Project structure |
| Error-First Returns | Crash prevention | Every function that can fail |
| Early Return | Deep nesting | All handlers |
| Centralized Error Handler | Consistent errors | `handlers/error.go` |
| Service Layer Separation | Testability | `services/` vs `handlers/` |
| Filter-Then-Sort Pipeline | Complex queries | `handlers/search.go` |
| Variadic Functions | Optional parameters | `ErrorHandler` |
| Client-Side Rendering | Fast search | `static/script.js` |

---

## What's Next?

Now that you understand the patterns, [Lesson 04](04-line-by-line.md) walks through every important function line-by-line, explaining the WHY behind each decision.
