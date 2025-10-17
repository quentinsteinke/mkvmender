package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/quentinsteinke/mkvmender/internal/api"
	"github.com/quentinsteinke/mkvmender/internal/hasher"
	"github.com/quentinsteinke/mkvmender/internal/models"
	"github.com/spf13/cobra"
)

func newUploadCmd() *cobra.Command {
	var (
		filename  string
		mediaType string
		title     string
		year      int
		season    int
		episode   int
		quality   string
		source    string
	)

	cmd := &cobra.Command{
		Use:   "upload <file>",
		Short: "Upload naming submission for a file",
		Long:  "Hash a file and upload your naming to help the community.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]

			// Hash the file
			fmt.Println("Hashing file...")
			result, err := hasher.HashFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to hash file: %w", err)
			}

			// Use provided filename or default to current filename
			if filename == "" {
				filename = filepath.Base(filePath)
			}

			// Validate media type
			var mt models.MediaType
			switch strings.ToLower(mediaType) {
			case "movie":
				mt = models.MediaTypeMovie
			case "tv":
				mt = models.MediaTypeTV
			default:
				return fmt.Errorf("invalid media type: must be 'movie' or 'tv'")
			}

			// Create API client
			client, err := api.NewClient()
			if err != nil {
				return fmt.Errorf("failed to create API client: %w", err)
			}

			// Build upload request
			uploadReq := &models.UploadRequest{
				Hash:      result.Hash,
				FileSize:  result.FileSize,
				MediaType: mt,
				Filename:  filename,
			}

			// Add metadata if provided
			if title != "" || year > 0 || season > 0 || episode > 0 || quality != "" || source != "" {
				metadata := &models.NamingMetadata{}

				if title != "" {
					metadata.Title = &title
				}
				if year > 0 {
					metadata.Year = &year
				}
				if season > 0 {
					metadata.Season = &season
				}
				if episode > 0 {
					metadata.Episode = &episode
				}
				if quality != "" {
					metadata.Quality = &quality
				}
				if source != "" {
					metadata.Source = &source
				}

				uploadReq.Metadata = metadata
			}

			// Upload
			fmt.Println("Uploading naming submission...")
			submission, err := client.Upload(uploadReq)
			if err != nil {
				return fmt.Errorf("upload failed: %w", err)
			}

			fmt.Printf("\nSubmission uploaded successfully!\n")
			fmt.Printf("Filename: %s\n", submission.Filename)
			fmt.Printf("Hash: %s\n", result.Hash)
			fmt.Printf("Submission ID: %d\n", submission.ID)

			return nil
		},
	}

	cmd.Flags().StringVarP(&filename, "name", "n", "", "Name for the file (default: current filename)")
	cmd.Flags().StringVarP(&mediaType, "type", "t", "", "Media type: 'movie' or 'tv' (required)")
	cmd.Flags().StringVar(&title, "title", "", "Title of the movie/show")
	cmd.Flags().IntVar(&year, "year", 0, "Release year")
	cmd.Flags().IntVar(&season, "season", 0, "Season number (for TV shows)")
	cmd.Flags().IntVar(&episode, "episode", 0, "Episode number (for TV shows)")
	cmd.Flags().StringVar(&quality, "quality", "", "Quality (e.g., 1080p, 4K)")
	cmd.Flags().StringVar(&source, "source", "", "Source (e.g., Blu-ray, DVD)")

	cmd.MarkFlagRequired("type")

	return cmd
}
