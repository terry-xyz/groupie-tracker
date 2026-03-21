const searchInput = document.getElementById('searchInput');
const minYearInput = document.getElementById('minYear');
const maxYearInput = document.getElementById('maxYear');
const minAlbumYearInput = document.getElementById('minAlbumYear');
const maxAlbumYearInput = document.getElementById('maxAlbumYear');
const sortBySelect = document.getElementById('sortBy');
const resetBtn = document.getElementById('resetFilters');
const artistGrid = document.getElementById('artistGrid');
const themeToggle = document.getElementById('themeToggle');
const suggestionsDropdown = document.getElementById('suggestionsDropdown');
const locationSearch = document.getElementById('locationSearch');
const locationCheckboxes = document.getElementById('locationCheckboxes');
const resultCount = document.getElementById('resultCount');
const minMembersSlider = document.getElementById('minMembers');
const maxMembersSlider = document.getElementById('maxMembers');
const memberRangeLabel = document.getElementById('memberRangeLabel');

// IIFE runs before DOM paint to prevent theme flash (FOUC) on page load
(function() {
    const savedTheme = localStorage.getItem('theme');
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    document.body.className = savedTheme || (prefersDark ? 'dark-theme' : 'light-theme');
})();

// initTheme syncs the theme toggle icon with the theme already applied by the IIFE above.
function initTheme() {
    const currentTheme = document.body.className;
    updateThemeIcon(currentTheme);
}

// updateThemeIcon sets the toggle button glyph to ☀️ in dark mode and 🌙 in light mode.
function updateThemeIcon(theme) {
    if (themeToggle) {
        themeToggle.textContent = theme === 'dark-theme' ? '☀️' : '🌙';
    }
}

// toggleTheme flips between dark-theme and light-theme, persisting the choice to localStorage.
function toggleTheme() {
    const currentTheme = document.body.className;
    const newTheme = currentTheme === 'dark-theme' ? 'light-theme' : 'dark-theme';
    document.body.className = newTheme;
    localStorage.setItem('theme', newTheme);
    updateThemeIcon(newTheme);
}

if (themeToggle) {
    themeToggle.addEventListener('click', toggleTheme);
}

initTheme();

// --- Debounce utility ---
// debounce wraps fn so it only fires after delay ms of silence — prevents flooding the
// search API on every keystroke by resetting the timer each time the user types.
function debounce(fn, delay) {
    let timer;
    return function() {
        const args = arguments;
        const context = this;
        clearTimeout(timer);
        timer = setTimeout(function() {
            fn.apply(context, args);
        }, delay);
    };
}

// --- Autocomplete suggestions ---
// fetchSuggestions requests categorised suggestions from /api/suggestions and renders the dropdown.
async function fetchSuggestions(query) {
    if (!query || query.length < 1) {
        hideSuggestions();
        return;
    }
    try {
        const response = await fetch('/api/suggestions?q=' + encodeURIComponent(query));
        if (!response.ok) return;
        const suggestions = await response.json();
        displaySuggestions(suggestions);
    } catch (e) {
        // Silently fail — autocomplete is non-critical; the user can still type and search
    }
}

// displaySuggestions renders the dropdown list from the suggestions array returned by the API.
function displaySuggestions(suggestions) {
    if (!suggestionsDropdown) return;
    if (!suggestions || suggestions.length === 0) {
        hideSuggestions();
        return;
    }

    // data-text stores the raw value so clicking a suggestion fills the search box correctly
    suggestionsDropdown.innerHTML = suggestions.map(function(s) {
        return '<div class="suggestion-item" data-text="' + escapeAttr(s.text) + '">' +
            '<span class="suggestion-text">' + escapeHTML(s.text) + '</span>' +
            '<span class="suggestion-category">' + escapeHTML(s.category) + '</span>' +
            '</div>';
    }).join('');
    suggestionsDropdown.style.display = 'block';

    // Attach click handlers after innerHTML is set so the elements exist in the DOM
    suggestionsDropdown.querySelectorAll('.suggestion-item').forEach(function(item) {
        item.addEventListener('click', function() {
            searchInput.value = this.getAttribute('data-text');
            hideSuggestions();
            applyFilters();
        });
    });
}

// hideSuggestions collapses the dropdown without clearing its contents.
function hideSuggestions() {
    if (suggestionsDropdown) {
        suggestionsDropdown.style.display = 'none';
    }
}

// escapeHTML prevents XSS by assigning str as textContent (browser encodes it), then reading
// it back as innerHTML — turns <script> into &lt;script&gt; without a manual replace list.
function escapeHTML(str) {
    var div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}

// escapeAttr escapes quotes so artist names with apostrophes don't break HTML attribute values.
function escapeAttr(str) {
    return str.replace(/"/g, '&quot;').replace(/'/g, '&#39;');
}

if (searchInput) {
    searchInput.addEventListener('input', debounce(function() {
        var val = searchInput.value.trim();
        fetchSuggestions(val);
        applyFilters();
    }, 300));

    searchInput.addEventListener('keyup', function(e) {
        if (e.key === 'Enter') {
            hideSuggestions();
            applyFilters();
        }
        if (e.key === 'Escape') {
            hideSuggestions();
        }
    });
}

// Hide suggestions when clicking outside
document.addEventListener('click', function(e) {
    if (searchInput && !searchInput.contains(e.target) &&
        suggestionsDropdown && !suggestionsDropdown.contains(e.target)) {
        hideSuggestions();
    }
});

// --- Location checkbox loading ---
// loadLocations fetches all unique locations from /api/locations, renders them as grouped
// checkboxes, and restores any previously checked boxes from savedLocations (may be null).
async function loadLocations(savedLocations) {
    if (!locationCheckboxes) return;
    try {
        const response = await fetch('/api/locations');
        if (!response.ok) return;
        const data = await response.json();

        var html = '';
        data.countries.forEach(function(country) {
            html += '<div class="location-country" data-country="' + escapeAttr(country.name) + '">';
            html += '<strong class="country-name">' + escapeHTML(formatLocationLabel(country.name)) + '</strong>';
            country.locations.forEach(function(loc) {
                html += '<label class="checkbox-label location-item" data-location="' + escapeAttr(loc) + '">' +
                    '<input type="checkbox" name="locations" value="' + escapeAttr(loc) + '"> ' +
                    escapeHTML(formatLocationLabel(loc)) + '</label>';
            });
            html += '</div>';
        });
        locationCheckboxes.innerHTML = html;

        // Restore saved location checkboxes, then re-apply so they take effect
        if (savedLocations && savedLocations.length > 0) {
            savedLocations.forEach(function(loc) {
                var cb = locationCheckboxes.querySelector('input[value="' + loc + '"]');
                if (cb) cb.checked = true;
            });
            applyFilters();
        }

        // Auto-trigger filters on checkbox change
        locationCheckboxes.querySelectorAll('input[type="checkbox"]').forEach(function(cb) {
            cb.addEventListener('change', function() {
                applyFilters();
            });
        });
    } catch (e) {
        locationCheckboxes.innerHTML = '<p class="loading-text">Failed to load locations.</p>';
    }
}

// Format "north_carolina-usa" to "North Carolina, Usa"
function formatLocationLabel(raw) {
    return raw.replace(/_/g, ' ').replace(/-/g, ', ').replace(/\b\w/g, function(c) {
        return c.toUpperCase();
    });
}

// Filter visible locations based on search input
if (locationSearch) {
    locationSearch.addEventListener('input', function() {
        var query = this.value.toLowerCase();
        if (!locationCheckboxes) return;
        locationCheckboxes.querySelectorAll('.location-item').forEach(function(item) {
            var loc = item.getAttribute('data-location') || '';
            item.style.display = loc.toLowerCase().indexOf(query) !== -1 ? '' : 'none';
        });
        // Hide country headers with no visible children
        locationCheckboxes.querySelectorAll('.location-country').forEach(function(group) {
            var hasVisible = group.querySelector('.location-item:not([style*="display: none"])');
            group.style.display = hasVisible ? '' : 'none';
        });
    });
}

// --- Filter buttons ---
if (resetBtn) {
    resetBtn.addEventListener('click', resetFilters);
}

// --- Member slider logic ---
// updateMemberLabel refreshes the "Members: X - Y" label above the sliders; shows "8+" when
// max is 8 because the backend treats 8 as "8 or more", not exactly 8.
function updateMemberLabel() {
    if (!memberRangeLabel || !minMembersSlider || !maxMembersSlider) return;
    var min = parseInt(minMembersSlider.value);
    var max = parseInt(maxMembersSlider.value);
    var minLabel = min.toString();
    var maxLabel = max === 8 ? '8+' : max.toString(); // 8 is the slider ceiling and means "8 or more"
    memberRangeLabel.textContent = minLabel + ' - ' + maxLabel;
}

// Initialize member label on page load
if (minMembersSlider && maxMembersSlider) {
    updateMemberLabel();
}

if (minMembersSlider) {
    minMembersSlider.addEventListener('input', function() {
        // If min overtakes max, push max up to match
        if (parseInt(minMembersSlider.value) > parseInt(maxMembersSlider.value)) {
            maxMembersSlider.value = minMembersSlider.value;
        }
        updateMemberLabel();
        debouncedApply();
    });
}

if (maxMembersSlider) {
    maxMembersSlider.addEventListener('input', function() {
        // If max drops below min, push min down to match
        if (parseInt(maxMembersSlider.value) < parseInt(minMembersSlider.value)) {
            minMembersSlider.value = maxMembersSlider.value;
        }
        updateMemberLabel();
        debouncedApply();
    });
}

// --- Auto-trigger filters on change ---
if (sortBySelect) {
    sortBySelect.addEventListener('change', applyFilters);
}

// Debounced input listeners for number fields
var debouncedApply = debounce(applyFilters, 400);
[minYearInput, maxYearInput, minAlbumYearInput, maxAlbumYearInput].forEach(function(input) {
    if (!input) return;
    input.addEventListener('input', debouncedApply);

    // Block non-digit keys (prevents -, +, ., e from being typed into number inputs)
    input.addEventListener('keydown', function(e) {
        var nav = ['Backspace', 'Delete', 'Tab', 'Escape', 'Enter',
                   'ArrowLeft', 'ArrowRight', 'ArrowUp', 'ArrowDown', 'Home', 'End'];
        if (nav.indexOf(e.key) !== -1) return; // always allow navigation/editing keys
        if ((e.ctrlKey || e.metaKey) && 'acvx'.indexOf(e.key.toLowerCase()) !== -1) return; // allow Ctrl+A/C/V/X
        if (!/^\d$/.test(e.key)) e.preventDefault(); // /^\d$/ matches exactly one digit 0-9; anything else is blocked
    });

    // Strip non-digits on paste or autofill — keydown alone can't catch clipboard content
    input.addEventListener('input', function() {
        var cleaned = this.value.replace(/[^\d]/g, ''); // [^\d] matches any non-digit character
        if (this.value !== cleaned) this.value = cleaned;
    });
});

function showLoading() {
    if (artistGrid) {
        artistGrid.innerHTML = '<div class="loading">Loading artists...</div>';
    }
}

function showError(message) {
    if (artistGrid) {
        artistGrid.innerHTML = '<div class="error-message">' + escapeHTML(message) + '</div>';
    }
}

// applyFilters reads all filter inputs, saves them to sessionStorage, then fetches /api/search
// and re-renders the artist grid without a page reload.
async function applyFilters() {
    var query = searchInput ? searchInput.value : '';
    var minYear = minYearInput ? minYearInput.value : '';
    var maxYear = maxYearInput ? maxYearInput.value : '';
    var minAlbumYear = minAlbumYearInput ? minAlbumYearInput.value : '';
    var maxAlbumYear = maxAlbumYearInput ? maxAlbumYearInput.value : '';
    var sort = sortBySelect ? sortBySelect.value : '';

    // Expand slider range into individual member counts (e.g., min=2 max=4 → ["2","3","4"])
    // because the API accepts a comma-separated list of exact counts, not a range
    var members = [];
    if (minMembersSlider && maxMembersSlider) {
        var min = parseInt(minMembersSlider.value);
        var max = parseInt(maxMembersSlider.value);
        for (var i = min; i <= max; i++) {
            members.push(i.toString());
        }
    }

    // Collect checked locations
    var locations = [];
    document.querySelectorAll('input[name="locations"]:checked').forEach(function(cb) {
        locations.push(cb.value);
    });

    // Persist filter state so it survives artist page navigation
    sessionStorage.setItem('gt_filters', JSON.stringify({
        q: query, minYear: minYear, maxYear: maxYear,
        minAlbumYear: minAlbumYear, maxAlbumYear: maxAlbumYear,
        sort: sort, minMembers: minMembersSlider ? minMembersSlider.value : '1',
        maxMembers: maxMembersSlider ? maxMembersSlider.value : '8',
        locations: locations
    }));

    // Only append non-empty params to avoid sending blank query values
    var params = new URLSearchParams();
    if (query) params.append('q', query);
    if (minYear) params.append('minYear', minYear);
    if (maxYear) params.append('maxYear', maxYear);
    if (minAlbumYear) params.append('minAlbumYear', minAlbumYear);
    if (maxAlbumYear) params.append('maxAlbumYear', maxAlbumYear);
    if (sort) params.append('sort', sort);
    if (members.length > 0) params.append('members', members.join(','));
    if (locations.length > 0) params.append('locations', locations.join(','));

    showLoading();

    try {
        var response = await fetch('/api/search?' + params);

        if (!response.ok) {
            throw new Error('Failed to fetch artists');
        }

        var artists = await response.json();

        if (artists.error) {
            showError(artists.error);
            return;
        }

        displayArtists(artists);

        // Update result count; ternary handles "1 artist" vs "N artists" pluralisation
        if (resultCount) {
            resultCount.textContent = 'Showing ' + artists.length + ' artist' + (artists.length !== 1 ? 's' : '');
        }
    } catch (error) {
        console.error('Error fetching artists:', error);
        showError('Unable to load artists. Please check your connection and try again.');
    }
}

// displayArtists replaces the artist grid contents with cards built from the API response array.
function displayArtists(artists) {
    if (!artistGrid) return;

    if (artists.length === 0) {
        artistGrid.innerHTML = '<p class="no-results">No artists found matching your criteria.</p>';
        return;
    }

    // Build card HTML from API response; uses JSON field names (camelCase) from Go's json tags
    artistGrid.innerHTML = artists.map(function(artist) {
        return '<div class="artist-card">' +
            '<img src="' + escapeAttr(artist.image) + '" alt="' + escapeAttr(artist.name) + '">' +
            '<h2>' + escapeHTML(artist.name) + '</h2>' +
            '<p>Created: ' + artist.creationDate + '</p>' +
            '<p>First Album: ' + escapeHTML(artist.firstAlbum) + '</p>' +
            '<a href="/artist/' + artist.id + '" class="btn">View Details</a>' +
            '</div>';
    }).join('');
}

// resetFilters clears every filter input, removes the saved sessionStorage state so it isn't
// restored on the next page visit, then re-fetches all artists without a page reload.
function resetFilters() {
    if (searchInput) searchInput.value = '';
    if (minYearInput) minYearInput.value = '';
    if (maxYearInput) maxYearInput.value = '';
    if (minAlbumYearInput) minAlbumYearInput.value = '';
    if (maxAlbumYearInput) maxAlbumYearInput.value = '';
    if (sortBySelect) sortBySelect.value = '';
    if (locationSearch) locationSearch.value = '';

    // Reset member sliders
    if (minMembersSlider) minMembersSlider.value = '1';
    if (maxMembersSlider) maxMembersSlider.value = '8';
    updateMemberLabel();

    // Uncheck all location checkboxes
    document.querySelectorAll('input[name="locations"]:checked').forEach(function(cb) {
        cb.checked = false;
    });

    // Show all location items again
    if (locationCheckboxes) {
        locationCheckboxes.querySelectorAll('.location-item').forEach(function(item) {
            item.style.display = '';
        });
        locationCheckboxes.querySelectorAll('.location-country').forEach(function(group) {
            group.style.display = '';
        });
    }

    sessionStorage.removeItem('gt_filters');
    applyFilters();
}

// Restore saved filter state when returning from an artist page.
// Only runs on the home page where locationCheckboxes exists.
if (locationCheckboxes) {
    var _saved = null;
    try { _saved = JSON.parse(sessionStorage.getItem('gt_filters')); } catch(e) {} // corrupt storage → treat as no saved state
    if (_saved) {
        if (searchInput) searchInput.value = _saved.q || '';
        if (minYearInput) minYearInput.value = _saved.minYear || '';
        if (maxYearInput) maxYearInput.value = _saved.maxYear || '';
        if (minAlbumYearInput) minAlbumYearInput.value = _saved.minAlbumYear || '';
        if (maxAlbumYearInput) maxAlbumYearInput.value = _saved.maxAlbumYear || '';
        if (sortBySelect) sortBySelect.value = _saved.sort || '';
        if (minMembersSlider) minMembersSlider.value = _saved.minMembers || '1'; // fall back to slider min if missing
        if (maxMembersSlider) maxMembersSlider.value = _saved.maxMembers || '8'; // fall back to slider max if missing
        updateMemberLabel();
        applyFilters(); // re-render grid with restored non-location filters immediately
    }
    // loadLocations is async; it restores saved location checkboxes after the DOM is built
    // and calls applyFilters() again if any locations were checked, so the grid updates once more
    loadLocations(_saved ? _saved.locations : null);
}
