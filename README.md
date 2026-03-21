# Groupie Tracker

A Go web app that displays music artist profiles, concert history, and tour maps by consuming an external REST API.

![Go Version](https://img.shields.io/badge/Go-1.24-blue)
![License](https://img.shields.io/badge/license-MIT-green)

## Features

- **Artist profiles** — members, formation year, first album, band type, concert stats
- **Search & autocomplete** — by name, member, location, or date with categorized suggestions
- **Filters** — creation year range, first album year range, member count, location
- **Sorting** — by name, newest, oldest, or default
- **Tour map** — interactive world map with geocoded concert markers and route polylines (Leaflet + OpenStreetMap)
- **Responsive UI** — glassmorphism design with dark/light theme toggle

## Installation

```bash
git clone https://platform.zone01.gr/git/ckotsalas/groupie-tracker.git
cd groupie-tracker
go run main.go
```

Open `http://localhost:8080`.

## Testing

```bash
go test ./...
go test ./... -cover
```

## Project Structure

```
groupie-tracker/
├── main.go
├── handlers/       # HTTP handlers + tests
├── models/         # Data structures
├── services/       # API client, geocoding + tests
├── templates/      # HTML templates
└── static/         # CSS, JS
```
