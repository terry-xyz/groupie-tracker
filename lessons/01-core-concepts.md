# Lesson 01: Core Concepts

## Reading Go Code: Quick Guide

Before diving into concepts, here's a cheat sheet for reading Go syntax:

```go
package main                    // This file belongs to the "main" package

import "fmt"                    // Load a library (like importing a tool)

var name string = "Alice"       // Declare a variable with explicit type
age := 25                       // Short declaration (Go guesses the type)

func greet(name string) string { // Function: takes a string, returns a string
    return "Hello, " + name
}

if age > 18 {                   // Condition: runs code if true
    fmt.Println("Adult")
}

for i := 0; i < 5; i++ {       // Loop: repeat 5 times (i goes 0,1,2,3,4)
    fmt.Println(i)
}

for _, item := range items {   // Loop over a list ("_" means "ignore index")
    fmt.Println(item)
}

type Dog struct {               // Struct: a blueprint with named fields
    Name  string
    Breed string
}

myDog := Dog{Name: "Rex", Breed: "Lab"}  // Create an instance
fmt.Println(myDog.Name)                   // Access a field: "Rex"
```

### Common Symbols Explained

| Symbol | Meaning | Example |
|--------|---------|---------|
| `:=` | Declare and assign | `x := 5` (create x, set to 5) |
| `=` | Reassign | `x = 10` (change x to 10) |
| `!=` | Not equal | `if err != nil` (if error exists) |
| `_` | "I don't need this" | `for _, v := range list` (ignore index) |
| `...` | Variable arguments | `func f(msgs ...string)` (any number of strings) |
| `*` | Pointer (address of) | `*Relation` (a reference to a Relation) |
| `&` | Get address | `&myVar` (get pointer to myVar) |
| `[]` | Slice (list) | `[]string` (list of strings) |
| `map[K]V` | Dictionary | `map[string][]string` (key→list of strings) |

---

## Concept 1: Structs (Data Blueprints)

### What (Simple Definition)
> A struct is a **blueprint** that groups related pieces of data together — like a form with labeled fields.

### Why (Why This Matters)
- **The problem:** An artist has a name, image, members, and creation date. Without structs, you'd have separate variables scattered everywhere.
- **Real-world analogy:** Think of a **driver's license**. It's one card with multiple fields: name, photo, date of birth, address. A struct is the same idea in code.

### Where (Files)
- `models/artist.go:3-9` — The `Artist` struct
- `models/artist.go:11-14` — The `Relation` struct
- `handlers/artist.go:14-21` — The `ArtistPageData` struct (no `Locations` field — geocoding is async)

### How (Code Walkthrough)

```go
// models/artist.go

// This is the "driver's license" for an artist.
// Every artist in the system has exactly these fields.
type Artist struct {
    ID           int      `json:"id"`           // Unique number (like a social security number)
    Name         string   `json:"name"`         // Band/artist name
    Image        string   `json:"image"`        // URL to their photo
    Members      []string `json:"members"`      // List of member names
    CreationDate int      `json:"creationDate"` // Year they formed (e.g., 1985)
    FirstAlbum   string   `json:"firstAlbum"`   // Date of first album release
}
```

**What are those backtick tags?** (`json:"id"`)
These are **struct tags** — instructions for Go's JSON decoder. They say: "When you see `"id"` in JSON data, put that value into the `ID` field." It's like a translation guide between the API's language and our code's language.

```go
// This is the concert history for an artist.
// The key is a location, the value is a list of dates they played there.
type Relation struct {
    ID             int                 `json:"id"`
    DatesLocations map[string][]string `json:"datesLocations"`
    // Example: {"new_york-usa": ["01-05-2020", "15-08-2021"]}
}
```

**Reading `map[string][]string`:**
- `map` = dictionary/lookup table
- `[string]` = keys are strings (location names)
- `[]string` = values are lists of strings (concert dates)

---

## Concept 2: HTTP Handlers (Request Responders)

### What (Simple Definition)
> A handler is a **function that runs when someone visits a specific URL** — like a receptionist who answers different phone extensions.

### Why (Why This Matters)
- **The problem:** When a user visits `/artist/3`, the server needs to know WHAT to do. Handlers are the "what to do" instructions.
- **Real-world analogy:** A **restaurant host**. When you walk in and say "table for two," the host knows the procedure: check availability, pick a table, seat you. A handler does the same for web requests.

### Where (Files)
- `handlers/home.go:10-29` — HomeHandler
- `handlers/artist.go:21-73` — ArtistHandler
- `handlers/search.go:13-56` — SearchHandler
- `handlers/error.go:10-25` — ErrorHandler

### How (Code Walkthrough)

```go
// handlers/home.go

// This function runs every time someone visits the home page (/).
// "w" is the response (what we send back)
// "r" is the request (what the user asked for)
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    // STEP 1: Make sure they're visiting exactly "/", not "/random"
    if r.URL.Path != "/" {
        ErrorHandler(w, http.StatusNotFound)
        return  // Stop here — don't continue
    }

    // STEP 2: Get all artists (served from 5-min in-memory cache after first load)
    artists, err := services.GetArtists()
    if err != nil {
        // Something went wrong → show error page
        log.Printf("Error fetching artists: %v", err)
        ErrorHandler(w, http.StatusInternalServerError, "Unable to load artists...")
        return
    }

    // STEP 3: Render with the pre-parsed template (parsed once at startup via sync.Once)
    getHomeTmpl().Execute(w, artists)
}
```

**Key insight:** Every handler follows the same pattern:
1. Validate the request
2. Get the data
3. Load the template
4. Send the response

If anything goes wrong at any step, show an error and **return immediately** (don't keep going).

---

## Concept 3: API Fetching (Getting External Data)

### What (Simple Definition)
> API fetching means **asking another server for data over the internet** — like ordering food for delivery instead of cooking it yourself.

### Why (Why This Matters)
- **The problem:** We don't have a database of artists. Someone else (groupietrackers.herokuapp.com) maintains that data. We need to ask them for it.
- **Real-world analogy:** You don't grow your own vegetables. You go to the **grocery store** (the API), pick what you need (make a request), and bring it home (parse the response).

### Where (Files)
- `services/api.go:10-30` — `FetchArtists()`
- `services/api.go:32-56` — `FetchRelation()`

### How (Code Walkthrough)

```go
// services/api.go

const baseURL = "https://groupietrackers.herokuapp.com/api"

func FetchArtists() ([]models.Artist, error) {
    // STEP 1: Create an HTTP client with a 10-second timeout
    //         (Don't wait forever if the API is down)
    client := &http.Client{Timeout: 10 * time.Second}

    // STEP 2: Ask the API for all artists
    resp, err := client.Get(baseURL + "/artists")
    if err != nil {
        return nil, fmt.Errorf("failed to connect to API: %v", err)
    }
    defer resp.Body.Close()  // Always clean up when done

    // STEP 3: Check if the API said "OK" (status 200)
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
    }

    // STEP 4: Convert JSON text into Go structs
    var artists []models.Artist
    if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
        return nil, fmt.Errorf("failed to parse artist data: %v", err)
    }

    return artists, nil
}
```

**Key insight:** The function returns TWO values: `([]models.Artist, error)`. This is Go's way of saying "here's the data, OR here's what went wrong." The caller MUST check the error before using the data.

**What's `defer`?** It means "do this LATER, right before the function ends." `defer resp.Body.Close()` ensures we close the network connection even if something crashes. Think of it like "remind me to lock the door when I leave."

---

## Concept 4: Templates (Dynamic HTML)

### What (Simple Definition)
> Templates are **HTML files with blanks to fill in** — like a Mad Libs where Go inserts the real data.

### Why (Why This Matters)
- **The problem:** Every artist has a different name, image, and members. We can't write 52 separate HTML pages. Templates let us write ONE page that works for ALL artists.
- **Real-world analogy:** A **mail merge**. You write one letter template: "Dear {{.Name}}, thank you for..." and the computer fills in each person's name.

### Where (Files)
- `templates/home.html` — Home page template
- `templates/artist.html` — Artist detail template
- `templates/error.html` — Error page template

### How (Code Walkthrough)

```html
<!-- templates/home.html (simplified) -->

<!-- Loop over every artist in the data -->
{{range .}}
<div class="artist-card">
    <!-- Insert the artist's image URL -->
    <img src="{{.Image}}" alt="{{.Name}}">

    <!-- Insert the artist's name -->
    <h2>{{.Name}}</h2>

    <!-- Insert member names, joined with commas -->
    <p>{{range $i, $m := .Members}}{{if $i}}, {{end}}{{$m}}{{end}}</p>

    <!-- Create a link using the artist's ID -->
    <a href="/artist/{{.ID}}">View Details</a>
</div>
{{end}}
```

**Template syntax cheat sheet:**

| Syntax | Meaning | Example |
|--------|---------|---------|
| `{{.Name}}` | Insert the Name field | `"Queen"` |
| `{{range .}}` | Loop over a list | Repeat for each artist |
| `{{end}}` | End a block (range, if) | — |
| `{{if .}}` | If truthy | Show only if data exists |
| `{{$i}}` | Loop variable | Index number |

---

## Concept 5: JSON Encoding/Decoding (Data Translation)

### What (Simple Definition)
> JSON is a **universal data format** that both humans and computers can read — like a standardized shipping label that any delivery company understands.

### Why (Why This Matters)
- **The problem:** The external API sends data as text. Our Go code needs Go structs. The browser's JavaScript needs JavaScript objects. JSON is the common language they all speak.
- **Real-world analogy:** **Currency exchange.** The API "speaks" JSON. Go "speaks" structs. JavaScript "speaks" objects. JSON is the exchange booth that converts between them.

### Where (Files)
- `services/api.go:25` — Decoding JSON → Go structs (API response)
- `handlers/search.go:52-55` — Encoding Go structs → JSON (search results)

### How (Code Walkthrough)

```go
// DECODING: JSON text → Go structs (reading FROM the API)
var artists []models.Artist
json.NewDecoder(resp.Body).Decode(&artists)
// The JSON: [{"id":1,"name":"Queen","members":["Freddie Mercury",...]}]
// Becomes: []Artist{{ID:1, Name:"Queen", Members:["Freddie Mercury",...]}}

// ENCODING: Go structs → JSON text (sending TO the browser)
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(filteredArtists)
// The struct: []Artist{{ID:1, Name:"Queen",...}}
// Becomes: [{"id":1,"name":"Queen",...}]
```

**Key insight:** `Decode` reads JSON IN. `Encode` writes JSON OUT. The struct tags (`json:"name"`) tell Go how to map between JSON field names and Go field names.

---

## Concept 6: Error Handling (The Go Way)

### What (Simple Definition)
> In Go, every operation that might fail returns an error value — and you MUST check it. It's like a doctor who checks for allergies before every procedure.

### Why (Why This Matters)
- **The problem:** APIs can be down. URLs can be wrong. Templates can be missing. If you ignore these failures, your app crashes or shows garbage.
- **Real-world analogy:** **A pilot's checklist.** Before takeoff, pilots check every system. They don't assume the engines work — they verify. Go's error handling is the same: verify before proceeding.

### Where (Files)
- Every single handler and service function uses this pattern

### How (Code Walkthrough)

```go
// The pattern: try something, check if it failed, handle the failure

// Step 1: Try the operation (returns result AND error)
artists, err := services.FetchArtists()

// Step 2: Check if something went wrong
if err != nil {
    // Step 3a: Handle the error (log it, show error page)
    log.Printf("Error: %v", err)
    ErrorHandler(w, http.StatusInternalServerError)
    return  // STOP — don't use the broken data
}

// Step 3b: If no error, safely use the result
fmt.Println(artists[0].Name)  // Only runs if fetch succeeded
```

**Why `return` after errors?** Without `return`, the code would continue and try to use `artists` — which is `nil` (empty) because the fetch failed. This would crash the server. The `return` is a safety fence.

---

## Concept 7: Routing (URL → Function Mapping)

### What (Simple Definition)
> Routing maps URLs to handler functions — like a **phone system menu**: "Press 1 for sales, press 2 for support."

### Why (Why This Matters)
- **The problem:** When someone types `localhost:8080/artist/5`, the server needs to know which function to run. Routing is the map.
- **Real-world analogy:** A **building directory**. Floor 1 = reception, Floor 2 = marketing, Floor 3 = engineering. Routing tells requests which "floor" to go to.

### Where (Files)
- `main.go:11-14` — Route definitions

### How (Code Walkthrough)

```go
// main.go

// "When someone visits /, run HomeHandler"
http.HandleFunc("/", handlers.HomeHandler)

// "When someone visits /artist/anything, run ArtistHandler"
http.HandleFunc("/artist/", handlers.ArtistHandler)

// "When someone visits /api/search, run SearchHandler"
http.HandleFunc("/api/search", handlers.SearchHandler)

// "Serve files from the static/ folder when someone visits /static/"
fs := http.FileServer(http.Dir("static"))
http.Handle("/static/", http.StripPrefix("/static/", fs))
```

**Key insight:** Go's default router uses **prefix matching**. `/artist/` matches any URL starting with `/artist/` — so `/artist/1`, `/artist/42`, and `/artist/banana` all go to `ArtistHandler`. The handler must then extract and validate the ID itself.

**Why does HomeHandler check `r.URL.Path != "/"`?** Because `/` matches EVERYTHING (it's a prefix of every URL). Without this check, visiting `/random-page` would show the home page instead of a 404 error.

---

## Summary: How These Concepts Connect

```
User types URL
       │
       ▼
   ROUTING (main.go)
   "Which handler should respond?"
       │
       ▼
   HANDLER (handlers/*.go)
   "What should I do with this request?"
       │
       ├──→ FETCH DATA (services/api.go)
       │    "Get data from external API"
       │         │
       │         ▼
       │    JSON DECODE
       │    "Convert API response → Go STRUCTS"
       │         │
       │         ▼
       │    STRUCTS (models/artist.go)
       │    "Data now organized in Go's format"
       │
       ▼
   TEMPLATE (templates/*.html)
   "Fill in the HTML blanks with data"
       │
       ▼
   HTML RESPONSE → Browser
```

---

## What's Next?

Now that you understand the individual concepts, [Lesson 02](02-data-flow.md) shows how data flows through the entire application — from the external API all the way to what the user sees on screen.
