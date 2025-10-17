// API Configuration
const API_BASE_URL = window.location.origin + '/api';

// State
let apiKey = localStorage.getItem('mkvmender_api_key') || '';

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    if (apiKey) {
        showSearchSection();
    } else {
        showLoginSection();
    }

    // Event Listeners
    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('search-form').addEventListener('submit', handleSearch);
    document.getElementById('login-btn').addEventListener('click', showLoginSection);
    document.getElementById('logout-btn').addEventListener('click', handleLogout);
});

// Navigation
function showLoginSection() {
    document.getElementById('login-section').style.display = 'block';
    document.getElementById('search-section').style.display = 'none';
    document.getElementById('login-btn').style.display = 'none';
    document.getElementById('logout-btn').style.display = 'none';
}

function showSearchSection() {
    document.getElementById('login-section').style.display = 'none';
    document.getElementById('search-section').style.display = 'block';
    document.getElementById('login-btn').style.display = 'none';
    document.getElementById('logout-btn').style.display = 'inline-block';
}

// Login
async function handleLogin(e) {
    e.preventDefault();

    const key = document.getElementById('api-key').value;
    const errorDiv = document.getElementById('login-error');

    // Test the API key by hitting the health endpoint with auth
    try {
        const response = await fetch(`${API_BASE_URL}/health`, {
            headers: {
                'Authorization': `Bearer ${key}`
            }
        });

        if (response.ok) {
            apiKey = key;
            localStorage.setItem('mkvmender_api_key', key);
            errorDiv.style.display = 'none';
            showSearchSection();
        } else {
            throw new Error('Invalid API key');
        }
    } catch (error) {
        errorDiv.textContent = 'Login failed. Please check your API key.';
        errorDiv.style.display = 'block';
    }
}

// Logout
function handleLogout() {
    apiKey = '';
    localStorage.removeItem('mkvmender_api_key');
    document.getElementById('api-key').value = '';
    document.getElementById('results-container').style.display = 'none';
    showLoginSection();
}

// Search
async function handleSearch(e) {
    e.preventDefault();

    const query = document.getElementById('search-query').value;
    const fuzzy = document.getElementById('fuzzy-search').checked;
    const sortBy = document.getElementById('sort-by').value;
    const errorDiv = document.getElementById('search-error');

    try {
        // Build query params
        const params = new URLSearchParams({
            q: query,
            sort: sortBy,
            fuzzy: fuzzy ? 'true' : 'false'
        });

        const response = await fetch(`${API_BASE_URL}/search?${params}`);

        if (!response.ok) {
            throw new Error('Search failed');
        }

        const data = await response.json();
        displayResults(data);
        errorDiv.style.display = 'none';
    } catch (error) {
        errorDiv.textContent = 'Search failed. Please try again.';
        errorDiv.style.display = 'block';
    }
}

// Display Results
function displayResults(data) {
    const resultsContainer = document.getElementById('results-container');
    const resultsTitle = document.getElementById('results-title');
    const resultsList = document.getElementById('results-list');

    if (!data.results || data.results.length === 0) {
        resultsTitle.textContent = 'No results found';
        resultsList.innerHTML = '<p style="color: var(--text-secondary);">Try a different search term or enable fuzzy matching.</p>';
        resultsContainer.style.display = 'block';
        return;
    }

    resultsTitle.textContent = `Found ${data.results.length} result(s) for "${data.query}"`;

    // Group results by title
    const grouped = groupResults(data.results);

    resultsList.innerHTML = grouped.map(group => {
        const icon = group.media_type === 'tv' ? '📺' : '🎬';
        const year = group.year ? ` (${group.year})` : '';

        return `
            <div class="result-item">
                <div class="result-header">
                    <span class="result-icon">${icon}</span>
                    <span class="result-title">${escapeHtml(group.title)}${year}</span>
                    <span class="result-year">${group.media_type}</span>
                </div>

                <div class="result-meta">
                    <span>📦 ${formatFileSize(group.file_size)}</span>
                    <span>🔢 ${group.submissions.length} submission(s)</span>
                    ${group.season ? `<span>Season ${group.season}</span>` : ''}
                    ${group.episode ? `<span>Episode ${group.episode}</span>` : ''}
                </div>

                <div class="submissions">
                    ${group.submissions.map(sub => `
                        <div class="submission-item">
                            <div class="submission-filename">${escapeHtml(sub.filename)}</div>
                            <div class="submission-meta">
                                <span>👤 ${escapeHtml(sub.username)}</span>
                                <span class="votes">
                                    <span class="vote-up">👍 ${sub.upvotes}</span>
                                    <span class="vote-down">👎 ${sub.downvotes}</span>
                                    <span>Score: ${sub.vote_score}</span>
                                </span>
                            </div>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
    }).join('');

    resultsContainer.style.display = 'block';
}

// Group results by title/year/media_type
function groupResults(results) {
    const groups = {};

    results.forEach(result => {
        const key = `${result.title}-${result.year || 'none'}-${result.media_type}-${result.hash}`;

        if (!groups[key]) {
            groups[key] = {
                title: result.title,
                year: result.year,
                media_type: result.media_type,
                season: result.season,
                episode: result.episode,
                hash: result.hash,
                file_size: result.file_size,
                submissions: [],
                submissionIds: new Set() // Track submission IDs to prevent duplicates
            };
        }

        // Only add submissions we haven't seen before
        result.submissions.forEach(sub => {
            if (!groups[key].submissionIds.has(sub.id)) {
                groups[key].submissionIds.add(sub.id);
                groups[key].submissions.push(sub);
            }
        });
    });

    // Clean up the submissionIds Set before returning
    return Object.values(groups).map(group => {
        delete group.submissionIds;
        return group;
    });
}

// Utilities
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';

    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
