const searchInput = document.getElementById('searchInput');
const minYearInput = document.getElementById('minYear');
const maxYearInput = document.getElementById('maxYear');
const sortBySelect = document.getElementById('sortBy');
const applyBtn = document.getElementById('applyFilters');
const resetBtn = document.getElementById('resetFilters');
const artistGrid = document.getElementById('artistGrid');
const themeToggle = document.getElementById('themeToggle');

// Apply theme immediately (before page renders)
(function() {
    const savedTheme = localStorage.getItem('theme');
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    document.body.className = savedTheme || (prefersDark ? 'dark-theme' : 'light-theme');
})();

// Theme Management
function initTheme() {
    const currentTheme = document.body.className;
    updateThemeIcon(currentTheme);
}

function updateThemeIcon(theme) {
    themeToggle.textContent = theme === 'dark-theme' ? '☀️' : '🌙';
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

if (applyBtn) {
    applyBtn.addEventListener('click', applyFilters);
}

if (resetBtn) {
    resetBtn.addEventListener('click', resetFilters);
}

if (searchInput) {
    searchInput.addEventListener('keyup', (e) => {
        if (e.key === 'Enter') {
            applyFilters();
        }
    });
}

function showLoading() {
    artistGrid.innerHTML = '<div class="loading">Loading artists...</div>';
}

function showError(message) {
    artistGrid.innerHTML = `<div class="error-message">${message}</div>`;
}

async function applyFilters() {
    const query = searchInput.value;
    const minYear = minYearInput.value;
    const maxYear = maxYearInput.value;
    const sort = sortBySelect.value;

    const params = new URLSearchParams();
    if (query) params.append('q', query);
    if (minYear) params.append('minYear', minYear);
    if (maxYear) params.append('maxYear', maxYear);
    if (sort) params.append('sort', sort);

    showLoading();

    try {
        const response = await fetch(`/api/search?${params}`);
        
        if (!response.ok) {
            throw new Error('Failed to fetch artists');
        }
        
        const artists = await response.json();
        
        if (artists.error) {
            showError(artists.error);
            return;
        }
        
        displayArtists(artists);
    } catch (error) {
        console.error('Error fetching artists:', error);
        showError('Unable to load artists. Please check your connection and try again.');
    }
}

function displayArtists(artists) {
    if (artists.length === 0) {
        artistGrid.innerHTML = '<p class="no-results">No artists found matching your criteria.</p>';
        return;
    }

    artistGrid.innerHTML = artists.map(artist => `
        <div class="artist-card">
            <img src="${artist.image}" alt="${artist.name}">
            <h2>${artist.name}</h2>
            <p>Created: ${artist.creationDate}</p>
            <p>First Album: ${artist.firstAlbum}</p>
            <a href="/artist/${artist.id}" class="btn">View Details</a>
        </div>
    `).join('');
}

function resetFilters() {
    searchInput.value = '';
    minYearInput.value = '';
    maxYearInput.value = '';
    sortBySelect.value = '';
    location.reload();
}
