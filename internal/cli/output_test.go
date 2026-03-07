package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/kuroponzu/kibela-cli/internal/kibela"
)

func TestFormatter_PrintNote_JSON(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, true)

	updatedAt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	note := &kibela.Note{
		ID:        "note-123",
		Title:     "Test Note",
		Content:   "Hello, world!",
		CoEditing: false,
		UpdatedAt: updatedAt,
		URL:       "https://test.kibe.la/notes/123",
		Path:      "/notes/123",
		Author: &kibela.User{
			ID:       "user-1",
			Account:  "testuser",
			RealName: "Test User",
		},
		Groups: []kibela.Group{
			{ID: "group-1", Name: "General"},
		},
	}

	err := formatter.PrintNote(note)
	if err != nil {
		t.Fatalf("PrintNote() error = %v", err)
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Check key fields
	if result["id"] != "note-123" {
		t.Errorf("JSON id = %v, want %v", result["id"], "note-123")
	}
	if result["title"] != "Test Note" {
		t.Errorf("JSON title = %v, want %v", result["title"], "Test Note")
	}
}

func TestFormatter_PrintNote_Human(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, false)

	updatedAt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	note := &kibela.Note{
		ID:        "note-123",
		Title:     "Test Note",
		Content:   "Hello, world!",
		CoEditing: true,
		UpdatedAt: updatedAt,
		URL:       "https://test.kibe.la/notes/123",
		Path:      "/notes/123",
		Author: &kibela.User{
			ID:       "user-1",
			Account:  "testuser",
			RealName: "Test User",
		},
		Groups: []kibela.Group{
			{ID: "group-1", Name: "General"},
		},
	}

	err := formatter.PrintNote(note)
	if err != nil {
		t.Fatalf("PrintNote() error = %v", err)
	}

	output := buf.String()

	// Check that key information is present
	checks := []string{
		"Title: Test Note",
		"ID: note-123",
		"URL: https://test.kibe.la/notes/123",
		"Author: Test User (@testuser)",
		"Groups: General",
		"CoEditing: true",
		"--- Content ---",
		"Hello, world!",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("Output missing expected text: %q", check)
		}
	}
}

func TestFormatter_PrintNoteCreated_Human(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, false)

	updatedAt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	note := &kibela.Note{
		ID:        "note-123",
		Title:     "New Note",
		UpdatedAt: updatedAt,
		URL:       "https://test.kibe.la/notes/123",
		Path:      "/notes/123",
	}

	err := formatter.PrintNoteCreated(note)
	if err != nil {
		t.Fatalf("PrintNoteCreated() error = %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "Note created successfully!") {
		t.Error("Output missing success message")
	}
	if !strings.Contains(output, "Title: New Note") {
		t.Error("Output missing title")
	}
}

func TestFormatter_PrintNoteUpdated_Human(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, false)

	updatedAt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	note := &kibela.Note{
		ID:        "note-123",
		Title:     "Updated Note",
		UpdatedAt: updatedAt,
		URL:       "https://test.kibe.la/notes/123",
	}

	err := formatter.PrintNoteUpdated(note)
	if err != nil {
		t.Fatalf("PrintNoteUpdated() error = %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "Note updated successfully!") {
		t.Error("Output missing success message")
	}
	if !strings.Contains(output, "Title: Updated Note") {
		t.Error("Output missing title")
	}
}

func TestNewFormatter(t *testing.T) {
	var buf bytes.Buffer

	// Test JSON format
	jsonFormatter := NewFormatter(&buf, true)
	if jsonFormatter.Format != FormatJSON {
		t.Errorf("NewFormatter(true) Format = %v, want FormatJSON", jsonFormatter.Format)
	}

	// Test Human format
	humanFormatter := NewFormatter(&buf, false)
	if humanFormatter.Format != FormatHuman {
		t.Errorf("NewFormatter(false) Format = %v, want FormatHuman", humanFormatter.Format)
	}
}
