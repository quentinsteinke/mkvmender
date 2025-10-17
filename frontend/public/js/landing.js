// GitHub API Configuration
const GITHUB_REPO = 'quentinsteinke/mkvmender';
const GITHUB_API_URL = `https://api.github.com/repos/${GITHUB_REPO}`;

// Cache for GitHub data (15 minutes)
const CACHE_DURATION = 15 * 60 * 1000;
let githubDataCache = null;
let cacheTimestamp = 0;

// Initialize on DOM load
document.addEventListener('DOMContentLoaded', () => {
    setupDownloadButtons();
    setupMobileMenu();
    setupTerminalAnimation();
    fetchGitHubStats();
    setupCopyButton();
    setupScrollAnimations();
});

// OS Detection and Download Button Setup
function setupDownloadButtons() {
    const buttons = [
        {
            btn: document.getElementById('hero-download-btn'),
            text: document.getElementById('hero-download-text')
        },
        {
            btn: document.getElementById('cta-download-btn'),
            text: document.getElementById('cta-download-text')
        }
    ];

    const userAgent = navigator.userAgent.toLowerCase();
    const platform = navigator.platform.toLowerCase();

    let downloadUrl = '';
    let label = 'Download';

    // Detect OS
    if (platform.includes('mac') || userAgent.includes('mac')) {
        // Default to Apple Silicon (most common since 2020)
        downloadUrl = 'https://github.com/quentinsteinke/mkvmender/releases/download/v1.0.0/mkvmender-darwin-arm64';
        label = 'Download for macOS';
    } else if (platform.includes('win') || userAgent.includes('win')) {
        downloadUrl = 'https://github.com/quentinsteinke/mkvmender/releases/download/v1.0.0/mkvmender-windows-amd64.exe';
        label = 'Download for Windows';
    } else if (platform.includes('linux') || userAgent.includes('linux')) {
        const isArm = userAgent.includes('arm') || platform.includes('arm') || userAgent.includes('aarch64');
        if (isArm) {
            downloadUrl = 'https://github.com/quentinsteinke/mkvmender/releases/download/v1.0.0/mkvmender-linux-arm64';
            label = 'Download for Linux (ARM64)';
        } else {
            downloadUrl = 'https://github.com/quentinsteinke/mkvmender/releases/download/v1.0.0/mkvmender-linux-amd64';
            label = 'Download for Linux (x86_64)';
        }
    }

    // Set button properties
    buttons.forEach(({ btn, text }) => {
        if (downloadUrl && btn && text) {
            btn.href = downloadUrl;
            text.textContent = label;
        } else if (btn) {
            // Fallback to releases page
            btn.href = 'https://github.com/quentinsteinke/mkvmender/releases';
            btn.target = '_blank';
            if (text) text.textContent = 'View Releases';
        }
    });
}

// Mobile Menu Toggle
function setupMobileMenu() {
    const menuBtn = document.getElementById('mobile-menu-btn');
    const menu = document.getElementById('mobile-menu');

    if (menuBtn && menu) {
        menuBtn.addEventListener('click', () => {
            menu.classList.toggle('hidden');
        });

        // Close menu when clicking on a link
        menu.querySelectorAll('a').forEach(link => {
            link.addEventListener('click', () => {
                menu.classList.add('hidden');
            });
        });
    }
}

// Terminal Typing Animation
function setupTerminalAnimation() {
    const terminal = document.getElementById('terminal-content');
    if (!terminal) return;

    const scenarios = [
        {
            lines: [
                { text: '$ mkvmender register', delay: 0 },
                { text: 'Enter username: neo', delay: 800, type: 'input' },
                { text: '✓ Account created successfully!', delay: 600, type: 'success' },
                { text: 'Your API key: eyJhbGc...', delay: 400, type: 'success' },
                { text: 'Save this key - you\'ll need it to login!', delay: 300, type: 'secondary' }
            ]
        },
        {
            lines: [
                { text: '$ mkvmender search "the matrix"', delay: 0 },
                { text: 'Searching database...', delay: 600, type: 'muted' },
                { text: '[1] 🎬 The Matrix (1999)', delay: 500, type: 'primary' },
                { text: '    42 community submissions', delay: 300, type: 'secondary' },
                { text: '[2] 📺 The Matrix (TV Series)', delay: 400, type: 'primary' },
                { text: '    8 community submissions', delay: 300, type: 'secondary' }
            ]
        },
        {
            lines: [
                { text: '$ mkvmender lookup movie.mkv', delay: 0 },
                { text: 'Calculating hash... ████████████ 100%', delay: 800, type: 'muted' },
                { text: 'Hash: a3b2c1d4e5f6...', delay: 400, type: 'muted' },
                { text: 'Found 3 naming submissions:', delay: 600, type: 'primary' },
                { text: '[1] The.Matrix.1999.1080p.mkv', delay: 400, type: 'success' },
                { text: '    👤 moviefan23  |  👍 42  👎 3', delay: 300, type: 'secondary' }
            ]
        }
    ];

    let currentScenario = 0;

    function runAnimation() {
        terminal.innerHTML = '';
        const scenario = scenarios[currentScenario];
        let delay = 0;

        scenario.lines.forEach((line, index) => {
            delay += line.delay;

            setTimeout(() => {
                const lineEl = document.createElement('div');
                lineEl.style.opacity = '0';
                lineEl.style.transform = 'translateY(5px)';

                // Apply styling based on type
                switch (line.type) {
                    case 'input':
                        lineEl.className = 'text-text-secondary ml-4';
                        break;
                    case 'success':
                        lineEl.className = 'text-success ml-4';
                        break;
                    case 'primary':
                        lineEl.className = 'text-primary ml-4';
                        break;
                    case 'secondary':
                        lineEl.className = 'text-text-secondary ml-8 text-sm';
                        break;
                    case 'muted':
                        lineEl.className = 'text-text-muted ml-4';
                        break;
                    default:
                        lineEl.className = 'text-text';
                }

                lineEl.textContent = line.text;
                terminal.appendChild(lineEl);

                // Fade in animation
                setTimeout(() => {
                    lineEl.style.transition = 'all 0.3s ease';
                    lineEl.style.opacity = '1';
                    lineEl.style.transform = 'translateY(0)';

                    // Auto-scroll to bottom
                    terminal.scrollTop = terminal.scrollHeight;
                }, 10);

                // Add cursor on last line
                if (index === scenario.lines.length - 1) {
                    setTimeout(() => {
                        const cursorEl = document.createElement('div');
                        cursorEl.className = 'text-text mt-4';
                        cursorEl.innerHTML = '<span class="text-success">$</span> <span class="animate-pulse">_</span>';
                        terminal.appendChild(cursorEl);
                    }, 500);
                }
            }, delay);
        });

        // Move to next scenario after animation completes
        const totalDelay = scenario.lines.reduce((sum, line) => sum + line.delay, 0);
        setTimeout(() => {
            currentScenario = (currentScenario + 1) % scenarios.length;
            setTimeout(runAnimation, 2000); // Wait 2 seconds before next scenario
        }, totalDelay + 3000); // Show result for 3 seconds
    }

    runAnimation();
}

// Fetch GitHub Statistics
async function fetchGitHubStats() {
    // Check cache first
    const now = Date.now();
    if (githubDataCache && (now - cacheTimestamp) < CACHE_DURATION) {
        updateStatsUI(githubDataCache);
        return;
    }

    try {
        // Fetch repository data
        const repoResponse = await fetch(GITHUB_API_URL);
        if (!repoResponse.ok) throw new Error('Failed to fetch repo data');
        const repoData = await repoResponse.json();

        // Fetch release data
        const releasesResponse = await fetch(`${GITHUB_API_URL}/releases/latest`);
        let downloadCount = 0;
        if (releasesResponse.ok) {
            const releaseData = await releasesResponse.json();
            downloadCount = releaseData.assets?.reduce((sum, asset) => sum + (asset.download_count || 0), 0) || 0;
        }

        // Fetch contributors
        const contributorsResponse = await fetch(`${GITHUB_API_URL}/contributors`);
        let contributorCount = 0;
        if (contributorsResponse.ok) {
            const contributors = await contributorsResponse.json();
            contributorCount = contributors.length;
        }

        const data = {
            stars: repoData.stargazers_count || 0,
            downloads: downloadCount,
            contributors: contributorCount
        };

        // Update cache
        githubDataCache = data;
        cacheTimestamp = now;

        updateStatsUI(data);
    } catch (error) {
        console.error('Error fetching GitHub stats:', error);
        // Set placeholder values on error
        updateStatsUI({
            stars: '-',
            downloads: '-',
            contributors: '-'
        });
    }
}

// Update Stats UI with Animation
function updateStatsUI(data) {
    // Animate numbers
    animateNumber('stat-stars', data.stars);
    animateNumber('stat-downloads', data.downloads);
    animateNumber('stat-contributors', data.contributors);

    // Update GitHub stars in nav
    const githubStars = document.getElementById('github-stars');
    if (githubStars && data.stars !== '-') {
        githubStars.textContent = `${formatNumber(data.stars)} stars`;
    }
}

// Animate Number Counter
function animateNumber(elementId, target) {
    const element = document.getElementById(elementId);
    if (!element) return;

    if (target === '-') {
        element.textContent = '-';
        return;
    }

    const duration = 1000;
    const start = 0;
    const increment = target / (duration / 16); // 60fps
    let current = start;

    const timer = setInterval(() => {
        current += increment;
        if (current >= target) {
            current = target;
            clearInterval(timer);
        }
        element.textContent = formatNumber(Math.floor(current));
    }, 16);
}

// Format Number with K suffix
function formatNumber(num) {
    if (num >= 1000) {
        return (num / 1000).toFixed(1) + 'K';
    }
    return num.toString();
}

// Copy to Clipboard
function setupCopyButton() {
    const copyBtn = document.getElementById('copy-install-btn');
    if (!copyBtn) return;

    copyBtn.addEventListener('click', async () => {
        const installCommand = 'curl -sSL https://raw.githubusercontent.com/quentinsteinke/mkvmender/main/install.sh | bash';

        try {
            await navigator.clipboard.writeText(installCommand);

            // Show feedback
            const originalHTML = copyBtn.innerHTML;
            copyBtn.innerHTML = `
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                </svg>
            `;
            copyBtn.classList.add('!border-success', 'text-success');

            setTimeout(() => {
                copyBtn.innerHTML = originalHTML;
                copyBtn.classList.remove('!border-success', 'text-success');
            }, 2000);
        } catch (error) {
            console.error('Failed to copy:', error);
            // Fallback: show alert
            alert('Copy failed. Command: ' + installCommand);
        }
    });
}

// Scroll Animations (Intersection Observer)
function setupScrollAnimations() {
    const observerOptions = {
        threshold: 0.1,
        rootMargin: '0px 0px -50px 0px'
    };

    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('animate-fade-in');
                observer.unobserve(entry.target);
            }
        });
    }, observerOptions);

    // Observe all cards and sections
    document.querySelectorAll('.card, section > div > div').forEach(el => {
        observer.observe(el);
    });
}
