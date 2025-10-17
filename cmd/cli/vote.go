package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/quentinsteinke/mkvmender/internal/api"
	"github.com/quentinsteinke/mkvmender/internal/hasher"
	"github.com/quentinsteinke/mkvmender/internal/models"
	"github.com/spf13/cobra"
)

func newVoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote <file>",
		Short: "Vote on naming submissions for a file",
		Long:  "Interactively view and vote on naming submissions for a media file.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]

			// Hash the file
			fmt.Println("Hashing file...")
			result, err := hasher.HashFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to hash file: %w", err)
			}

			// Create API client
			client, err := api.NewClient()
			if err != nil {
				return fmt.Errorf("failed to create API client: %w", err)
			}

			// Lookup naming options
			fmt.Println("Looking up naming options...")
			response, err := client.Lookup(result.Hash)
			if err != nil {
				return fmt.Errorf("lookup failed: %w", err)
			}

			if len(response.Submissions) == 0 {
				fmt.Println("No naming submissions found for this file.")
				return nil
			}

			// Display options
			fmt.Printf("\nFound %d naming option(s):\n\n", len(response.Submissions))
			for i, submission := range response.Submissions {
				fmt.Printf("[%d] %s\n", i+1, submission.Filename)
				fmt.Printf("    Submitted by: %s\n", submission.Username)
				fmt.Printf("    Votes: %d (↑%d ↓%d)\n", submission.VoteScore, submission.Upvotes, submission.Downvotes)
				if submission.Metadata != nil && submission.Metadata.Title != nil {
					fmt.Printf("    Title: %s", *submission.Metadata.Title)
					if submission.Metadata.Year != nil {
						fmt.Printf(" (%d)", *submission.Metadata.Year)
					}
					fmt.Println()
				}
				fmt.Println()
			}

			// Prompt for selection
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Select a submission to vote on (1-" + fmt.Sprint(len(response.Submissions)) + ") or 'q' to quit: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "q" || input == "" {
				fmt.Println("Cancelled.")
				return nil
			}

			selection, err := strconv.Atoi(input)
			if err != nil || selection < 1 || selection > len(response.Submissions) {
				return fmt.Errorf("invalid selection")
			}

			selectedSubmission := response.Submissions[selection-1]

			// Ask for vote type
			fmt.Print("\nUpvote or downvote? (up/down): ")
			voteInput, _ := reader.ReadString('\n')
			voteInput = strings.TrimSpace(strings.ToLower(voteInput))

			var voteType models.VoteType
			switch voteInput {
			case "up", "upvote", "u", "1":
				voteType = models.VoteUp
			case "down", "downvote", "d", "-1":
				voteType = models.VoteDown
			default:
				return fmt.Errorf("invalid vote type (use 'up' or 'down')")
			}

			// Submit vote
			if err := client.Vote(selectedSubmission.ID, voteType); err != nil {
				return fmt.Errorf("vote failed: %w", err)
			}

			voteAction := "upvoted"
			if voteType == models.VoteDown {
				voteAction = "downvoted"
			}

			fmt.Printf("\n✓ Successfully %s: %s\n", voteAction, selectedSubmission.Filename)

			// Show updated results
			fmt.Println("\nFetching updated vote counts...")
			updatedResponse, err := client.Lookup(result.Hash)
			if err == nil && len(updatedResponse.Submissions) > 0 {
				fmt.Println("\nUpdated rankings:")
				for i, submission := range updatedResponse.Submissions {
					prefix := "  "
					if submission.ID == selectedSubmission.ID {
						prefix = "→ "
					}
					fmt.Printf("%s[%d] %s - Votes: %d (↑%d ↓%d)\n",
						prefix, i+1, submission.Filename,
						submission.VoteScore, submission.Upvotes, submission.Downvotes)
				}
			}

			return nil
		},
	}

	return cmd
}
