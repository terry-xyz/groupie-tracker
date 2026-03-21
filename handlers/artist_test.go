package handlers

import (
	"groupie-tracker/models"
	"testing"
)

// Verifies band type classification by member count
func TestGetBandType(t *testing.T) {
	tests := []struct {
		count    int
		expected string
	}{
		{1, "Solo Artist"},
		{2, "Duo"},
		{3, "Trio"},
		{4, "Quartet"},
		{5, "Quintet"},
		{6, "Band"},  // any count above 5 falls through to the default "Band" label
		{7, "Band"},
		{10, "Band"},
	}

	for _, test := range tests {
		result := getBandType(test.count)
		if result != test.expected {
			t.Errorf("getBandType(%d) = %q, expected %q", test.count, result, test.expected)
		}
	}
}

// Verifies total concert count summed across all locations
func TestCalculateTotalConcerts(t *testing.T) {
	tests := []struct {
		relation *models.Relation
		expected int
	}{
		{nil, 0}, // nil relation means no concert data available, so count is 0
		{
			&models.Relation{DatesLocations: map[string][]string{}},
			0, // empty map: locations exist but none have dates
		},
		{
			&models.Relation{DatesLocations: map[string][]string{
				"london-uk": {"01-01-2020", "02-01-2020"},
			}},
			2, // single location with 2 dates
		},
		{
			&models.Relation{DatesLocations: map[string][]string{
				"london-uk":          {"01-01-2020"},
				"north_carolina-usa": {"05-05-2019", "06-05-2019", "07-05-2019"},
				"paris-france":       {"10-10-2018"},
			}},
			5, // 1 + 3 + 1 dates across three locations
		},
	}

	for _, test := range tests {
		result := calculateTotalConcerts(test.relation)
		if result != test.expected {
			t.Errorf("calculateTotalConcerts() = %d, expected %d", result, test.expected)
		}
	}
}

// Verifies unique country count extracted from location keys
func TestCalculateTotalCountries(t *testing.T) {
	tests := []struct {
		relation *models.Relation
		expected int
	}{
		{nil, 0}, // nil relation means no location data, so country count is 0
		{
			&models.Relation{DatesLocations: map[string][]string{}},
			0, // empty map: no locations to extract countries from
		},
		{
			&models.Relation{DatesLocations: map[string][]string{
				"london-uk": {"01-01-2020"},
			}},
			1, // single location means single country
		},
		{
			// 5 locations across only 3 unique countries — verifies deduplication
			&models.Relation{DatesLocations: map[string][]string{
				"london-uk":          {"01-01-2020"},
				"manchester-uk":      {"02-02-2020"}, // same country as london-uk, must not be double-counted
				"north_carolina-usa": {"03-03-2020"},
				"new_york-usa":       {"04-04-2020"}, // same country as north_carolina-usa, must not be double-counted
				"paris-france":       {"05-05-2020"},
			}},
			3, // uk, usa, france
		},
	}

	for _, test := range tests {
		result := calculateTotalCountries(test.relation)
		if result != test.expected {
			t.Errorf("calculateTotalCountries() = %d, expected %d", result, test.expected)
		}
	}
}
