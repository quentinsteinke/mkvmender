package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/quentinsteinke/mkvmender/internal/api"
	"github.com/quentinsteinke/mkvmender/internal/hasher"
	"github.com/spf13/cobra"
)

func newBatchCmd() *cobra.Command {
	var dryRun bool
	var extensions []string

	cmd := &cobra.Command{
		Use:   "batch <directory>",
		Short: "Process all media files in a directory",
		Long:  "Recursively process all media files in a directory and look up naming options.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			directory := args[0]

			// Check if directory exists
			info, err := os.Stat(directory)
			if err != nil {
				return fmt.Errorf("failed to access directory: %w", err)
			}
			if !info.IsDir() {
				return fmt.Errorf("path is not a directory")
			}

			// Default extensions if not specified
			if len(extensions) == 0 {
				extensions = []string{".mkv", ".mp4", ".avi", ".m4v"}
			}

			// Create API client
			client, err := api.NewClient()
			if err != nil {
				return fmt.Errorf("failed to create API client: %w", err)
			}

			// Find all media files
			var files []string
			err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					ext := strings.ToLower(filepath.Ext(path))
					for _, validExt := range extensions {
						if ext == strings.ToLower(validExt) {
							files = append(files, path)
							break
						}
					}
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("failed to walk directory: %w", err)
			}

			if len(files) == 0 {
				fmt.Println("No media files found.")
				return nil
			}

			fmt.Printf("Found %d media file(s)\n\n", len(files))

			// Process each file
			for i, file := range files {
				fmt.Printf("[%d/%d] Processing: %s\n", i+1, len(files), filepath.Base(file))

				// Hash the file
				result, err := hasher.HashFile(file)
				if err != nil {
					fmt.Printf("  Error hashing file: %v\n\n", err)
					continue
				}

				// Lookup
				response, err := client.Lookup(result.Hash)
				if err != nil {
					fmt.Printf("  Error looking up: %v\n\n", err)
					continue
				}

				if len(response.Submissions) == 0 {
					fmt.Printf("  No naming submissions found\n\n")
					continue
				}

				// Show top result
				top := response.Submissions[0]
				fmt.Printf("  Best match: %s (votes: %d)\n", top.Filename, top.VoteScore)

				if !dryRun {
					fmt.Printf("  (Use 'mkvmender rename' to apply)\n")
				}
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview without making changes")
	cmd.Flags().StringSliceVarP(&extensions, "ext", "e", nil, "File extensions to process (default: .mkv,.mp4,.avi,.m4v)")

	return cmd
}
