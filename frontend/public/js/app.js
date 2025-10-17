// API Configuration
const API_BASE_URL = window.location.origin + '/api';

// State
let apiKey = localStorage.getItem('mkvmender_api_key') || '';
let userRole = localStorage.getItem('mkvmender_user_role') || 'user';
let isAdmin = userRole === 'admin' || userRole === 'moderator';

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    // Set up OS-specific download button
    setupDownloadButton();

    if (apiKey) {
        // Verify API key and fetch user data
        verifyAndShowApp();
    } else {
        showLoginSection();
    }

    // Event Listeners
    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('search-form').addEventListener('submit', handleSearch);
    document.getElementById('login-btn').addEventListener('click', showLoginSection);
    document.getElementById('logout-btn').addEventListener('click', handleLogout);
    document.getElementById('admin-panel-btn').addEventListener('click', showAdminSection);
    document.getElementById('search-panel-btn').addEventListener('click', showSearchSection);
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
        // macOS detection is tricky because browsers report "MacIntel" even on Apple Silicon
        // We'll default to Apple Silicon (ARM64) since:
        // 1. Most Macs sold since 2020 are Apple Silicon
        // 2. Users can click "All platforms" if they need Intel version

        os = 'mac-arm';
        downloadUrl = 'https://github.com/quentinsteinke/mkvmender/releases/download/v1.0.0/mkvmender-darwin-arm64';
        label = 'Download for macOS';
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

// Verify API key and show appropriate section
async function verifyAndShowApp() {
    try {
        // Fetch user data to verify API key and get role
        const response = await fetch(`${API_BASE_URL}/verify`, {
            headers: {
                'Authorization': `Bearer ${apiKey}`
            }
        });

        if (response.ok) {
            const userData = await response.json();

            // Check if user is active
            if (!userData.is_active) {
                alert('Your account has been suspended. Please contact support.');
                handleLogout();
                return;
            }

            // Set user role
            userRole = userData.role;
            localStorage.setItem('mkvmender_user_role', userRole);

            // Check if user is admin or moderator
            isAdmin = userRole === 'admin' || userRole === 'moderator';

            // Initialize admin panel if user is admin
            if (isAdmin && typeof initAdminPanel === 'function') {
                initAdminPanel(apiKey);
            }

            showSearchSection();
        } else {
            // Invalid API key, show login
            handleLogout();
        }
    } catch (error) {
        console.error('Error verifying API key:', error);
        handleLogout();
    }
}

// Navigation
function showLoginSection() {
    document.getElementById('login-section').classList.remove('hidden');
    document.getElementById('search-section').classList.add('hidden');
    document.getElementById('admin-section').classList.add('hidden');
    document.getElementById('login-btn').classList.remove('hidden');
    document.getElementById('logout-btn').classList.add('hidden');
    document.getElementById('admin-panel-btn').classList.add('hidden');
    document.getElementById('search-panel-btn').classList.add('hidden');
}

function showSearchSection() {
    document.getElementById('login-section').classList.add('hidden');
    document.getElementById('search-section').classList.remove('hidden');
    document.getElementById('admin-section').classList.add('hidden');
    document.getElementById('login-btn').classList.add('hidden');
    document.getElementById('logout-btn').classList.remove('hidden');

    // Show admin panel button if user is admin
    if (isAdmin) {
        document.getElementById('admin-panel-btn').classList.remove('hidden');
        document.getElementById('search-panel-btn').classList.add('hidden');
    }
}

function showAdminSection() {
    document.getElementById('login-section').classList.add('hidden');
    document.getElementById('search-section').classList.add('hidden');
    document.getElementById('admin-section').classList.remove('hidden');
    document.getElementById('login-btn').classList.add('hidden');
    document.getElementById('logout-btn').classList.remove('hidden');
    document.getElementById('admin-panel-btn').classList.add('hidden');
    document.getElementById('search-panel-btn').classList.remove('hidden');
}

// Login
async function handleLogin(e) {
    e.preventDefault();

    const key = document.getElementById('api-key').value;
    const errorDiv = document.getElementById('login-error');

    // Test the API key by hitting the verify endpoint
    try {
        const response = await fetch(`${API_BASE_URL}/verify`, {
            headers: {
                'Authorization': `Bearer ${key}`
            }
        });

        if (response.ok) {
            apiKey = key;
            localStorage.setItem('mkvmender_api_key', key);
            errorDiv.classList.add('hidden');
            await verifyAndShowApp();
        } else {
            throw new Error('Invalid API key');
        }
    } catch (error) {
        errorDiv.textContent = 'Login failed. Please check your API key.';
        errorDiv.classList.remove('hidden');
    }
}

// Logout
function handleLogout() {
    apiKey = '';
    userRole = 'user';
    isAdmin = false;
    localStorage.removeItem('mkvmender_api_key');
    localStorage.removeItem('mkvmender_user_role');
    document.getElementById('api-key').value = '';
    document.getElementById('results-container').classList.add('hidden');
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
        errorDiv.classList.add('hidden');
    } catch (error) {
        errorDiv.textContent = 'Search failed. Please try again.';
        errorDiv.classList.remove('hidden');
    }
}

// Display Results
function displayResults(data) {
    const resultsContainer = document.getElementById('results-container');
    const resultsTitle = document.getElementById('results-title');
    const resultsList = document.getElementById('results-list');

    if (!data.results || data.results.length === 0) {
        resultsTitle.textContent = 'No results found';
        resultsList.innerHTML = '<p class="text-text-secondary">Try a different search term or enable fuzzy matching.</p>';
        resultsContainer.classList.remove('hidden');
        return;
    }

    resultsTitle.textContent = `Found ${data.results.length} result(s) for "${data.query}"`;

    // Group results by title
    const grouped = groupResults(data.results);

    resultsList.innerHTML = grouped.map(group => {
        const icon = group.media_type === 'tv' ? '📺' : '🎬';
        const year = group.year ? ` (${group.year})` : '';

        return `
            <div class="card card-hover animate-fade-in">
                <div class="flex items-center gap-3 mb-3">
                    <span class="text-3xl">${icon}</span>
                    <div class="flex-1">
                        <h4 class="text-xl font-bold">${escapeHtml(group.title)}${year}</h4>
                        <span class="text-sm text-text-secondary uppercase">${group.media_type}</span>
                    </div>
                </div>

                <div class="flex flex-wrap gap-4 text-sm text-text-secondary mb-4">
                    <span class="flex items-center gap-1">📦 ${formatFileSize(group.file_size)}</span>
                    <span class="flex items-center gap-1">🔢 ${group.submissions.length} submission(s)</span>
                    ${group.season ? `<span>Season ${group.season}</span>` : ''}
                    ${group.episode ? `<span>Episode ${group.episode}</span>` : ''}
                </div>

                <div class="space-y-3 pl-4">
                    ${group.submissions.map(sub => `
                        <div class="bg-bg-tertiary p-4 rounded-lg border border-transparent hover:border-border transition-all">
                            <div class="flex flex-col md:flex-row md:items-start gap-3 mb-3">
                                <div class="flex-1 font-mono text-sm text-primary break-words">${escapeHtml(sub.filename)}</div>
                                <div class="flex gap-2 flex-shrink-0">
                                    <button class="vote-btn upvote px-3 py-2 bg-transparent border border-border rounded-lg hover:border-success hover:bg-success-bg hover:text-success transition-all text-sm flex items-center gap-2" data-submission-id="${sub.id}" data-vote-type="1">
                                        👍 <span class="vote-count font-semibold">${sub.upvotes}</span>
                                    </button>
                                    <button class="vote-btn downvote px-3 py-2 bg-transparent border border-border rounded-lg hover:border-error hover:bg-error-bg hover:text-error transition-all text-sm flex items-center gap-2" data-submission-id="${sub.id}" data-vote-type="-1">
                                        👎 <span class="vote-count font-semibold">${sub.downvotes}</span>
                                    </button>
                                </div>
                            </div>
                            <div class="flex gap-4 text-sm text-text-secondary">
                                <span>👤 ${escapeHtml(sub.username)}</span>
                                <span class="vote-score font-semibold">Score: ${sub.vote_score}</span>
                            </div>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
    }).join('');

    resultsContainer.classList.remove('hidden');

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
    if (type === 'success') {
        feedback.className = 'mt-2 px-3 py-2 bg-success-bg border border-success text-success rounded-lg text-sm';
    } else {
        feedback.className = 'mt-2 px-3 py-2 bg-error-bg border border-error text-error rounded-lg text-sm';
    }
    feedback.textContent = message;

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
