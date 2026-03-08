# Implementation Tasks

## Phase 1: Setup & API Understanding ✓

### Task 1.1: Explore the API
- [ ] Visit `https://groupietrackers.herokuapp.com/api`
- [ ] Check each endpoint:
  - [ ] `/api/artists`
  - [ ] `/api/locations`
  - [ ] `/api/dates`
  - [ ] `/api/relation`
- [ ] Understand the JSON structure of each endpoint
- [ ] Note the relationships between data

### Task 1.2: Project Structure
- [ ] Create project folders:
  ```
  groupie-tracker/
  ├── main.go              (entry point)
  ├── handlers/            (HTTP handlers)
  ├── models/              (data structures)
  ├── services/            (API fetching logic)
  ├── templates/           (HTML files)
  ├── static/              (CSS, JS, images)
  └── tests/               (unit tests)
  ```

---

## Phase 2: Backend Development

### Task 2.1: Define Data Models
- [ ] Create Go structs for:
  - [ ] Artist (name, image, year, first album, members)
  - [ ] Location (concert venues)
  - [ ] Date (concert dates)
  - [ ] Relation (linking all data)

### Task 2.2: API Client Service
- [ ] Create function to fetch data from API
- [ ] Parse JSON responses into Go structs
- [ ] Handle errors (network issues, invalid JSON, etc.)
- [ ] Consider caching data to reduce API calls

### Task 2.3: HTTP Server Setup
- [ ] Initialize HTTP server using `net/http`
- [ ] Define routes:
  - [ ] `/` - Home page (list all artists)
  - [ ] `/artist/{id}` - Individual artist details
  - [ ] `/api/search` - Search functionality (client-server feature)
  - [ ] `/static/` - Serve CSS/JS files
- [ ] Implement handlers for each route

### Task 2.4: Template Rendering
- [ ] Use `html/template` package
- [ ] Create templates for:
  - [ ] Home page
  - [ ] Artist detail page
  - [ ] Error pages (404, 500)
- [ ] Pass data from backend to templates

---

## Phase 3: Frontend Development

### Task 3.1: Home Page
- [ ] Display all artists in a grid/list
- [ ] Show: name, image, formation year
- [ ] Make each artist clickable (links to detail page)
- [ ] Add basic styling (CSS)

### Task 3.2: Artist Detail Page
- [ ] Show full artist information
- [ ] Display members list
- [ ] Show concert locations (map or list)
- [ ] Show concert dates
- [ ] Link locations with dates using relation data

### Task 3.3: Client-Server Feature
Choose ONE interactive feature:
- [ ] **Option A**: Search bar (search by artist name, member, location)
- [ ] **Option B**: Filter (by year, location, date range)
- [ ] **Option C**: Sort (by name, year, etc.)
- [ ] Implement using JavaScript fetch/AJAX
- [ ] Server responds with filtered/searched data

---

## Phase 4: Error Handling & Stability

### Task 4.1: Error Pages
- [ ] Create 404 page (not found)
- [ ] Create 500 page (server error)
- [ ] Handle invalid routes gracefully

### Task 4.2: Input Validation
- [ ] Validate user inputs (search queries, IDs)
- [ ] Prevent crashes from bad data
- [ ] Return meaningful error messages

### Task 4.3: API Error Handling
- [ ] Handle API unavailability
- [ ] Handle timeout errors
- [ ] Display user-friendly messages when data can't be fetched

---

## Phase 5: Testing & Polish

### Task 5.1: Unit Tests
- [ ] Test API fetching functions
- [ ] Test data parsing
- [ ] Test handlers
- [ ] Run: `go test ./...`

### Task 5.2: Code Quality
- [ ] Follow Go best practices
- [ ] Add comments to complex functions
- [ ] Format code: `go fmt ./...`
- [ ] Check for errors: `go vet ./...`

### Task 5.3: Final Testing
- [ ] Test all pages manually
- [ ] Test on different browsers
- [ ] Verify no crashes occur
- [ ] Test error scenarios

---

## Phase 6: Documentation

### Task 6.1: README
- [ ] How to run the project
- [ ] Dependencies (should be none except Go stdlib)
- [ ] Features implemented
- [ ] Screenshots (optional)

### Task 6.2: Code Documentation
- [ ] Add package comments
- [ ] Document exported functions
- [ ] Explain complex logic

---

## Estimated Timeline

- **Phase 1**: 1-2 hours
- **Phase 2**: 4-6 hours
- **Phase 3**: 4-6 hours
- **Phase 4**: 2-3 hours
- **Phase 5**: 2-3 hours
- **Phase 6**: 1 hour

**Total**: ~15-20 hours
