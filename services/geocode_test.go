package services

import (
	"testing"
)

// TestFormatLocationName verifies location parsing from raw API keys
func TestFormatLocationName(t *testing.T) {
	tests := []struct {
		raw             string
		expectedCity    string
		expectedCountry string
	}{
		{"north_carolina-usa", "North Carolina", "Usa"},
		{"london-uk", "London", "Uk"},
		{"new_york-usa", "New York", "Usa"},
		{"paris-france", "Paris", "France"},
		{"singleword", "Singleword", ""},
		{"los_angeles-california-usa", "Los Angeles-california", "Usa"}, // last dash splits country; inner dashes remain
	}

	for _, test := range tests {
		city, country := FormatLocationName(test.raw)
		if city != test.expectedCity {
			t.Errorf("FormatLocationName(%q) city = %q, expected %q", test.raw, city, test.expectedCity)
		}
		if country != test.expectedCountry {
			t.Errorf("FormatLocationName(%q) country = %q, expected %q", test.raw, country, test.expectedCountry)
		}
	}
}

// TestTitleCase verifies custom title-casing
func TestTitleCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello world", "Hello World"},
		{"HELLO", "HELLO"},
		{"a", "A"},
		{"", ""},
		{"one two three", "One Two Three"},
	}

	for _, test := range tests {
		result := TitleCase(test.input)
		if result != test.expected {
			t.Errorf("TitleCase(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}
