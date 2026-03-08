# Groupie Tracker - Project Overview

## What is this project?

You need to build a **web application** that displays information about music bands and artists by fetching data from an external API and presenting it in a user-friendly way.

## The Big Picture

```
API (groupietrackers.herokuapp.com) 
    ↓ (fetch data)
Your Go Backend Server
    ↓ (process & serve)
Web Frontend (HTML/CSS/JS)
    ↓ (display to)
User's Browser
```

## What You're Building

A website that shows:
- Band/artist information (names, photos, members, etc.)
- Concert locations (where they performed/will perform)
- Concert dates (when they performed/will perform)
- All this data connected together in a meaningful way

## Key Requirements

### 1. **Backend: Go Language**
   - Must use only standard Go packages (no external frameworks)
   - Handle API requests to fetch data
   - Serve HTML pages to users
   - Must be stable (no crashes)

### 2. **Data Source: External API**
   - Base URL: `https://groupietrackers.herokuapp.com/api`
   - Four endpoints:
     - `/artists` - Band info
     - `/locations` - Concert venues
     - `/dates` - Concert dates
     - `/relation` - Links everything together

### 3. **Frontend: HTML/CSS/JavaScript**
   - Display data in creative ways (cards, tables, lists, etc.)
   - Must be user-friendly
   - All pages must work without errors

### 4. **Client-Server Interaction**
   - Implement at least one feature where the client triggers an action
   - The server responds with information
   - Example: Search, filter, click for details, etc.

## Technical Stack

- **Language**: Go (backend only)
- **Packages**: Standard library only
- **Data Format**: JSON
- **Frontend**: HTML, CSS, JavaScript
- **Architecture**: Client-Server model

## What You'll Learn

- Working with REST APIs
- JSON data manipulation
- Go web server development
- HTML templating
- Client-server communication
- Error handling
- Data visualization
