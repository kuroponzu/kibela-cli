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

func newCreateCmd() *cobra.Command {
	var (
		title       string
		content     string
		contentFile string
		groupID     string
		folderID    string
		coEditing   bool
		draft       bool
		stdin       bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new note in Kibela",
		Long:  `Create a new note in Kibela with the specified title and content.`,
		Example: `  # Create a note with inline content
  kibela create --title "My Note" --content "Hello, world!" --group-id abc123

  # Create a note from a file
  kibela create --title "My Note" --content-file ./article.md --group-id abc123

  # Create a note from stdin
  cat article.md | kibela create --title "My Note" --group-id abc123 --stdin

  # Create a draft note
  kibela create --title "Draft Note" --content "WIP" --group-id abc123 --draft

  # Create a co-editing note
  kibela create --title "Shared Note" --content "Let's edit together" --group-id abc123 --co-editing`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate required fields
			if title == "" {
				return fmt.Errorf("--title is required")
			}

			if groupID == "" {
				return fmt.Errorf("--group-id is required")
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

			if contentSources == 0 {
				return fmt.Errorf("one of --content, --content-file, or --stdin must be specified")
			}

			if contentSources > 1 {
				return fmt.Errorf("only one of --content, --content-file, or --stdin can be specified")
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

			input := &kibela.CreateNoteInput{
				Title:     title,
				Content:   content,
				GroupIDs:  []string{groupID},
				CoEditing: coEditing,
				Draft:     draft,
				FolderID:  folderID,
			}

			note, err := client.CreateNote(ctx, input)
			if err != nil {
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
			if err := formatter.PrintNoteCreated(note); err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to format output: %s\n", err)
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Note title (required)")
	cmd.Flags().StringVar(&content, "content", "", "Note content")
	cmd.Flags().StringVar(&contentFile, "content-file", "", "Path to file containing note content")
	cmd.Flags().StringVar(&groupID, "group-id", "", "Group ID to create the note in (required)")
	cmd.Flags().StringVar(&folderID, "folder-id", "", "Folder ID to place the note in")
	cmd.Flags().BoolVar(&coEditing, "co-editing", false, "Enable co-editing")
	cmd.Flags().BoolVar(&draft, "draft", false, "Create as draft")
	cmd.Flags().BoolVar(&stdin, "stdin", false, "Read content from stdin")

	return cmd
}
