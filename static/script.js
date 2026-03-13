const searchInput = document.getElementById('searchInput');
const minYearInput = document.getElementById('minYear');
const maxYearInput = document.getElementById('maxYear');
const minAlbumYearInput = document.getElementById('minAlbumYear');
const maxAlbumYearInput = document.getElementById('maxAlbumYear');
const sortBySelect = document.getElementById('sortBy');
const applyBtn = document.getElementById('applyFilters');
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

function initTheme() {
    const currentTheme = document.body.className;
    updateThemeIcon(currentTheme);
}

function updateThemeIcon(theme) {
    if (themeToggle) {
        themeToggle.textContent = theme === 'dark-theme' ? '☀️' : '🌙';
    }
}

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
        // Silently fail on autocomplete errors
    }
}

function displaySuggestions(suggestions) {
    if (!suggestionsDropdown) return;
    if (!suggestions || suggestions.length === 0) {
        hideSuggestions();
        return;
    }

    suggestionsDropdown.innerHTML = suggestions.map(function(s) {
        return '<div class="suggestion-item" data-text="' + escapeAttr(s.text) + '">' +
            '<span class="suggestion-text">' + escapeHTML(s.text) + '</span>' +
            '<span class="suggestion-category">' + escapeHTML(s.category) + '</span>' +
            '</div>';
    }).join('');
    suggestionsDropdown.style.display = 'block';

    // Attach click handlers
    suggestionsDropdown.querySelectorAll('.suggestion-item').forEach(function(item) {
        item.addEventListener('click', function() {
            searchInput.value = this.getAttribute('data-text');
            hideSuggestions();
            applyFilters();
        });
    });
}

function hideSuggestions() {
    if (suggestionsDropdown) {
        suggestionsDropdown.style.display = 'none';
    }
}

function escapeHTML(str) {
    var div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}

function escapeAttr(str) {
    return str.replace(/"/g, '&quot;').replace(/'/g, '&#39;');
}

if (searchInput) {
    searchInput.addEventListener('input', debounce(function() {
        fetchSuggestions(searchInput.value.trim());
    }, 250));

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
async function loadLocations() {
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
if (applyBtn) {
    applyBtn.addEventListener('click', applyFilters);
}

if (resetBtn) {
    resetBtn.addEventListener('click', resetFilters);
}

// --- Member slider logic ---
function updateMemberLabel() {
    if (!memberRangeLabel || !minMembersSlider || !maxMembersSlider) return;
    var min = parseInt(minMembersSlider.value);
    var max = parseInt(maxMembersSlider.value);
    
    // Ensure min doesn't exceed max
    if (min > max) {
        minMembersSlider.value = max;
        min = max;
    }
    
    var minLabel = min.toString();
    var maxLabel = max === 8 ? '8+' : max.toString();
    memberRangeLabel.textContent = minLabel + ' - ' + maxLabel;
}

// Initialize member label on page load
if (minMembersSlider && maxMembersSlider) {
    updateMemberLabel();
}

if (minMembersSlider) {
    minMembersSlider.addEventListener('input', function() {
        updateMemberLabel();
        debouncedApply();
    });
}

if (maxMembersSlider) {
    maxMembersSlider.addEventListener('input', function() {
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
    if (input) {
        input.addEventListener('input', debouncedApply);
    }
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

async function applyFilters() {
    var query = searchInput ? searchInput.value : '';
    var minYear = minYearInput ? minYearInput.value : '';
    var maxYear = maxYearInput ? maxYearInput.value : '';
    var minAlbumYear = minAlbumYearInput ? minAlbumYearInput.value : '';
    var maxAlbumYear = maxAlbumYearInput ? maxAlbumYearInput.value : '';
    var sort = sortBySelect ? sortBySelect.value : '';

    // Collect member range from sliders
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

        // Update result count
        if (resultCount) {
            resultCount.textContent = 'Showing ' + artists.length + ' artist' + (artists.length !== 1 ? 's' : '');
        }
    } catch (error) {
        console.error('Error fetching artists:', error);
        showError('Unable to load artists. Please check your connection and try again.');
    }
}

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

    // Clear result count
    if (resultCount) {
        resultCount.textContent = '';
    }

    location.reload(); // Reload restores the server-rendered artist grid
}

// Load locations on page load (only on home page)
if (locationCheckboxes) {
    loadLocations();
}
