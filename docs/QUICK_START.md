# Quick Start Guide

## TL;DR - What You Need to Do

Build a website in Go that:
1. Fetches band data from an API
2. Displays it nicely in HTML
3. Has at least one interactive feature (search, filter, etc.)
4. Never crashes

---

## Step-by-Step Getting Started

### Step 1: Explore the API (5 minutes)
Open these URLs in your browser:
- https://groupietrackers.herokuapp.com/api/artists
- https://groupietrackers.herokuapp.com/api/locations
- https://groupietrackers.herokuapp.com/api/dates
- https://groupietrackers.herokuapp.com/api/relation

Look at the JSON structure. This is the data you'll work with.

### Step 2: Create Project Structure (5 minutes)
```bash
cd groupie-tracker
mkdir -p handlers models services templates static tests
touch main.go
```

### Step 3: Start with a Simple Server (15 minutes)
Create a basic Go server that just says "Hello World"

### Step 4: Fetch API Data (30 minutes)
Write code to fetch and parse the artists JSON

### Step 5: Display Data (1 hour)
Create HTML templates and show the artists on a webpage

### Step 6: Add Details Page (1 hour)
Make each artist clickable to see full details

### Step 7: Add Interactive Feature (1-2 hours)
Implement search or filter functionality

### Step 8: Error Handling (1 hour)
Make sure nothing crashes, add error pages

### Step 9: Testing (1 hour)
Write basic tests and manually test everything

---

## Minimum Viable Product (MVP)

To pass this project, you MUST have:

✅ **Backend in Go** (using only standard library)
✅ **Fetch data from the API** (all 4 endpoints)
✅ **Home page** showing all artists
✅ **Detail page** for each artist with locations and dates
✅ **One client-server feature** (e.g., search)
✅ **Error handling** (no crashes, 404/500 pages)
✅ **Clean code** following Go best practices

---

## Common Pitfalls to Avoid

❌ Using external frameworks (only stdlib allowed!)
❌ Not handling errors (will cause crashes)
❌ Forgetting the client-server interaction feature
❌ Not testing edge cases
❌ Hardcoding data instead of fetching from API
❌ Not validating user input

---

## Recommended Order of Implementation

1. **Basic server** → Get something running
2. **API fetching** → Get the data
3. **Home page** → Display artists list
4. **Detail page** → Show full info
5. **Styling** → Make it look decent
6. **Interactive feature** → Add search/filter
7. **Error handling** → Make it robust
8. **Testing** → Verify everything works

---

## Key Go Packages You'll Use

- `net/http` - HTTP server and client
- `html/template` - Render HTML
- `encoding/json` - Parse JSON
- `fmt` - Formatting
- `log` - Logging
- `strings` - String manipulation
- `strconv` - String conversions

---

## Success Criteria

Your project is complete when:
- ✅ Server runs without crashing
- ✅ All artists are displayed on home page
- ✅ Clicking an artist shows their details
- ✅ Concert locations and dates are visible
- ✅ Search/filter feature works
- ✅ Error pages appear for invalid URLs
- ✅ Code is clean and well-organized
- ✅ Basic tests are written

---

## Need Help?

1. Check `TECHNICAL_GUIDE.md` for code examples
2. Check `TASKS.md` for detailed breakdown
3. Check `PROJECT_OVERVIEW.md` for concept explanation
4. Test the API endpoints in your browser first
5. Start simple, add features incrementally

---

## Time Estimate

- **Minimum**: 10-12 hours (basic functionality)
- **Comfortable**: 15-20 hours (with good styling and features)
- **Polished**: 25+ hours (with advanced features and tests)

Good luck! 🚀
