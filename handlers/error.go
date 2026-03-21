package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// ErrorHandler renders the error template with the given HTTP status code and an optional custom message.
func ErrorHandler(w http.ResponseWriter, status int, customMsg ...string) {
	w.WriteHeader(status) // Must be set before writing the response body
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, "Error loading page", status) // Fallback to plain text if the error template itself fails
		return
	}

	message := http.StatusText(status) // Default to standard HTTP status text (e.g., "Not Found")
	if len(customMsg) > 0 && customMsg[0] != "" {
		message = customMsg[0] // Variadic param lets callers optionally override the default message
	}

	data := map[string]interface{}{
		"Status":  status,
		"Message": message,
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error rendering error template (status %d): %v", status, err)
	}
}
