package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/quentinsteinke/mkvmender/internal/api"
	"github.com/spf13/cobra"
)

func newRegisterCmd() *cobra.Command {
	var username string
	var baseURL string

	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a new user account",
		Long:  "Create a new user account and receive an API key.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load existing config for base URL
			config, err := api.LoadConfig()
			if err != nil {
				config = api.DefaultConfig()
			}

			if baseURL != "" {
				config.BaseURL = baseURL
			}

			// Prompt for username if not provided
			if username == "" {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter username: ")
				input, _ := reader.ReadString('\n')
				username = strings.TrimSpace(input)
			}

			if username == "" {
				return fmt.Errorf("username is required")
			}

			// Create client without API key
			client := api.New(config.BaseURL, "")

			// Register user
			fmt.Println("Registering user...")
			user, err := client.Register(username)
			if err != nil {
				return fmt.Errorf("registration failed: %w", err)
			}

			fmt.Printf("\nRegistration successful!\n")
			fmt.Printf("Username: %s\n", user.Username)
			fmt.Printf("API Key:  %s\n", user.APIKey)
			fmt.Println("\nIMPORTANT: Save your API key securely. You will need it to authenticate.")

			// Offer to save config
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("\nSave API key to config file? (y/n): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))

			if input == "y" || input == "yes" {
				config.APIKey = user.APIKey
				if err := api.SaveConfig(config); err != nil {
					return fmt.Errorf("failed to save config: %w", err)
				}

				configPath, _ := api.ConfigPath()
				fmt.Printf("Configuration saved to: %s\n", configPath)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "Username for registration")
	cmd.Flags().StringVar(&baseURL, "url", "", "Base URL for the API server")

	return cmd
}
