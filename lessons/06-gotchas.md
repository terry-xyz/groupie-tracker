# Lesson 06: Gotchas & Common Mistakes

## Why This Chapter Exists

Every codebase has "gotchas" — things that are obvious AFTER you know them but can trip you up badly the first time. Each gotcha here represents a **real bug someone could introduce**. Understanding them makes you a better debugger and a more careful programmer.

---

## Gotcha 1: The "/" Route Matches Everything

### The Bug
Go's default HTTP router uses **prefix matching**. The pattern `"/"` is a prefix of every URL. Without a guard, HomeHandler would handle `/artist/5`, `/api/search`, `/favicon.ico`, and literally every other path.

```go
// WRONG: No path check — every URL shows the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    artists, _ := services.FetchArtists()
    tmpl, _ := template.ParseFiles("templates/home.html")
    tmpl.Execute(w, artists)
}

// RIGHT: Guard clause ensures only exact "/" matches
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        ErrorHandler(w, http.StatusNotFound)
        return
    }
    // ... rest of handler
}
```

### How to Spot This Bug
- Visiting any random URL (like `/asdfgh`) shows the home page instead of a 404
- Artist pages or search work, but the URL bar shows `/artist/5` while the home page content appears
- The browser's favicon request (`/favicon.ico`) triggers a full API fetch

### Code Location
`handlers/home.go:11-14`

---

## Gotcha 2: Ignoring Errors Causes Nil Pointer Crashes

### The Bug
If you ignore an error from `FetchArtists()` using `_`, the `artists` variable is `nil`. Trying to loop over `nil` won't crash (it just doesn't iterate), but trying to access `artists[0]` WILL crash with a **nil pointer dereference** — one of Go's most common runtime panics.

```go
// WRONG: Ignoring the error
artists, _ := services.FetchArtists()
fmt.Println(artists[0].Name)  // CRASH if API is down (artists is nil)

// RIGHT: Always check errors
artists, err := services.FetchArtists()
if err != nil {
    ErrorHandler(w, http.StatusInternalServerError, "Unable to load artists")
    return
}
fmt.Println(artists[0].Name)  // Safe — we know artists is not nil
```

### How to Spot This Bug
- Error message: `runtime error: index out of range` or `nil pointer dereference`
- The app crashes randomly (only when the API is slow or down)
- Works fine in development, crashes in production (because the API is less reliable at scale)

### Code Location
Every handler that calls a service function checks for errors.

---

## Gotcha 3: Forgetting `defer resp.Body.Close()`

### The Bug
Every HTTP response has a "body" (the data). If you don't close it, the network connection stays open forever. After enough leaked connections, your app can't make new requests.

```go
// WRONG: Connection leaks
resp, err := client.Get(url)
if err != nil {
    return nil, err
}
// No Close() — connection stays open forever!
var data MyType
json.NewDecoder(resp.Body).Decode(&data)
return data, nil

// RIGHT: defer ensures Close() runs even if the function crashes
resp, err := client.Get(url)
if err != nil {
    return nil, err
}
defer resp.Body.Close()  // Will run when the function returns, no matter what
```

### How to Spot This Bug
- App works fine at first, then starts failing after many requests
- Error messages like "too many open files" or "connection refused"
- Memory usage slowly increases over time

### Code Location
`services/api.go:17` and `services/api.go:39` — both use `defer resp.Body.Close()`

---

## Gotcha 4: JSON Tag Mismatch

### The Bug
If the `json:"..."` tag in a struct doesn't match the key in the JSON data, Go silently ignores the field (sets it to zero/empty). No error, no warning — just missing data.

```go
// The API sends: {"name": "Queen"}

// WRONG: Tag doesn't match API
type Artist struct {
    Name string `json:"artist_name"`  // Looks for "artist_name" — won't find "name"
}
// Result: Artist.Name is "" (empty string) — NO ERROR

// RIGHT: Tag matches the API's JSON key
type Artist struct {
    Name string `json:"name"`  // Matches "name" in JSON
}
// Result: Artist.Name is "Queen"
```

### How to Spot This Bug
- Artist names, images, or other fields appear blank
- No error messages anywhere (this is a silent failure)
- Data works in some fields but not others

### Code Location
`models/artist.go` — all struct tags must match the API's JSON keys exactly

---

## Gotcha 5: The Country Extraction Assumption

### The Bug
`calculateTotalCountries` splits location strings by `"-"` and takes the **last part** as the country. This works for `"london-england"` and `"new_york-usa"`, but could break for countries with hyphens in their name.

```go
// Current implementation
parts := strings.Split(location, "-")
country := parts[len(parts)-1]  // Last element

// Works for:
// "london-england"       → "england"    ✓
// "new_york-usa"         → "usa"        ✓
// "sao_paulo-brazil"     → "brazil"     ✓

// Would break for (hypothetical):
// "port-au-prince-haiti" → "haiti"      ✓ (actually works!)
// But what about edge cases in the actual API data?
```

### How to Spot This Bug
- Country counts seem wrong on artist profile pages
- Certain locations show unexpected country names

### Code Location
`handlers/artist.go` — `calculateTotalCountries` function

---

## Gotcha 6: Template Parse Errors Are Silent

### The Bug
If `template.ParseFiles()` succeeds but `tmpl.Execute()` fails (e.g., because the template references a field that doesn't exist), the error is silently discarded in this codebase. Partial HTML may have already been written to the response.

```go
// The current code doesn't check Execute's error
tmpl.Execute(w, data)

// More defensive approach:
if err := tmpl.Execute(w, data); err != nil {
    log.Printf("Template execution error: %v", err)
    // But at this point, partial HTML may already be sent to the browser
    // This is a fundamental challenge with streaming templates
}
```

### How to Spot This Bug
- Pages render partially — some data shows, some doesn't
- No error page appears, but the page looks "broken"
- Console shows no error (unless you add the check)

### Code Location
`handlers/home.go:29`, `handlers/artist.go:71` — `tmpl.Execute()` calls

---

## Gotcha 7: No Caching — Every Request Hits the API

### The Bug (More of a Design Choice)
Every page load calls `FetchArtists()`, which makes a fresh HTTP request to the external API. If 100 users load the home page, that's 100 API calls for the same data. This works fine for a small project but would be a performance problem at scale.

```go
// Current: fresh API call every time
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    artists, err := services.FetchArtists()  // Hits external API
    // ...
}

// With caching (not implemented, but would look like):
var cachedArtists []models.Artist
var cacheTime time.Time

func FetchArtistsCached() ([]models.Artist, error) {
    if time.Since(cacheTime) < 5*time.Minute && cachedArtists != nil {
        return cachedArtists, nil  // Return cached data
    }
    // ... fetch fresh data and update cache
}
```

### How to Spot This Issue
- Pages load slowly (waiting for external API on every request)
- External API rate-limits your app
- App breaks completely when the external API is down (no fallback)

### Code Location
`services/api.go` — both `FetchArtists` and `FetchRelation` make fresh calls every time

---

## Gotcha 8: `strconv.Atoi` and Invalid Input

### The Bug
When extracting the artist ID from `/artist/{id}`, if someone visits `/artist/` (empty ID) or `/artist/3.5` (float), `strconv.Atoi` returns an error. The code handles this correctly, but forgetting the check would crash the app.

```go
// URL: /artist/         → idStr = ""        → Atoi fails
// URL: /artist/3.5      → idStr = "3.5"     → Atoi fails
// URL: /artist/-1       → idStr = "-1"      → Atoi succeeds! ID = -1
// URL: /artist/9999999  → idStr = "9999999" → Atoi succeeds! But no artist with that ID

id, err := strconv.Atoi(idStr)
if err != nil {
    ErrorHandler(w, http.StatusBadRequest, "Invalid artist ID")
    return
}
// Note: negative IDs and very large IDs pass this check
// They're caught later when searching for the artist (not found → 404)
```

### How to Spot This Bug
- Visiting `/artist/` shows an error page (which is correct!)
- But visiting `/artist/-1` might not show the error you expect

### Code Location
`handlers/artist.go:23-27`

---

## Gotcha 9: Template Re-parsing on Every Request

### The Bug (Performance)
Every handler call parses the HTML template from disk:

```go
// Current: parse template on EVERY request
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    // ...
    tmpl, err := template.ParseFiles("templates/home.html")  // Reads file from disk
    // ...
}
```

This reads and parses the file from disk on every single request. For a small project this is fine, and it has the benefit of reflecting template changes without restarting the server. But at scale, you'd want to parse templates once at startup:

```go
// Better for production: parse once, reuse many times
var homeTemplate = template.Must(template.ParseFiles("templates/home.html"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    // ...
    homeTemplate.Execute(w, artists)  // No file read — template already parsed
}
```

### How to Spot This Issue
- Slow page loads under high traffic
- High disk I/O on the server

---

## Gotcha 10: XSS Risk in JavaScript Template Literals

### The Bug
In `script.js`, `displayArtists()` inserts data directly into HTML using template literals:

```javascript
// Current code
card.innerHTML = `
    <img src="${artist.image}" alt="${artist.name}">
    <h2>${artist.name}</h2>
`;
```

If an artist's name contained HTML like `<script>alert('hacked')</script>`, it would execute as code. This is called **Cross-Site Scripting (XSS)**.

**In this project, the risk is LOW** because the data comes from a trusted external API (not user input). But if you ever let users submit data that gets displayed this way, you'd need to sanitize it:

```javascript
// Safer approach: use textContent instead of innerHTML for text
const h2 = document.createElement('h2');
h2.textContent = artist.name;  // Treated as text, not HTML
card.appendChild(h2);
```

### How to Spot This Bug
- Strange characters appearing in artist names
- Browser developer console showing script errors
- Any user-submitted data appearing in `innerHTML`

### Code Location
`static/script.js` — `displayArtists()` function

---

## The Gotcha Mindset

When reading or writing code, always ask:

1. **"What could go wrong?"** — API down, bad input, missing data, network timeout
2. **"What order do things happen?"** — Does the check come BEFORE the action?
3. **"What does the data actually look like?"** — Is `"creationDate"` an int or string? Is `"firstAlbum"` a date or a formatted string?
4. **"What would happen if...?"**
   - ...the list is empty?
   - ...the string is blank?
   - ...the network is slow?
   - ...two requests arrive at the same time?

This **defensive thinking** is what separates robust code from fragile code. You don't need to handle every edge case — but you should be **aware** of them.

---

## What's Next?

[Lesson 07](07-glossary.md) is your reference guide — every term, abbreviation, and concept decoded in one place.
