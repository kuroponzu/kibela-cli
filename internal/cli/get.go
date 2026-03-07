package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kuroponzu/kibela-cli/internal/config"
	"github.com/kuroponzu/kibela-cli/internal/kibela"
	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	var (
		noteID   string
		notePath string
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a note from Kibela",
		Long:  `Get a note from Kibela by ID or path.`,
		Example: `  # Get a note by ID
  kibela get --id abc123

  # Get a note by path
  kibela get --path "/notes/12345"

  # Get a note and output as JSON
  kibela get --id abc123 --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if noteID == "" && notePath == "" {
				return fmt.Errorf("either --id or --path must be specified")
			}

			if noteID != "" && notePath != "" {
				return fmt.Errorf("only one of --id or --path can be specified")
			}

			cfg, err := loadConfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(config.ExitConfig)
			}

			client := kibela.NewClient(cfg)
			ctx := context.Background()

			var note *kibela.Note

			if noteID != "" {
				note, err = client.GetNoteByID(ctx, noteID)
			} else {
				note, err = client.GetNoteByPath(ctx, notePath)
			}

			if err != nil {
				if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "404") {
					fmt.Fprintf(os.Stderr, "Error: Note not found\n")
					os.Exit(config.ExitNotFound)
				}
				if strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "401") {
					fmt.Fprintf(os.Stderr, "Error: Authentication failed\n")
					os.Exit(config.ExitAuth)
				}
				if strings.Contains(err.Error(), "forbidden") || strings.Contains(err.Error(), "403") {
					fmt.Fprintf(os.Stderr, "Error: Permission denied\n")
					os.Exit(config.ExitPermission)
				}
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(config.ExitGraphQL)
			}

			formatter := NewFormatter(os.Stdout, jsonOutput)
			if err := formatter.PrintNote(note); err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to format output: %s\n", err)
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&noteID, "id", "", "Note ID")
	cmd.Flags().StringVar(&notePath, "path", "", "Note path (e.g., /notes/12345)")

	return cmd
}
