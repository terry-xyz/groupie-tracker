package handlers

import (
	"html/template"
	"log"
	"net/http"
	"sync"
)

var (
	errorTmpl     *template.Template
	errorTmplOnce sync.Once // ensures the template is parsed exactly once across all goroutines
)

// getErrorTmpl parses templates/error.html on the first call and caches it; returns nil if
// parsing failed so ErrorHandler can fall back to plain text instead of panicking.
func getErrorTmpl() *template.Template {
	errorTmplOnce.Do(func() {
		var err error
		errorTmpl, err = template.ParseFiles("templates/error.html")
		if err != nil {
			log.Printf("Error parsing error template: %v", err)
			// errorTmpl stays nil — ErrorHandler checks for nil before calling Execute
		}
	})
	return errorTmpl
}

// ErrorHandler renders the error template with the given HTTP status code and an optional custom message.
func ErrorHandler(w http.ResponseWriter, status int, customMsg ...string) {
	tmpl := getErrorTmpl()
	w.WriteHeader(status) // Must be set before writing the response body
	if tmpl == nil {
		http.Error(w, "Error loading page", status) // Fallback to plain text if the error template itself failed
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
