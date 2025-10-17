package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/quentinsteinke/mkvmender/internal/api"
	"github.com/quentinsteinke/mkvmender/internal/hasher"
	"github.com/spf13/cobra"
)

func newRenameCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "rename <file>",
		Short: "Interactively rename a media file",
		Long:  "Looks up naming options and allows you to select one to rename the file.",
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
				fmt.Printf("[%d] %s (votes: %d)\n", i+1, submission.Filename, submission.VoteScore)
			}

			// Prompt for selection
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("\nSelect an option (1-" + fmt.Sprint(len(response.Submissions)) + ") or 'q' to quit: ")
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

			// Get directory and new path
			dir := filepath.Dir(filePath)
			ext := filepath.Ext(filePath)
			newFilename := selectedSubmission.Filename

			// Ensure the new filename has the same extension
			if !strings.HasSuffix(newFilename, ext) {
				newFilename = strings.TrimSuffix(newFilename, filepath.Ext(newFilename)) + ext
			}

			newPath := filepath.Join(dir, newFilename)

			// Preview rename
			fmt.Printf("\nRename:\n")
			fmt.Printf("  From: %s\n", filepath.Base(filePath))
			fmt.Printf("  To:   %s\n", newFilename)

			if dryRun {
				fmt.Println("\n[DRY RUN] No changes made.")
				return nil
			}

			// Confirm
			fmt.Print("\nConfirm rename? (y/n): ")
			confirm, _ := reader.ReadString('\n')
			confirm = strings.TrimSpace(strings.ToLower(confirm))

			if confirm != "y" && confirm != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}

			// Perform rename
			if err := os.Rename(filePath, newPath); err != nil {
				return fmt.Errorf("failed to rename file: %w", err)
			}

			fmt.Println("\nFile renamed successfully!")
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview rename without making changes")

	return cmd
}
