package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "mkvmender",
		Short: "MKV Mender - Community-driven media file renaming tool",
		Long: `MKV Mender helps you rename your movie and TV show rips by matching
file hashes with community-submitted naming data.`,
		Version: version,
	}

	// Add commands
	rootCmd.AddCommand(newHashCmd())
	rootCmd.AddCommand(newLookupCmd())
	rootCmd.AddCommand(newRenameCmd())
	rootCmd.AddCommand(newUploadCmd())
	rootCmd.AddCommand(newVoteCmd())
	rootCmd.AddCommand(newBatchCmd())
	rootCmd.AddCommand(newSearchCmd())
	rootCmd.AddCommand(newLoginCmd())
	rootCmd.AddCommand(newRegisterCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
