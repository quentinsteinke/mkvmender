package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/quentinsteinke/mkvmender/internal/api"
	"github.com/quentinsteinke/mkvmender/internal/hasher"
	"github.com/quentinsteinke/mkvmender/internal/models"
	"github.com/spf13/cobra"
)

func newSearchCmd() *cobra.Command {
	var sortBy string
	var noFuzzy bool

	cmd := &cobra.Command{
		Use:   "search <title>",
		Short: "Search for movies and TV shows by title",
		Long:  "Search the database for movies and TV shows, browse seasons and episodes, and view naming submissions.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := strings.Join(args, " ")

			// Create API client
			client, err := api.NewClient()
			if err != nil {
				return fmt.Errorf("failed to create API client: %w", err)
			}

			// Search
			fmt.Printf("Searching for '%s'...\n\n", query)
			response, err := client.Search(query, sortBy, !noFuzzy)
			if err != nil {
				return fmt.Errorf("search failed: %w", err)
			}

			if len(response.Results) == 0 {
				fmt.Println("No results found.")
				return nil
			}

			// Group results by title and media type
			grouped := groupSearchResults(response.Results)

			// Display grouped results
			fmt.Printf("Found %d result(s):\n\n", len(grouped))
			for i, group := range grouped {
				yearStr := ""
				if group.Year != nil {
					yearStr = fmt.Sprintf(" (%d)", *group.Year)
				}
				mediaIcon := "🎬"
				if group.MediaType == models.MediaTypeTV {
					mediaIcon = "📺"
				}
				fmt.Printf("[%d] %s %s%s - %s\n", i+1, mediaIcon, group.Title, yearStr, group.MediaType)
			}

			// Prompt for selection
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("\nSelect a title (1-" + fmt.Sprint(len(grouped)) + ") or 'q' to quit: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "q" || input == "" {
				return nil
			}

			selection, err := strconv.Atoi(input)
			if err != nil || selection < 1 || selection > len(grouped) {
				return fmt.Errorf("invalid selection")
			}

			selectedGroup := grouped[selection-1]

			// Handle based on media type
			if selectedGroup.MediaType == models.MediaTypeMovie {
				return displayMovieSubmissions(selectedGroup)
			} else {
				return displayTVShowSeasons(selectedGroup, reader)
			}
		},
	}

	cmd.Flags().StringVarP(&sortBy, "sort", "s", "relevance", "Sort results by: relevance, votes, date, title")
	cmd.Flags().BoolVar(&noFuzzy, "no-fuzzy", false, "Disable fuzzy matching (use exact string matching)")

	return cmd
}

type searchGroup struct {
	Title     string
	Year      *int
	MediaType models.MediaType
	Results   []models.SearchResult
}

func groupSearchResults(results []models.SearchResult) []searchGroup {
	groups := make(map[string]*searchGroup)

	for _, result := range results {
		key := result.Title
		if result.Year != nil {
			key = fmt.Sprintf("%s-%d", result.Title, *result.Year)
		}
		key = fmt.Sprintf("%s-%s", key, result.MediaType)

		if _, exists := groups[key]; !exists {
			groups[key] = &searchGroup{
				Title:     result.Title,
				Year:      result.Year,
				MediaType: result.MediaType,
				Results:   []models.SearchResult{},
			}
		}

		groups[key].Results = append(groups[key].Results, result)
	}

	// Convert to slice and sort
	var grouped []searchGroup
	for _, group := range groups {
		grouped = append(grouped, *group)
	}

	sort.Slice(grouped, func(i, j int) bool {
		if grouped[i].Title != grouped[j].Title {
			return grouped[i].Title < grouped[j].Title
		}
		if grouped[i].Year != nil && grouped[j].Year != nil {
			return *grouped[i].Year > *grouped[j].Year
		}
		return grouped[i].MediaType < grouped[j].MediaType
	})

	return grouped
}

func displayMovieSubmissions(group searchGroup) error {
	if len(group.Results) == 0 {
		fmt.Println("No submissions found.")
		return nil
	}

	// Movies should have just one result with all submissions
	result := group.Results[0]

	yearStr := ""
	if result.Year != nil {
		yearStr = fmt.Sprintf(" (%d)", *result.Year)
	}

	fmt.Printf("\n%s%s\n", result.Title, yearStr)
	fmt.Printf("Hash: %s\n", result.Hash)
	fmt.Printf("Size: %s\n\n", hasher.FormatFileSize(result.FileSize))

	if len(result.Submissions) == 0 {
		fmt.Println("No naming submissions found.")
		return nil
	}

	fmt.Printf("Found %d naming submission(s):\n\n", len(result.Submissions))
	for i, submission := range result.Submissions {
		fmt.Printf("[%d] %s\n", i+1, submission.Filename)
		fmt.Printf("    Submitted by: %s\n", submission.Username)
		fmt.Printf("    Votes: %d (↑%d ↓%d)\n", submission.VoteScore, submission.Upvotes, submission.Downvotes)
		if submission.Metadata != nil {
			if submission.Metadata.Quality != nil {
				fmt.Printf("    Quality: %s\n", *submission.Metadata.Quality)
			}
			if submission.Metadata.Source != nil {
				fmt.Printf("    Source: %s\n", *submission.Metadata.Source)
			}
		}
		fmt.Println()
	}

	return nil
}

func displayTVShowSeasons(group searchGroup, reader *bufio.Reader) error {
	// Group results by season
	seasons := make(map[int][]models.SearchResult)
	for _, result := range group.Results {
		if result.Season != nil {
			seasons[*result.Season] = append(seasons[*result.Season], result)
		}
	}

	if len(seasons) == 0 {
		fmt.Println("No seasons found.")
		return nil
	}

	// Sort seasons
	var seasonNumbers []int
	for season := range seasons {
		seasonNumbers = append(seasonNumbers, season)
	}
	sort.Ints(seasonNumbers)

	yearStr := ""
	if group.Year != nil {
		yearStr = fmt.Sprintf(" (%d)", *group.Year)
	}

	fmt.Printf("\n%s%s - TV Show\n\n", group.Title, yearStr)
	fmt.Printf("Found %d season(s):\n\n", len(seasons))
	for i, seasonNum := range seasonNumbers {
		episodeCount := len(seasons[seasonNum])
		fmt.Printf("[%d] Season %d (%d episode(s))\n", i+1, seasonNum, episodeCount)
	}

	// Prompt for season selection
	fmt.Print("\nSelect a season (1-" + fmt.Sprint(len(seasonNumbers)) + ") or 'q' to quit: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "q" || input == "" {
		return nil
	}

	selection, err := strconv.Atoi(input)
	if err != nil || selection < 1 || selection > len(seasonNumbers) {
		return fmt.Errorf("invalid selection")
	}

	selectedSeason := seasonNumbers[selection-1]
	return displayTVShowEpisodes(group.Title, selectedSeason, seasons[selectedSeason], reader)
}

func displayTVShowEpisodes(title string, seasonNum int, episodes []models.SearchResult, reader *bufio.Reader) error {
	// Sort episodes by episode number
	sort.Slice(episodes, func(i, j int) bool {
		if episodes[i].Episode != nil && episodes[j].Episode != nil {
			return *episodes[i].Episode < *episodes[j].Episode
		}
		return false
	})

	fmt.Printf("\n%s - Season %d\n\n", title, seasonNum)
	fmt.Printf("Found %d episode(s):\n\n", len(episodes))

	for i, episode := range episodes {
		epNum := "?"
		if episode.Episode != nil {
			epNum = fmt.Sprint(*episode.Episode)
		}
		submissionCount := len(episode.Submissions)
		fmt.Printf("[%d] Episode %s (%d submission(s))\n", i+1, epNum, submissionCount)
	}

	// Prompt for episode selection
	fmt.Print("\nSelect an episode (1-" + fmt.Sprint(len(episodes)) + ") or 'q' to quit: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "q" || input == "" {
		return nil
	}

	selection, err := strconv.Atoi(input)
	if err != nil || selection < 1 || selection > len(episodes) {
		return fmt.Errorf("invalid selection")
	}

	selectedEpisode := episodes[selection-1]
	return displayEpisodeSubmissions(title, seasonNum, selectedEpisode)
}

func displayEpisodeSubmissions(title string, seasonNum int, episode models.SearchResult) error {
	epNum := "?"
	if episode.Episode != nil {
		epNum = fmt.Sprint(*episode.Episode)
	}

	fmt.Printf("\n%s - S%02dE%s\n", title, seasonNum, epNum)
	fmt.Printf("Hash: %s\n", episode.Hash)
	fmt.Printf("Size: %s\n\n", hasher.FormatFileSize(episode.FileSize))

	if len(episode.Submissions) == 0 {
		fmt.Println("No naming submissions found.")
		return nil
	}

	fmt.Printf("Found %d naming submission(s):\n\n", len(episode.Submissions))
	for i, submission := range episode.Submissions {
		fmt.Printf("[%d] %s\n", i+1, submission.Filename)
		fmt.Printf("    Submitted by: %s\n", submission.Username)
		fmt.Printf("    Votes: %d (↑%d ↓%d)\n", submission.VoteScore, submission.Upvotes, submission.Downvotes)
		if submission.Metadata != nil {
			if submission.Metadata.Quality != nil {
				fmt.Printf("    Quality: %s\n", *submission.Metadata.Quality)
			}
			if submission.Metadata.Source != nil {
				fmt.Printf("    Source: %s\n", *submission.Metadata.Source)
			}
		}
		fmt.Println()
	}

	return nil
}
