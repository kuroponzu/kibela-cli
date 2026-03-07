package cli

import (
	"fmt"
	"os"

	"github.com/kuroponzu/kibela-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	verbose    bool
)

// NewRootCmd creates the root command.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "kibela",
		Short: "A CLI tool for Kibela",
		Long:  `kibela is a CLI tool for interacting with Kibela API. It allows you to get, create, and update notes.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip config check for help command
			if cmd.Name() == "help" {
				return nil
			}
			return nil
		},
	}

	// Global flags
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Add subcommands
	rootCmd.AddCommand(newGetCmd())
	rootCmd.AddCommand(newCreateCmd())
	rootCmd.AddCommand(newUpdateCmd())

	return rootCmd
}

// Execute runs the root command.
func Execute() {
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// loadConfig loads and validates configuration.
func loadConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// exitWithError prints an error message and exits with the specified code.
func exitWithError(msg string, code int) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	os.Exit(code)
}
