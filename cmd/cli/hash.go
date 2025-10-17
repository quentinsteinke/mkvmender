package main

import (
	"fmt"

	"github.com/quentinsteinke/mkvmender/internal/hasher"
	"github.com/spf13/cobra"
)

func newHashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hash <file>",
		Short: "Compute hash of a media file",
		Long:  "Computes the SHA-256 hash of a media file and displays file information.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]

			result, err := hasher.HashFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to hash file: %w", err)
			}

			fmt.Printf("File: %s\n", filePath)
			fmt.Printf("Hash: %s\n", result.Hash)
			fmt.Printf("Size: %s (%d bytes)\n", hasher.FormatFileSize(result.FileSize), result.FileSize)

			return nil
		},
	}

	return cmd
}
