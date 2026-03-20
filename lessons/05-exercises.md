# Lesson 05: Exercises

Practice is how you move from "I read about it" to "I can do it." These exercises are ordered from absolute beginner to advanced. **Do them WITHOUT AI help** — that's the whole point.

---

## Level 0: Warmup (Absolute Beginners)

> **Skills:** Basic file navigation, reading code structure

### Exercise 0.1: Find a File
**Task:** List all the `.go` files in the `handlers/` directory.
<details>
<summary>Hint</summary>
Use your file explorer or terminal: <code>ls handlers/*.go</code>
</details>
<details>
<summary>Answer</summary>
<code>home.go</code>, <code>artist.go</code>, <code>artist_geo.go</code>, <code>search.go</code>, <code>locations.go</code>, <code>suggestions.go</code>, <code>error.go</code>, <code>search_test.go</code>, <code>suggestions_test.go</code>
</details>

### Exercise 0.2: Read a Struct
**Task:** Open `models/artist.go`. What fields does the `Artist` struct have? List them with their types.
<details>
<summary>Answer</summary>

| Field | Type |
|-------|------|
| ID | int |
| Name | string |
| Image | string |
| Members | []string |
| CreationDate | int |
| FirstAlbum | string |
</details>

### Exercise 0.3: Count Functions
**Task:** How many functions are defined in `handlers/search.go`? List their names.
<details>
<summary>Hint</summary>
Look for lines that start with <code>func</code>.
</details>
<details>
<summary>Answer</summary>
4 functions: <code>SearchHandler</code>, <code>matchesSearch</code>, <code>matchesYearFilter</code>, <code>sortArtists</code>
</details>

### Exercise 0.4: Understand a Comment
**Task:** Find a comment in `services/api.go` that explains WHY something is done (not just WHAT). What line is it on?
<details>
<summary>Hint</summary>
Look for comments that explain a reason, not just describe the next line.
</details>

### Exercise 0.5: Trace an Import
**Task:** In `handlers/home.go`, find one import from Go's standard library and one from this project. What are they?
<details>
<summary>Answer</summary>
Standard library: <code>"net/http"</code> (or <code>"html/template"</code> or <code>"log"</code>)<br>
This project: <code>"groupie-tracker/services"</code>
</details>

---

## Level 1: Find (Read the Code)

> **Skills:** Locating specific code, understanding definitions

### Exercise 1.1: The Entry Point
**Task:** What port does the server listen on? Which file and line defines it?
<details>
<summary>Answer</summary>
Port 8080, defined in <code>main.go</code> in the <code>http.ListenAndServe(":8080", nil)</code> call.
</details>

### Exercise 1.2: The API URL
**Task:** What is the base URL for the external API? Where is it defined?
<details>
<summary>Answer</summary>
<code>https://groupietrackers.herokuapp.com/api</code>, defined as a constant in <code>services/api.go</code>.
</details>

### Exercise 1.3: Band Classification
**Task:** If an artist has exactly 3 members, what `BandType` do they get? What about 7 members?
<details>
<summary>Hint</summary>
Look at the <code>getBandType</code> function in <code>handlers/artist.go</code>.
</details>
<details>
<summary>Answer</summary>
3 members → "Trio", 7 members → "Band" (the default case)
</details>

### Exercise 1.4: Error Messages
**Task:** What error message does the user see when they visit `/artist/banana`?
<details>
<summary>Answer</summary>
"Invalid artist ID" with status 400 (Bad Request). Found in <code>handlers/artist.go</code>.
</details>

### Exercise 1.5: Theme Persistence
**Task:** How does the app remember which theme (dark/light) the user chose? Where is this code?
<details>
<summary>Answer</summary>
It uses <code>localStorage.setItem('theme', ...)</code> in <code>static/script.js</code>. localStorage persists data in the browser even after page refresh.
</details>

---

## Level 2: Trace (Follow the Data)

> **Skills:** Data flow, mental execution, predicting behavior
> This is the MOST important skill for understanding code.

### Exercise 2.1: Trace a Home Page Load
**Task:** When a user visits `http://localhost:8080/`, trace the exact sequence of function calls that execute. List them in order.
<details>
<summary>Answer</summary>

1. `main.go` routes `/` to `handlers.HomeHandler`
2. `HomeHandler` checks `r.URL.Path != "/"` → false (it IS `/`)
3. `HomeHandler` calls `services.FetchArtists()`
4. `FetchArtists` makes HTTP GET to `https://groupietrackers.herokuapp.com/api/artists`
5. `FetchArtists` decodes JSON → `[]Artist`
6. `HomeHandler` calls `template.ParseFiles("templates/home.html")`
7. `HomeHandler` calls `tmpl.Execute(w, artists)`
8. Browser receives rendered HTML
</details>

### Exercise 2.2: Trace a Search
**Task:** A user types "mercury" in the search box and clicks Search. What happens step by step? Which functions run, and in what order?
<details>
<summary>Answer</summary>

1. JavaScript `applyFilters()` runs
2. Builds URL: `/api/search?q=mercury&minYear=&maxYear=&sort=default`
3. `fetch(url)` sends GET request to Go server
4. `SearchHandler` parses query params: `q="mercury"`
5. `FetchArtists()` fetches all 52 artists
6. Loop: for each artist, calls `matchesSearch(artist, "mercury")`
   - For "Queen": checks name "Queen" → no match, checks members "Freddie Mercury" → contains "mercury" → MATCH
   - For most others: no match
7. `matchesYearFilter` → returns true (no year filters set)
8. `sortArtists(filtered, "default")` → sort by ID
9. `json.Encode(filtered)` → send JSON to browser
10. JavaScript `displayArtists(data)` rebuilds the grid
</details>

### Exercise 2.3: Predict the Output
**Task:** What does `calculateTotalCountries` return for this data?
```go
DatesLocations: map[string][]string{
    "london-england":     {"01-01-2020"},
    "manchester-england": {"15-03-2020"},
    "paris-france":       {"20-06-2020"},
    "lyon-france":        {"25-06-2020"},
    "tokyo-japan":        {"10-09-2020"},
}
```
<details>
<summary>Hint</summary>
The function splits each key by "-" and takes the LAST part as the country, then counts unique countries.
</details>
<details>
<summary>Answer</summary>
3 countries: "england", "france", "japan". Even though there are 5 locations, only 3 unique countries.
</details>

### Exercise 2.4: What Goes Wrong?
**Task:** What happens if the external API at `groupietrackers.herokuapp.com` is completely down? Trace what the user sees when they visit the home page.
<details>
<summary>Answer</summary>

1. `HomeHandler` calls `FetchArtists()`
2. `FetchArtists` tries `client.Get(...)` → times out after 10 seconds
3. Returns `nil, error("failed to connect to API: ...")`
4. `HomeHandler` receives `err != nil`
5. Logs: "Error fetching artists: failed to connect to API: ..."
6. Calls `ErrorHandler(w, 500, "Unable to load artists. Please try again later.")`
7. User sees error page with "500: Unable to load artists. Please try again later."
</details>

### Exercise 2.5: Where Does the Data Transform?
**Task:** The external API sends `"creationDate": 1970`. By the time the user sees "56 Years Active" on the profile page, the data has been transformed. List every file where this number is touched or transformed.
<details>
<summary>Answer</summary>

1. `services/api.go` — JSON decoded from `"creationDate": 1970` into `Artist.CreationDate = 1970`
2. `handlers/artist.go` — `currentYear - artist.CreationDate` = `2026 - 1970` = `56`, stored in `ArtistPageData.YearsActive`
3. `templates/artist.html` — `{{.YearsActive}}` rendered as "56" in the HTML
</details>

---

## Level 3: Modify (Small Changes)

> **Skills:** Making targeted changes, understanding ripple effects, testing

### Exercise 3.1: Change the Port
**Task:** Change the server to run on port `3000` instead of `8080`. Which file(s) need to change?
<details>
<summary>Hint</summary>
Only ONE file needs to change. But also update the log message so it's not misleading.
</details>
<details>
<summary>Answer</summary>
Change <code>main.go</code>: replace <code>":8080"</code> with <code>":3000"</code> in <code>http.ListenAndServe</code>, and update the log message.
</details>

### Exercise 3.2: Add a Timeout Message
**Task:** When `FetchArtists()` fails, the error message says "Unable to load artists. Please try again later." Change it to also display the specific error (for debugging). What are the risks?
<details>
<summary>Hint</summary>
Think about security — should users see internal error details?
</details>
<details>
<summary>Answer</summary>
You could change the message to include the error, but this is a SECURITY RISK. Internal error messages might reveal server details (IP addresses, API endpoints, stack traces). The current approach is correct: log the detailed error (for developers) and show a vague message to users.
</details>

### Exercise 3.3: Add a New Sort Option
**Task:** Add a "members" sort option that sorts artists by the number of members (fewest first). Which files need to change?
<details>
<summary>Answer</summary>

1. `handlers/search.go` — add a new case in `sortArtists`:
   ```go
   case "members":
       sort.Slice(artists, func(i, j int) bool {
           return len(artists[i].Members) < len(artists[j].Members)
       })
   ```
2. `templates/home.html` — add `<option value="members">Members</option>` to the sort dropdown
3. `static/script.js` — no change needed (it already reads the select value dynamically)
</details>

### Exercise 3.4: Fix a Potential Bug
**Task:** In `matchesYearFilter`, what happens if someone passes `minYear=abc`? Is this a bug?
<details>
<summary>Answer</summary>
Not a bug! If <code>strconv.Atoi("abc")</code> fails, <code>err</code> is not nil, so the condition <code>err == nil && artist.CreationDate < min</code> is false (short-circuit evaluation). The function correctly ignores invalid year values and returns true. This is intentional defensive coding.
</details>

### Exercise 3.5: Add a Member Count Display
**Task:** On the home page, show the number of members below each artist's name (e.g., "4 members"). Which file(s) need to change, and what Go template syntax would you use?
<details>
<summary>Answer</summary>
Only `templates/home.html` needs to change. Add:
```html
<p>{{len .Members}} members</p>
```
`len` is a built-in Go template function that returns the length of a slice.
</details>

---

## Level 4: Extend (Add Features)

> **Skills:** Planning additions, maintaining consistency, following existing patterns

### Exercise 4.1: Add a "Random Artist" Button
**Task:** Plan (don't code yet) a feature that shows a random artist when a button is clicked. Write pseudo-code and list which files need to change.
<details>
<summary>Answer</summary>

**Pseudo-code:**
```
1. Add a "Random Artist" button to home.html
2. When clicked, JavaScript picks a random artist ID
3. Navigate to /artist/{randomID}
```

**Files to change:**
- `templates/home.html` — add the button
- `static/script.js` — add click handler that generates random ID and navigates

**OR** (server-side approach):
- `handlers/random.go` — new handler that picks a random artist and redirects
- `main.go` — register the new route
- `templates/home.html` — add link to `/random`
</details>

### Exercise 4.2: Add a "First Album Year" Filter
**Task:** Add a filter for the first album year (similar to the creation date filter). Think about: What's tricky about the `FirstAlbum` field?
<details>
<summary>Hint</summary>
<code>FirstAlbum</code> is a string like "14-12-1973" (DD-MM-YYYY), not an integer. You'd need to extract the year.
</details>
<details>
<summary>Answer</summary>

**The tricky part:** `FirstAlbum` is `"14-12-1973"` — you need to split by `"-"` and take the LAST part as the year, then convert to integer.

**Files to change:**
1. `handlers/search.go` — add `matchesAlbumYearFilter()` function, add params to SearchHandler
2. `templates/home.html` — add min/max album year inputs
3. `static/script.js` — include new params in `applyFilters()` URL
</details>

### Exercise 4.3: Add Concert Count to Home Page
**Task:** Show the total number of concerts on each artist's card on the home page. What's the challenge, and how would you solve it?
<details>
<summary>Answer</summary>

**Challenge:** The home page only fetches artist data, not relation data. Concert counts are in the relation endpoint.

**Options:**
1. **Fetch relations too** — modify `HomeHandler` to also call `FetchRelation` for each artist. Problem: 52 extra API calls = slow.
2. **Create a combined data type** — fetch all relations once, match them to artists, pass combined data to template.
3. **Client-side approach** — use AJAX to fetch concert data lazily when cards are visible.

Option 2 is best: fetch `/api/relation` once, build a map of `artistID → concertCount`, pass it alongside artists to the template.
</details>

### Exercise 4.4: Add Unit Tests
**Task:** Write a test for `getBandType`. What edge cases should you test?
<details>
<summary>Answer</summary>

```go
func TestGetBandType(t *testing.T) {
    tests := []struct {
        count    int
        expected string
    }{
        {0, "Band"},       // Edge case: zero members
        {1, "Solo Artist"},
        {2, "Duo"},
        {3, "Trio"},
        {4, "Quartet"},
        {5, "Quintet"},
        {6, "Band"},       // 6+ = generic "Band"
        {100, "Band"},     // Edge case: very large number
    }

    for _, tt := range tests {
        result := getBandType(tt.count)
        if result != tt.expected {
            t.Errorf("getBandType(%d) = %q, want %q", tt.count, result, tt.expected)
        }
    }
}
```

Edge cases: 0 members (unrealistic but possible), boundary values (5 vs 6), very large numbers.
</details>

---

## Level 5: Break & Fix (Debugging)

> **Skills:** Understanding WHY code exists, debugging, defensive thinking

### Exercise 5.1: Remove the Path Check
**Task:** What happens if you remove the `if r.URL.Path != "/"` check from `HomeHandler`? Test your prediction by visiting `/favicon.ico`.
<details>
<summary>Answer</summary>
Without the check, visiting `/favicon.ico` (which browsers automatically request) would show the home page instead of a 404. EVERY unknown URL would show the home page because `/` is a prefix of all paths. The check ensures only exact matches to `/` show the home page.
</details>

### Exercise 5.2: Break the JSON Tags
**Task:** What happens if you change `json:"name"` to `json:"artist_name"` in the Artist struct? Which features break, and why?
<details>
<summary>Answer</summary>
**Everything breaks.** The external API sends `"name": "Queen"`, but now Go looks for `"artist_name"` in the JSON. Since there's no `"artist_name"` field in the API response, the `Name` field stays empty (`""`). All artist names would appear blank on the home page and artist pages. The search would also break because it searches `artist.Name` which is now always empty.
</details>

### Exercise 5.3: Remove defer
**Task:** What happens if you remove `defer resp.Body.Close()` from `FetchArtists`? Will the app crash immediately?
<details>
<summary>Answer</summary>
It won't crash immediately. But each request leaks a network connection (the HTTP response body stays open). Over time, you'll run out of available connections, and new API calls will fail with "too many open files" or similar errors. `defer` prevents this resource leak.
</details>

### Exercise 5.4: Simulate an API Failure
**Task:** Change the `baseURL` in `services/api.go` to `"https://doesnotexist.invalid/api"`. What error does the user see? Is it user-friendly?
<details>
<summary>Answer</summary>
The user sees "500: Unable to load artists. Please try again later." — which IS user-friendly. The actual DNS error is logged to the console but not shown to the user. This is correct behavior: vague messages for users, detailed logs for developers.
</details>

### Exercise 5.5: Concurrent Map Bug
**Task:** (Advanced) If two users search at the same time, could there be a race condition in `sortArtists`? Why or why not?
<details>
<summary>Answer</summary>
No race condition. Each call to `SearchHandler` creates its OWN `filtered` slice with `var filtered []models.Artist` and `append`. These are local variables, not shared state. Each request works on its own copy of the data. If they were modifying a global/shared slice, THEN there would be a race condition.
</details>

---

## The Vibecoding Graduation Test

Can you do these **WITHOUT AI assistance?**

### 1. Pseudo-code First
Write pseudo-code for a "favorites" feature where users can bookmark artists.

### 2. Predict Files
Which files would implementing "favorites" touch? List them and explain why.

### 3. Explain the Why
Why does `FetchRelation` download ALL relations and then search for the right one, instead of fetching just one?

### 4. Find the Edge Case
What happens if the external API returns an artist with an empty `Members` slice (`[]`)? Trace through `getBandType` and the home page template.

### 5. Rubber Duck
Explain the `matchesSearch` function out loud (to a rubber duck, a pet, or an imaginary friend). Cover:
- What it receives
- What it returns
- The three cases it checks
- Why it converts strings to lowercase

---

**If you can do all five, you've graduated from vibecoding to understanding.**

You don't need to memorize every line. You need to know:
- WHERE to look
- HOW to trace data flow
- WHY the code is written this way

That's the difference between copying code and understanding code.

---

## What's Next?

[Lesson 06](06-gotchas.md) covers the common mistakes and tricky spots in this codebase — the things that are obvious AFTER you know them.
