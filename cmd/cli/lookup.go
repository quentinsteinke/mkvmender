package main

import (
	"fmt"

	"github.com/quentinsteinke/mkvmender/internal/api"
	"github.com/quentinsteinke/mkvmender/internal/hasher"
	"github.com/spf13/cobra"
)

func newLookupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lookup <file>",
		Short: "Look up naming options for a media file",
		Long:  "Hashes the file and looks up available naming submissions from the database.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]

			// Hash the file
			fmt.Println("Hashing file...")
			result, err := hasher.HashFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to hash file: %w", err)
			}

			fmt.Printf("Hash: %s\n", result.Hash)
			fmt.Printf("Size: %s\n\n", hasher.FormatFileSize(result.FileSize))

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
				fmt.Println("Consider uploading your own naming using 'mkvmender upload'")
				return nil
			}

			fmt.Printf("Found %d naming option(s):\n\n", len(response.Submissions))
			for i, submission := range response.Submissions {
				fmt.Printf("[%d] %s\n", i+1, submission.Filename)
				fmt.Printf("    Submitted by: %s\n", submission.Username)
				fmt.Printf("    Votes: %d (↑%d ↓%d)\n", submission.VoteScore, submission.Upvotes, submission.Downvotes)
				fmt.Printf("    Media Type: %s\n", submission.MediaType)

				if submission.Metadata != nil {
					if submission.Metadata.Title != nil {
						fmt.Printf("    Title: %s\n", *submission.Metadata.Title)
					}
					if submission.Metadata.Year != nil {
						fmt.Printf("    Year: %d\n", *submission.Metadata.Year)
					}
					if submission.Metadata.Season != nil && submission.Metadata.Episode != nil {
						fmt.Printf("    Episode: S%02dE%02d\n", *submission.Metadata.Season, *submission.Metadata.Episode)
					}
					if submission.Metadata.Quality != nil {
						fmt.Printf("    Quality: %s\n", *submission.Metadata.Quality)
					}
				}
				fmt.Println()
			}

			return nil
		},
	}

	return cmd
}
