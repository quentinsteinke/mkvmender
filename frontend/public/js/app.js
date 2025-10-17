// API Configuration
const API_BASE_URL = window.location.origin + '/api';

// State
let apiKey = localStorage.getItem('mkvmender_api_key') || '';

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    // Set up OS-specific download button
    setupDownloadButton();

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

// Detect OS and set download button
function setupDownloadButton() {
    const downloadBtn = document.getElementById('download-btn');
    const downloadText = document.getElementById('download-text');
    const userAgent = navigator.userAgent.toLowerCase();
    const platform = navigator.platform.toLowerCase();

    let os = 'unknown';
    let downloadUrl = '';
    let label = 'Download';

    // Detect OS
    if (platform.includes('mac') || userAgent.includes('mac')) {
        // Detect Apple Silicon vs Intel
        const isAppleSilicon = userAgent.includes('arm') || platform.includes('arm');
        if (isAppleSilicon) {
            os = 'mac-arm';
            downloadUrl = 'https://github.com/quentinsteinke/mkvmender/releases/download/v1.0.0/mkvmender-darwin-arm64';
            label = 'Download for macOS (Apple Silicon)';
        } else {
            os = 'mac-intel';
            downloadUrl = 'https://github.com/quentinsteinke/mkvmender/releases/download/v1.0.0/mkvmender-darwin-amd64';
            label = 'Download for macOS (Intel)';
        }
    } else if (platform.includes('win') || userAgent.includes('win')) {
        os = 'windows';
        downloadUrl = 'https://github.com/quentinsteinke/mkvmender/releases/download/v1.0.0/mkvmender-windows-amd64.exe';
        label = 'Download for Windows';
    } else if (platform.includes('linux') || userAgent.includes('linux')) {
        // Detect ARM vs x86_64
        const isArm = userAgent.includes('arm') || platform.includes('arm') || userAgent.includes('aarch64');
        if (isArm) {
            os = 'linux-arm';
            downloadUrl = 'https://github.com/quentinsteinke/mkvmender/releases/download/v1.0.0/mkvmender-linux-arm64';
            label = 'Download for Linux (ARM64)';
        } else {
            os = 'linux-x86';
            downloadUrl = 'https://github.com/quentinsteinke/mkvmender/releases/download/v1.0.0/mkvmender-linux-amd64';
            label = 'Download for Linux (x86_64)';
        }
    }

    // Set button properties
    if (downloadUrl) {
        downloadBtn.href = downloadUrl;
        downloadText.textContent = label;
    } else {
        // Fallback to releases page if OS not detected
        downloadBtn.href = 'https://github.com/quentinsteinke/mkvmender/releases';
        downloadText.textContent = 'Download';
        downloadBtn.removeAttribute('download');
        downloadBtn.target = '_blank';
    }
}

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
                            <div class="submission-header">
                                <div class="submission-filename">${escapeHtml(sub.filename)}</div>
                                <div class="vote-buttons">
                                    <button class="vote-btn upvote" data-submission-id="${sub.id}" data-vote-type="1">
                                        👍 <span class="vote-count">${sub.upvotes}</span>
                                    </button>
                                    <button class="vote-btn downvote" data-submission-id="${sub.id}" data-vote-type="-1">
                                        👎 <span class="vote-count">${sub.downvotes}</span>
                                    </button>
                                </div>
                            </div>
                            <div class="submission-meta">
                                <span>👤 ${escapeHtml(sub.username)}</span>
                                <span class="vote-score">Score: ${sub.vote_score}</span>
                            </div>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
    }).join('');

    resultsContainer.style.display = 'block';

    // Attach vote button event listeners
    attachVoteListeners();
}

// Attach event listeners to vote buttons
function attachVoteListeners() {
    const voteButtons = document.querySelectorAll('.vote-btn');
    voteButtons.forEach(button => {
        button.addEventListener('click', async (e) => {
            const submissionId = parseInt(button.dataset.submissionId);
            const voteType = parseInt(button.dataset.voteType);
            await handleVote(submissionId, voteType, button);
        });
    });
}

// Handle voting
async function handleVote(submissionId, voteType, button) {
    // Disable button during request
    button.disabled = true;

    try {
        const response = await fetch(`${API_BASE_URL}/vote`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${apiKey}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                submission_id: submissionId,
                vote_type: voteType
            })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Vote failed');
        }

        const result = await response.json();

        // Update both vote count buttons
        const submissionItem = button.closest('.submission-item');
        const upvoteBtn = submissionItem.querySelector('.vote-btn.upvote .vote-count');
        const downvoteBtn = submissionItem.querySelector('.vote-btn.downvote .vote-count');

        if (upvoteBtn) {
            upvoteBtn.textContent = result.upvotes;
        }
        if (downvoteBtn) {
            downvoteBtn.textContent = result.downvotes;
        }

        // Update the score display
        const scoreSpan = submissionItem.querySelector('.vote-score');
        if (scoreSpan) {
            scoreSpan.textContent = `Score: ${result.vote_score}`;
        }

        // Add voted state to the clicked button
        button.classList.add('voted');

        // Show success feedback
        showVoteFeedback(submissionItem, 'Vote recorded!', 'success');

    } catch (error) {
        console.error('Vote failed:', error);
        showVoteFeedback(button.closest('.submission-item'), error.message, 'error');
    } finally {
        // Re-enable button
        button.disabled = false;
    }
}

// Show vote feedback message
function showVoteFeedback(element, message, type) {
    const feedback = document.createElement('div');
    feedback.className = type === 'success' ? 'success' : 'error';
    feedback.textContent = message;
    feedback.style.marginTop = '0.5rem';
    feedback.style.fontSize = '0.85rem';

    element.appendChild(feedback);

    // Remove feedback after 3 seconds
    setTimeout(() => {
        feedback.remove();
    }, 3000);
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
