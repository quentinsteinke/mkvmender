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

    // Static install command for hero section
    terminal.textContent = 'curl -sSL https://raw.githubusercontent.com/quentinsteinke/mkvmender/main/install.sh | bash';
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
