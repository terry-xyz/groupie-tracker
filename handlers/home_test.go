package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Verifies that any path other than "/" returns 404 (HomeHandler is not a catch-all)
func TestHomeHandlerRejectsNonRootPaths(t *testing.T) {
	// Mix of unregistered, partial, and nested paths — all must be rejected without hitting the API
	paths := []string{"/unknown", "/artist", "/foo/bar", "/api"}

	for _, path := range paths {
		req := httptest.NewRequest("GET", path, nil)
		w := httptest.NewRecorder()

		HomeHandler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("HomeHandler(%q) = %d, expected 404", path, w.Code)
		}
	}
}
