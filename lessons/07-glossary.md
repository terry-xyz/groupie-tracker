# Lesson 07: Glossary

## Start Here: Essential Terms

These terms appear constantly throughout the codebase. Master these first.

| Term | Plain English | Example in This Project |
|------|---------------|------------------------|
| **Handler** | A function that responds to a web request | `HomeHandler` serves the home page |
| **Route** | A URL pattern mapped to a handler | `"/artist/"` → `ArtistHandler` |
| **Struct** | A blueprint grouping related data fields | `Artist` has Name, Members, etc. |
| **Slice** | A resizable list | `[]string{"Alice", "Bob"}` |
| **Map** | A dictionary (key → value lookup) | `map[string][]string` |
| **Template** | HTML with placeholders for dynamic data | `{{.Name}}` becomes `"Queen"` |
| **JSON** | Text format for exchanging data | `{"name":"Queen","id":1}` |
| **API** | A way for programs to talk to each other | Our app talks to the groupietrackers API |
| **Error** | A value indicating something went wrong | `err != nil` means "failure" |
| **nil** | "Nothing" — empty/absent value | `var x *Artist = nil` (no artist) |

---

## Domain Terms (The Problem Space)

These terms describe the music domain this app operates in.

| Term | Meaning | Where Used |
|------|---------|-----------|
| **Artist** | A musician or band | The main entity in the app |
| **Member** | A person in a band | `Artist.Members` field |
| **Creation Date** | Year the band was formed | `Artist.CreationDate` (e.g., 1970) |
| **First Album** | Release date of debut album | `Artist.FirstAlbum` (e.g., "14-12-1973") |
| **Relation** | Concert history (where & when they played) | `Relation.DatesLocations` |
| **Band Type** | Classification by member count | Solo, Duo, Trio, Quartet, Quintet, Band |

---

## Code Terms (The Implementation)

### Go Language Basics

| Term | Plain English | Example |
|------|---------------|---------|
| **package** | A folder of related code | `package handlers` |
| **import** | Load code from another package | `import "net/http"` |
| **func** | Define a function | `func FetchArtists() {...}` |
| **type** | Define a new data type | `type Artist struct {...}` |
| **var** | Declare a variable | `var name string` |
| `:=` | Short variable declaration (type inferred) | `x := 5` |
| `defer` | "Run this when the function ends" | `defer resp.Body.Close()` |
| `range` | Iterate over a collection | `for _, a := range artists` |
| `_` (blank identifier) | "I don't need this value" | `for _, item := range list` |
| `&` (address-of) | Get a pointer to a value | `&artists[i]` |
| `*` (dereference) | The value a pointer points to | `*artist` (the actual Artist) |
| `...` (variadic) | Accept zero or more arguments | `func f(msgs ...string)` |
| `interface{}` | Any type (like `any` in TypeScript) | `map[string]interface{}` |

### Web & HTTP Terms

| Term | Plain English | Example |
|------|---------------|---------|
| **HTTP** | The protocol browsers use to talk to servers | `http.Get(url)` |
| **GET** | Request to READ data | `GET /artist/3` |
| **Status Code** | Server's response classification | 200=OK, 404=Not Found, 500=Error |
| **Response Writer** | Tool for sending data back to the browser | `w http.ResponseWriter` |
| **Request** | What the browser asked for | `r *http.Request` |
| **Header** | Metadata about the response | `Content-Type: application/json` |
| **Query Parameter** | Data in the URL after `?` | `/api/search?q=queen&sort=name` |
| **AJAX** | Fetch data without reloading the page | `fetch('/api/search?q=queen')` |
| **Endpoint** | A specific URL that accepts requests | `/api/search` is an endpoint |

### Template Terms

| Term | Plain English | Example |
|------|---------------|---------|
| `{{.}}` | The current data context ("dot") | In a loop, `.` is the current item |
| `{{.Name}}` | Access a field of the current data | Outputs the Name field |
| `{{range .}}` | Loop over a list | Repeat for each item |
| `{{if .}}` | Conditional rendering | Show only if value is truthy |
| `{{end}}` | End a range or if block | Required closing tag |
| `{{len .}}` | Get length of a list | `{{len .Members}}` → 4 |

### JavaScript Terms

| Term | Plain English | Example |
|------|---------------|---------|
| `async` | Function that can pause and wait | `async function applyFilters()` |
| `await` | Pause until an async operation finishes | `await fetch(url)` |
| `fetch` | Browser's built-in HTTP client | `fetch('/api/search?q=queen')` |
| `localStorage` | Browser storage that persists | `localStorage.setItem('theme', 'dark')` |
| `innerHTML` | The HTML content inside an element | `div.innerHTML = '<h1>Hi</h1>'` |
| `addEventListener` | React to user actions | `btn.addEventListener('click', fn)` |
| `?.` | Optional chaining (don't crash if null) | `element?.addEventListener(...)` |
| Template literal | String with `${variable}` insertion | `` `Hello ${name}` `` |

---

## Variable Naming (Why Names Look the Way They Do)

| Name | What It Means | Why This Name |
|------|---------------|---------------|
| `w` | Response **w**riter | Go convention for HTTP handlers |
| `r` | HTTP **r**equest | Go convention for HTTP handlers |
| `err` | **Err**or value | Universal Go convention |
| `tmpl` | **T**e**mpl**ate | Common abbreviation |
| `resp` | HTTP **resp**onse | Common abbreviation |
| `idStr` | ID as a **str**ing | Before converting to integer |
| `rel` | **Rel**ation | Short for Relation in loops |
| `i`, `j` | Loop **i**ndex, secondary index | Universal programming convention |
| `a` | An **a**rtist (in a loop) | Short loop variable |
| `q` | Search **q**uery | Common search parameter name |
| `fs` | **F**ile **s**erver | Abbreviation for the static file server |

---

## Abbreviations Decoded

| Abbreviation | Full Form | Context |
|--------------|-----------|---------|
| API | Application Programming Interface | How our app talks to the data source |
| JSON | JavaScript Object Notation | Data format (`{"key":"value"}`) |
| HTTP | HyperText Transfer Protocol | How browsers and servers communicate |
| URL | Uniform Resource Locator | A web address (`http://localhost:8080`) |
| HTML | HyperText Markup Language | The structure of web pages |
| CSS | Cascading Style Sheets | The styling of web pages |
| JS | JavaScript | The programming language of the browser |
| AJAX | Asynchronous JavaScript And XML | Fetching data without page reload |
| MVC | Model-View-Controller | Architecture pattern |
| XSS | Cross-Site Scripting | A security vulnerability |
| DNS | Domain Name System | Translates `google.com` to an IP address |

---

## File Types & Suffixes

| Extension | What It Is | Example |
|-----------|-----------|---------|
| `.go` | Go source code | `main.go` |
| `.html` | HTML template | `templates/home.html` |
| `.css` | Stylesheet | `static/style.css` |
| `.js` | JavaScript | `static/script.js` |
| `.mod` | Go module definition | `go.mod` |
| `_test.go` | Go test file | `search_test.go` |
| `.md` | Markdown documentation | `README.md` |

**Special Go convention:** Files ending in `_test.go` are automatically recognized as test files by `go test`. They're never compiled into the final binary.

---

## Magic Numbers (Why These Specific Values?)

| Value | Where | Why |
|-------|-------|-----|
| `8080` | `main.go` | Common alternative HTTP port (80 is the standard but requires admin privileges) |
| `10 * time.Second` | `services/api.go` | HTTP timeout — 10 seconds balances between "fast enough for users" and "patient enough for slow networks" |
| `200` | `services/api.go` | HTTP status code for "OK, here's your data" (`http.StatusOK`) |
| `400` | `handlers/artist.go` | HTTP status code for "Bad Request" — the user sent invalid data |
| `404` | `handlers/error.go` | HTTP status code for "Not Found" — the requested resource doesn't exist |
| `500` | `handlers/error.go` | HTTP status code for "Internal Server Error" — something broke on the server |

---

## Function Naming Patterns

Go has naming conventions that tell you what a function does:

| Pattern | Meaning | Example |
|---------|---------|---------|
| `Fetch*` | Get data from an external source | `FetchArtists()`, `FetchRelation()` |
| `*Handler` | Responds to an HTTP request | `HomeHandler`, `ArtistHandler` |
| `matches*` | Returns true/false (a predicate) | `matchesSearch()`, `matchesYearFilter()` |
| `calculate*` | Computes a derived value | `calculateTotalConcerts()` |
| `get*` | Retrieves/looks up a value | `getBandType()` |
| `sort*` | Orders a collection | `sortArtists()` |
| `Test*` | A test function (run by `go test`) | `TestMatchesSearch()` |

**Capitalization matters in Go:**
- `FetchArtists` (capital F) = **exported** (public) — can be used from other packages
- `matchesSearch` (lowercase m) = **unexported** (private) — only usable within its own package

---

## Error Messages Decoded

| Error Message | What It Means | What To Do |
|---------------|---------------|------------|
| `"failed to connect to API"` | External API is unreachable | Check internet connection, API might be down |
| `"API returned status code X"` | API responded but with an error code | Check if the API URL is correct |
| `"failed to parse artist data"` | JSON response couldn't be converted to Go structs | API format might have changed |
| `"relation not found for artist X"` | No concert data exists for this artist | The API might not have data for all artists |
| `"Invalid artist ID"` | URL had a non-numeric ID like `/artist/abc` | User typed a wrong URL |
| `"Artist not found"` | Numeric ID doesn't match any artist | ID is out of range |
| `runtime error: nil pointer dereference` | Code tried to access a field on a nil value | A variable is nil that shouldn't be — check error handling |
| `runtime error: index out of range` | Tried to access an item beyond the list's length | List is empty or shorter than expected |

---

## Glossary Complete!

**You now have:**
- Essential vocabulary for every concept in this codebase
- Decoder rings for abbreviations and naming conventions
- Quick reference for error messages
- Understanding of magic numbers and file types

**Congratulations!** You're ready to:
1. Navigate the codebase with confidence
2. Understand WHY code is written the way it is
3. Debug issues by tracing data flow
4. Communicate about this code using shared vocabulary
5. Make changes without fear of breaking things you don't understand

**The journey from "vibecoding" to understanding is complete.**

Go back and revisit any lesson anytime:
- [00 - Overview](00-overview.md) — Project structure and big picture
- [01 - Core Concepts](01-core-concepts.md) — Go syntax and key concepts
- [02 - Data Flow](02-data-flow.md) — How data moves through the app
- [03 - Patterns](03-patterns.md) — Reusable design patterns
- [04 - Line by Line](04-line-by-line.md) — Detailed code walkthrough
- [05 - Exercises](05-exercises.md) — Hands-on practice
- [06 - Gotchas](06-gotchas.md) — Common mistakes to avoid
