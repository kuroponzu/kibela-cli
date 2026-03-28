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

func newGroupsCmd() *cobra.Command {
	var (
		first       int
		showArchived bool
	)

	cmd := &cobra.Command{
		Use:   "groups",
		Short: "List groups from Kibela",
		Long:  `List all groups from Kibela.`,
		Example: `  # List groups
  kibela groups

  # List groups with limit
  kibela groups --first 50

  # Include archived groups
  kibela groups --archived

  # Output as JSON
  kibela groups --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(config.ExitConfig)
			}

			client := kibela.NewClient(cfg)
			ctx := context.Background()

			groups, err := client.GetGroups(ctx, first)
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

			// Filter out archived groups unless --archived is specified
			if !showArchived {
				filtered := make([]kibela.Group, 0, len(groups))
				for _, g := range groups {
					if !g.IsArchived {
						filtered = append(filtered, g)
					}
				}
				groups = filtered
			}

			formatter := NewFormatter(os.Stdout, jsonOutput)
			if err := formatter.PrintGroups(groups); err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to format output: %s\n", err)
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().IntVar(&first, "first", 100, "Number of groups to retrieve")
	cmd.Flags().BoolVar(&showArchived, "archived", false, "Include archived groups")

	return cmd
}
