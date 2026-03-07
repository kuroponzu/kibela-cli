package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kuroponzu/kibela-cli/internal/config"
	"github.com/kuroponzu/kibela-cli/internal/kibela"
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	var (
		noteID      string
		title       string
		content     string
		contentFile string
		coEditing   *bool
		draft       *bool
		stdin       bool
	)

	// Use separate bool vars for flags since cobra doesn't directly support *bool
	var coEditingFlag bool
	var coEditingSet bool
	var draftFlag bool
	var draftSet bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing note in Kibela",
		Long:  `Update an existing note in Kibela. Only specified fields will be updated.`,
		Example: `  # Update note title
  kibela update --id abc123 --title "New Title"

  # Update note content
  kibela update --id abc123 --content "New content"

  # Update note content from file
  kibela update --id abc123 --content-file ./updated.md

  # Update note content from stdin
  cat updated.md | kibela update --id abc123 --stdin

  # Update both title and content
  kibela update --id abc123 --title "New Title" --content "New content"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate required fields
			if noteID == "" {
				return fmt.Errorf("--id is required")
			}

			// Check if coEditing flag was set
			if cmd.Flags().Changed("co-editing") {
				coEditingSet = true
			}
			if cmd.Flags().Changed("draft") {
				draftSet = true
			}

			// Determine content source
			contentSources := 0
			if content != "" {
				contentSources++
			}
			if contentFile != "" {
				contentSources++
			}
			if stdin {
				contentSources++
			}

			if contentSources > 1 {
				return fmt.Errorf("only one of --content, --content-file, or --stdin can be specified")
			}

			// Check if at least one update field is specified
			hasUpdate := title != "" || contentSources > 0 || coEditingSet || draftSet
			if !hasUpdate {
				return fmt.Errorf("at least one of --title, --content, --content-file, --stdin, --co-editing, or --draft must be specified")
			}

			// Get content from file if specified
			if contentFile != "" {
				data, err := os.ReadFile(contentFile)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: Failed to read file: %s\n", err)
					os.Exit(config.ExitInput)
				}
				content = string(data)
			}

			// Get content from stdin if specified
			if stdin {
				reader := bufio.NewReader(os.Stdin)
				var sb strings.Builder
				for {
					line, err := reader.ReadString('\n')
					if err != nil {
						if err == io.EOF {
							sb.WriteString(line)
							break
						}
						fmt.Fprintf(os.Stderr, "Error: Failed to read stdin: %s\n", err)
						os.Exit(config.ExitInput)
					}
					sb.WriteString(line)
				}
				content = sb.String()
			}

			cfg, err := loadConfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(config.ExitConfig)
			}

			client := kibela.NewClient(cfg)
			ctx := context.Background()

			input := &kibela.UpdateNoteInput{
				ID: noteID,
			}

			if title != "" {
				input.Title = &title
			}

			if content != "" {
				input.Content = &content
			}

			if coEditingSet {
				coEditing = &coEditingFlag
				input.CoEditing = coEditing
			}

			if draftSet {
				draft = &draftFlag
				input.Draft = draft
			}

			note, err := client.UpdateNote(ctx, input)
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
			if err := formatter.PrintNoteUpdated(note); err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to format output: %s\n", err)
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&noteID, "id", "", "Note ID (required)")
	cmd.Flags().StringVar(&title, "title", "", "New note title")
	cmd.Flags().StringVar(&content, "content", "", "New note content")
	cmd.Flags().StringVar(&contentFile, "content-file", "", "Path to file containing new note content")
	cmd.Flags().BoolVar(&coEditingFlag, "co-editing", false, "Enable/disable co-editing")
	cmd.Flags().BoolVar(&draftFlag, "draft", false, "Set draft status")
	cmd.Flags().BoolVar(&stdin, "stdin", false, "Read content from stdin")

	return cmd
}
