// Documentation JavaScript
// Handles search, TOC generation, scroll spy, and copy buttons

// Search index
let searchIndex = [];
let searchDebounceTimer;

// Initialize on DOM load
document.addEventListener('DOMContentLoaded', () => {
    generateTOC();
    setupSearch();
    setupScrollSpy();
    setupCopyButtons();
    setupMobileMenu();
    setupKeyboardShortcuts();
    buildSearchIndex();
});

// Generate Table of Contents from headings
function generateTOC() {
    const content = document.querySelector('.docs-content');
    const toc = document.getElementById('docs-toc');
    const tocMobile = document.getElementById('docs-toc-mobile');

    if (!content || !toc) return;

    // Find all section headings (h2)
    const headings = content.querySelectorAll('h2');

    headings.forEach(heading => {
        const section = heading.closest('section');
        if (!section) return;

        const id = section.id;
        const text = heading.textContent;

        // Create TOC link
        const link = document.createElement('a');
        link.href = `#${id}`;
        link.textContent = text;
        link.className = 'block py-2 px-3 rounded-lg text-sm text-text-secondary hover:text-text hover:bg-bg-hover transition-colors toc-link';
        link.dataset.section = id;

        // Add to both desktop and mobile TOC
        toc.appendChild(link.cloneNode(true));
        if (tocMobile) tocMobile.appendChild(link);

        // Smooth scroll on click
        link.addEventListener('click', (e) => {
            e.preventDefault();
            scrollToSection(id);
        });
    });
}

// Scroll to section with offset for sticky header
function scrollToSection(id) {
    const section = document.getElementById(id);
    if (!section) return;

    const headerHeight = 80; // Height of sticky header
    const elementPosition = section.getBoundingClientRect().top;
    const offsetPosition = elementPosition + window.pageYOffset - headerHeight;

    window.scrollTo({
        top: offsetPosition,
        behavior: 'smooth'
    });

    // Close mobile TOC if open
    document.getElementById('mobile-toc')?.classList.add('hidden');

    // Update URL without scrolling
    history.pushState(null, null, `#${id}`);
}

// Scroll spy - highlight active section in TOC
function setupScrollSpy() {
    const sections = document.querySelectorAll('.docs-section');
    const tocLinks = document.querySelectorAll('.toc-link');

    if (sections.length === 0 || tocLinks.length === 0) return;

    const observerOptions = {
        root: null,
        rootMargin: '-100px 0px -66%',
        threshold: 0
    };

    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const id = entry.target.id;

                // Remove active class from all links
                tocLinks.forEach(link => {
                    link.classList.remove('!bg-primary/10', '!text-primary', '!border-l-2', '!border-primary');
                });

                // Add active class to current link
                const activeLink = document.querySelector(`.toc-link[data-section="${id}"]`);
                if (activeLink) {
                    activeLink.classList.add('!bg-primary/10', '!text-primary', '!border-l-2', '!border-primary');
                }
            }
        });
    }, observerOptions);

    sections.forEach(section => observer.observe(section));
}

// Build search index from all content
function buildSearchIndex() {
    const sections = document.querySelectorAll('.docs-section');

    sections.forEach(section => {
        const id = section.id;
        const heading = section.querySelector('h2');
        const title = heading ? heading.textContent : '';

        // Get all text content, excluding code blocks initially
        const content = Array.from(section.querySelectorAll('p, li, h3'))
            .map(el => el.textContent)
            .join(' ');

        searchIndex.push({
            id,
            title,
            content: content.toLowerCase(),
            preview: content.slice(0, 150) + '...'
        });
    });
}

// Setup search functionality
function setupSearch() {
    const searchInputs = [
        document.getElementById('docs-search'),
        document.getElementById('docs-search-mobile')
    ];
    const searchResults = [
        document.getElementById('search-results'),
        document.getElementById('search-results-mobile')
    ];

    searchInputs.forEach((input, index) => {
        if (!input) return;

        input.addEventListener('input', (e) => {
            clearTimeout(searchDebounceTimer);
            const query = e.target.value.trim();

            if (query.length < 2) {
                searchResults[index]?.classList.add('hidden');
                return;
            }

            // Debounce search
            searchDebounceTimer = setTimeout(() => {
                performSearch(query, searchResults[index]);
            }, 300);
        });

        // Close results when clicking outside
        document.addEventListener('click', (e) => {
            if (!input.contains(e.target) && !searchResults[index]?.contains(e.target)) {
                searchResults[index]?.classList.add('hidden');
            }
        });
    });
}

// Perform fuzzy search
function performSearch(query, resultsContainer) {
    if (!resultsContainer) return;

    const queryLower = query.toLowerCase();
    const results = [];

    // Simple fuzzy search - check if query terms are in content
    searchIndex.forEach(item => {
        const titleMatch = item.title.toLowerCase().includes(queryLower);
        const contentMatch = item.content.includes(queryLower);

        if (titleMatch || contentMatch) {
            // Calculate relevance score
            let score = 0;
            if (titleMatch) score += 10;
            if (contentMatch) score += 1;

            results.push({
                ...item,
                score,
                titleMatch
            });
        }
    });

    // Sort by relevance
    results.sort((a, b) => b.score - a.score);

    // Display results
    if (results.length === 0) {
        resultsContainer.innerHTML = '<div class="p-4 text-text-secondary text-sm">No results found</div>';
        resultsContainer.classList.remove('hidden');
        return;
    }

    const html = results.slice(0, 5).map(result => `
        <a href="#${result.id}"
           class="block p-4 hover:bg-bg-hover border-b border-border last:border-b-0 transition-colors"
           onclick="document.getElementById('search-results')?.classList.add('hidden'); document.getElementById('search-results-mobile')?.classList.add('hidden');">
            <div class="font-medium text-text mb-1">${highlightMatch(result.title, query)}</div>
            <div class="text-sm text-text-secondary line-clamp-2">${highlightMatch(result.preview, query)}</div>
        </a>
    `).join('');

    resultsContainer.innerHTML = html;
    resultsContainer.classList.remove('hidden');
}

// Highlight search query in text
function highlightMatch(text, query) {
    const regex = new RegExp(`(${escapeRegex(query)})`, 'gi');
    return text.replace(regex, '<mark class="bg-primary/30 text-text">$1</mark>');
}

// Escape special regex characters
function escapeRegex(string) {
    return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

// Add copy buttons to all code blocks
function setupCopyButtons() {
    const codeBlocks = document.querySelectorAll('.terminal');

    codeBlocks.forEach(block => {
        // Create copy button
        const button = document.createElement('button');
        button.className = 'absolute top-3 right-3 p-2 bg-bg-hover rounded-lg border border-border hover:border-primary transition-colors opacity-0 group-hover:opacity-100';
        button.innerHTML = `
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
            </svg>
        `;

        // Make block relative and add group class for hover effect
        block.classList.add('relative', 'group');
        block.appendChild(button);

        // Copy functionality
        button.addEventListener('click', async () => {
            const code = Array.from(block.querySelectorAll('div'))
                .map(div => div.textContent)
                .filter(text => !text.startsWith('#')) // Remove comments
                .map(text => text.replace(/^\$ /, '')) // Remove $ prompt
                .join('\n')
                .trim();

            try {
                await navigator.clipboard.writeText(code);

                // Show success feedback
                const originalHTML = button.innerHTML;
                button.innerHTML = `
                    <svg class="w-4 h-4 text-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                    </svg>
                `;
                button.classList.add('!border-success');

                setTimeout(() => {
                    button.innerHTML = originalHTML;
                    button.classList.remove('!border-success');
                }, 2000);
            } catch (error) {
                console.error('Failed to copy:', error);
            }
        });
    });
}

// Mobile menu functionality
function setupMobileMenu() {
    const menuBtn = document.getElementById('mobile-menu-btn');
    const tocBtn = document.getElementById('mobile-toc-btn');
    const toc = document.getElementById('mobile-toc');
    const tocOverlay = document.getElementById('mobile-toc-overlay');
    const tocClose = document.getElementById('mobile-toc-close');

    // Mobile menu button (if exists on this page)
    if (menuBtn) {
        menuBtn.addEventListener('click', () => {
            // Toggle some menu if needed
        });
    }

    // Mobile TOC toggle
    if (tocBtn && toc) {
        tocBtn.addEventListener('click', () => {
            toc.classList.remove('hidden');
        });
    }

    // Close TOC
    [tocOverlay, tocClose].forEach(el => {
        if (el) {
            el.addEventListener('click', () => {
                toc?.classList.add('hidden');
            });
        }
    });

    // Close TOC when clicking a link
    const tocLinks = document.querySelectorAll('#docs-toc-mobile a');
    tocLinks.forEach(link => {
        link.addEventListener('click', () => {
            toc?.classList.add('hidden');
        });
    });
}

// Keyboard shortcuts
function setupKeyboardShortcuts() {
    document.addEventListener('keydown', (e) => {
        // Cmd/Ctrl + K to focus search
        if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
            e.preventDefault();
            const searchInput = document.getElementById('docs-search') || document.getElementById('docs-search-mobile');
            searchInput?.focus();
        }

        // Escape to close search results
        if (e.key === 'Escape') {
            document.getElementById('search-results')?.classList.add('hidden');
            document.getElementById('search-results-mobile')?.classList.add('hidden');
            document.getElementById('mobile-toc')?.classList.add('hidden');
        }
    });
}
