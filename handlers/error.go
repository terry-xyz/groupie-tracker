package handlers

import (
	"html/template"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, status int, customMsg ...string) {
	w.WriteHeader(status)
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, "Error loading page", status)
		return
	}

	message := http.StatusText(status)
	if len(customMsg) > 0 && customMsg[0] != "" {
		message = customMsg[0]
	}

	data := map[string]interface{}{
		"Status":  status,
		"Message": message,
	}
	tmpl.Execute(w, data)
}
