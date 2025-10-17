package database

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/quentinsteinke/mkvmender/internal/models"
)

// SortBy represents different sorting options for search results
type SortBy string

const (
	SortByRelevance SortBy = "relevance" // Fuzzy match score
	SortByVotes     SortBy = "votes"     // Highest votes first
	SortByDate      SortBy = "date"      // Newest first
	SortByTitle     SortBy = "title"     // Alphabetical
)

// SearchResult represents a search result with grouped submissions
type SearchResult struct {
	Title        string
	Year         *int
	MediaType    models.MediaType
	Season       *int
	Episode      *int
	Hash         string
	FileSize     int64
	Submissions  []models.SubmissionWithVotes
	FuzzyScore   int // Internal: fuzzy match score for sorting
}

// SearchOptions represents search parameters
type SearchOptions struct {
	Query     string
	SortBy    SortBy
	Limit     int
	UseFuzzy  bool
}

// SearchByTitle searches for submissions by title with fuzzy matching and sorting
func (db *DB) SearchByTitle(query string, sortBy SortBy, useFuzzy bool) ([]SearchResult, error) {
	// First, get all possible titles from the database
	titlesQuery := `
		SELECT DISTINCT nm.title
		FROM naming_metadata nm
		WHERE nm.title IS NOT NULL
	`

	rows, err := db.conn.Query(titlesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get titles: %w", err)
	}
	defer rows.Close()

	var allTitles []string
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			continue
		}
		allTitles = append(allTitles, title)
	}

	// Find matching titles using fuzzy matching
	var matchingTitles []string
	queryLower := strings.ToLower(query)

	if useFuzzy {
		// Use fuzzy matching
		matches := fuzzy.RankFindFold(query, allTitles)
		for _, match := range matches {
			matchingTitles = append(matchingTitles, match.Target)
		}
	} else {
		// Use simple LIKE matching
		for _, title := range allTitles {
			if strings.Contains(strings.ToLower(title), queryLower) {
				matchingTitles = append(matchingTitles, title)
			}
		}
	}

	if len(matchingTitles) == 0 {
		return []SearchResult{}, nil
	}

	// Build SQL query for matching titles
	placeholders := make([]string, len(matchingTitles))
	args := make([]interface{}, len(matchingTitles))
	for i, title := range matchingTitles {
		placeholders[i] = "?"
		args[i] = title
	}

	sqlQuery := fmt.Sprintf(`
		SELECT DISTINCT
			nm.title,
			nm.year,
			nm.season,
			nm.episode,
			fh.media_type,
			fh.hash,
			fh.file_size,
			ns.created_at
		FROM naming_metadata nm
		JOIN naming_submissions ns ON nm.submission_id = ns.id
		JOIN file_hashes fh ON ns.hash_id = fh.id
		WHERE nm.title IN (%s)
		ORDER BY nm.title, nm.year DESC, nm.season, nm.episode
	`, strings.Join(placeholders, ","))

	rows, err = db.conn.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer rows.Close()

	var results []SearchResult
	seen := make(map[string]bool) // To track unique hash combinations

	for rows.Next() {
		var r SearchResult
		var title *string
		var mediaTypeStr string
		var createdAt string

		err := rows.Scan(
			&title,
			&r.Year,
			&r.Season,
			&r.Episode,
			&mediaTypeStr,
			&r.Hash,
			&r.FileSize,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan result: %w", err)
		}

		if title != nil {
			r.Title = *title
		}
		r.MediaType = models.MediaType(mediaTypeStr)

		// Calculate fuzzy score for sorting
		if useFuzzy {
			r.FuzzyScore = fuzzy.RankMatchFold(query, r.Title)
		} else {
			// For non-fuzzy, score by position of match
			idx := strings.Index(strings.ToLower(r.Title), queryLower)
			if idx == 0 {
				r.FuzzyScore = 100 // Exact prefix match
			} else if idx > 0 {
				r.FuzzyScore = 50 - idx // Later matches get lower scores
			}
		}

		// Create unique key for deduplication
		key := fmt.Sprintf("%s-%v-%v-%v", r.Hash, r.Year, r.Season, r.Episode)
		if seen[key] {
			continue
		}
		seen[key] = true

		// Get all submissions for this hash
		submissions, err := db.GetSubmissionsByHash(r.Hash)
		if err != nil {
			continue
		}
		r.Submissions = submissions

		results = append(results, r)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	// Apply sorting
	sortResults(results, sortBy)

	return results, nil
}

// sortResults sorts search results based on the specified sort order
func sortResults(results []SearchResult, sortBy SortBy) {
	switch sortBy {
	case SortByRelevance:
		sort.Slice(results, func(i, j int) bool {
			// Higher fuzzy score = better match = comes first
			if results[i].FuzzyScore != results[j].FuzzyScore {
				return results[i].FuzzyScore > results[j].FuzzyScore
			}
			// Tie-breaker: alphabetical
			return results[i].Title < results[j].Title
		})

	case SortByVotes:
		sort.Slice(results, func(i, j int) bool {
			// Get max vote score from submissions
			maxI := getMaxVoteScore(results[i].Submissions)
			maxJ := getMaxVoteScore(results[j].Submissions)
			if maxI != maxJ {
				return maxI > maxJ
			}
			// Tie-breaker: alphabetical
			return results[i].Title < results[j].Title
		})

	case SortByDate:
		sort.Slice(results, func(i, j int) bool {
			// Get newest submission date
			newestI := getNewestDate(results[i].Submissions)
			newestJ := getNewestDate(results[j].Submissions)
			return newestI.After(newestJ)
		})

	case SortByTitle:
		sort.Slice(results, func(i, j int) bool {
			if results[i].Title != results[j].Title {
				return results[i].Title < results[j].Title
			}
			// Tie-breaker: year (descending)
			if results[i].Year != nil && results[j].Year != nil {
				return *results[i].Year > *results[j].Year
			}
			return false
		})
	}
}

// getMaxVoteScore returns the highest vote score from submissions
func getMaxVoteScore(submissions []models.SubmissionWithVotes) int {
	maxScore := 0
	for _, sub := range submissions {
		if sub.VoteScore > maxScore {
			maxScore = sub.VoteScore
		}
	}
	return maxScore
}

// getNewestDate returns the newest creation date from submissions
func getNewestDate(submissions []models.SubmissionWithVotes) time.Time {
	var newest time.Time
	for _, sub := range submissions {
		if sub.CreatedAt.After(newest) {
			newest = sub.CreatedAt
		}
	}
	return newest
}

// GroupedTVShow represents a TV show with organized seasons and episodes
type GroupedTVShow struct {
	Title   string
	Year    *int
	Seasons map[int][]SearchResult // season number -> episodes
}

// GroupTVShowResults organizes search results for a TV show by season
func GroupTVShowResults(results []SearchResult) []GroupedTVShow {
	shows := make(map[string]*GroupedTVShow) // key: title-year

	for _, result := range results {
		if result.MediaType != models.MediaTypeTV {
			continue
		}

		key := result.Title
		if result.Year != nil {
			key = fmt.Sprintf("%s-%d", result.Title, *result.Year)
		}

		if _, exists := shows[key]; !exists {
			shows[key] = &GroupedTVShow{
				Title:   result.Title,
				Year:    result.Year,
				Seasons: make(map[int][]SearchResult),
			}
		}

		if result.Season != nil {
			shows[key].Seasons[*result.Season] = append(shows[key].Seasons[*result.Season], result)
		}
	}

	// Convert map to slice
	var grouped []GroupedTVShow
	for _, show := range shows {
		grouped = append(grouped, *show)
	}

	return grouped
}
