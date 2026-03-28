package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/kuroponzu/kibela-cli/internal/kibela"
)

// OutputFormat represents the output format type.
type OutputFormat int

const (
	FormatHuman OutputFormat = iota
	FormatJSON
)

// Formatter handles output formatting.
type Formatter struct {
	Format OutputFormat
	Writer io.Writer
}

// NewFormatter creates a new formatter.
func NewFormatter(w io.Writer, jsonOutput bool) *Formatter {
	format := FormatHuman
	if jsonOutput {
		format = FormatJSON
	}
	return &Formatter{
		Format: format,
		Writer: w,
	}
}

// PrintNote prints a note in the configured format.
func (f *Formatter) PrintNote(note *kibela.Note) error {
	if f.Format == FormatJSON {
		return f.printJSON(note)
	}
	return f.printNoteHuman(note)
}

// PrintNoteCreated prints a created note in the configured format.
func (f *Formatter) PrintNoteCreated(note *kibela.Note) error {
	if f.Format == FormatJSON {
		return f.printJSON(note)
	}
	return f.printNoteCreatedHuman(note)
}

// PrintNoteUpdated prints an updated note in the configured format.
func (f *Formatter) PrintNoteUpdated(note *kibela.Note) error {
	if f.Format == FormatJSON {
		return f.printJSON(note)
	}
	return f.printNoteUpdatedHuman(note)
}

func (f *Formatter) printJSON(v interface{}) error {
	encoder := json.NewEncoder(f.Writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

func (f *Formatter) printNoteHuman(note *kibela.Note) error {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Title: %s\n", note.Title))
	sb.WriteString(fmt.Sprintf("ID: %s\n", note.ID))
	sb.WriteString(fmt.Sprintf("URL: %s\n", note.URL))
	sb.WriteString(fmt.Sprintf("Path: %s\n", note.Path))

	if note.Author != nil {
		sb.WriteString(fmt.Sprintf("Author: %s (@%s)\n", note.Author.RealName, note.Author.Account))
	}

	if len(note.Groups) > 0 {
		groupNames := make([]string, len(note.Groups))
		for i, g := range note.Groups {
			groupNames[i] = g.Name
		}
		sb.WriteString(fmt.Sprintf("Groups: %s\n", strings.Join(groupNames, ", ")))
	}

	if len(note.Folders) > 0 {
		folderNames := make([]string, len(note.Folders))
		for i, f := range note.Folders {
			folderNames[i] = f.FullName
		}
		sb.WriteString(fmt.Sprintf("Folders: %s\n", strings.Join(folderNames, ", ")))
	}

	sb.WriteString(fmt.Sprintf("CoEditing: %t\n", note.CoEditing))
	sb.WriteString(fmt.Sprintf("UpdatedAt: %s\n", note.UpdatedAt.Format("2006-01-02 15:04:05")))

	if note.PublishedAt != nil {
		sb.WriteString(fmt.Sprintf("PublishedAt: %s\n", note.PublishedAt.Format("2006-01-02 15:04:05")))
	}

	sb.WriteString("\n--- Content ---\n")
	sb.WriteString(note.Content)
	sb.WriteString("\n")

	_, err := f.Writer.Write([]byte(sb.String()))
	return err
}

func (f *Formatter) printNoteCreatedHuman(note *kibela.Note) error {
	var sb strings.Builder

	sb.WriteString("Note created successfully!\n\n")
	sb.WriteString(fmt.Sprintf("Title: %s\n", note.Title))
	sb.WriteString(fmt.Sprintf("ID: %s\n", note.ID))
	sb.WriteString(fmt.Sprintf("URL: %s\n", note.URL))
	sb.WriteString(fmt.Sprintf("Path: %s\n", note.Path))

	_, err := f.Writer.Write([]byte(sb.String()))
	return err
}

func (f *Formatter) printNoteUpdatedHuman(note *kibela.Note) error {
	var sb strings.Builder

	sb.WriteString("Note updated successfully!\n\n")
	sb.WriteString(fmt.Sprintf("Title: %s\n", note.Title))
	sb.WriteString(fmt.Sprintf("ID: %s\n", note.ID))
	sb.WriteString(fmt.Sprintf("URL: %s\n", note.URL))
	sb.WriteString(fmt.Sprintf("UpdatedAt: %s\n", note.UpdatedAt.Format("2006-01-02 15:04:05")))

	_, err := f.Writer.Write([]byte(sb.String()))
	return err
}

// PrintGroups prints groups in the configured format.
func (f *Formatter) PrintGroups(groups []kibela.Group) error {
	if f.Format == FormatJSON {
		return f.printJSON(groups)
	}
	return f.printGroupsHuman(groups)
}

func (f *Formatter) printGroupsHuman(groups []kibela.Group) error {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Found %d groups:\n\n", len(groups)))

	for _, g := range groups {
		defaultStr := ""
		if g.IsDefault {
			defaultStr = " (default)"
		}
		archivedStr := ""
		if g.IsArchived {
			archivedStr = " [archived]"
		}
		sb.WriteString(fmt.Sprintf("- %s (ID: %s)%s%s\n", g.Name, g.ID, defaultStr, archivedStr))
	}

	_, err := f.Writer.Write([]byte(sb.String()))
	return err
}
