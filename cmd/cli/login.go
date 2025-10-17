package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/quentinsteinke/mkvmender/internal/api"
	"github.com/spf13/cobra"
)

func newLoginCmd() *cobra.Command {
	var apiKey string
	var baseURL string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Configure API credentials",
		Long:  "Set your API key and server URL for authentication.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load existing config
			config, err := api.LoadConfig()
			if err != nil {
				config = api.DefaultConfig()
			}

			// Prompt for API key if not provided
			if apiKey == "" {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter your API key: ")
				input, _ := reader.ReadString('\n')
				apiKey = strings.TrimSpace(input)
			}

			if apiKey == "" {
				return fmt.Errorf("API key is required")
			}

			// Update config
			config.APIKey = apiKey
			if baseURL != "" {
				config.BaseURL = baseURL
			}

			// Save config
			if err := api.SaveConfig(config); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			// Test the connection
			client := api.New(config.BaseURL, config.APIKey)
			if err := client.Health(); err != nil {
				fmt.Printf("Warning: Could not connect to server at %s\n", config.BaseURL)
			} else {
				fmt.Println("Successfully connected to server!")
			}

			configPath, _ := api.ConfigPath()
			fmt.Printf("\nConfiguration saved to: %s\n", configPath)

			return nil
		},
	}

	cmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "API key for authentication")
	cmd.Flags().StringVarP(&baseURL, "url", "u", "", "Base URL for the API server")

	return cmd
}
