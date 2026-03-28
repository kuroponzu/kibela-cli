package kibela

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNote_JSONMarshal(t *testing.T) {
	publishedAt := time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	note := &Note{
		ID:          "note-123",
		Title:       "Test Note",
		Content:     "Hello, world!",
		ContentHTML: "<p>Hello, world!</p>",
		CoEditing:   true,
		PublishedAt: &publishedAt,
		UpdatedAt:   updatedAt,
		URL:         "https://test.kibe.la/notes/123",
		Path:        "/notes/123",
		Author: &User{
			ID:       "user-1",
			Account:  "testuser",
			RealName: "Test User",
		},
		Groups: []Group{
			{ID: "group-1", Name: "General"},
		},
		Folders: []Folder{
			{ID: "folder-1", FullName: "Documents/Notes"},
		},
	}

	data, err := json.Marshal(note)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	// Check key fields
	if result["id"] != "note-123" {
		t.Errorf("id = %v, want %v", result["id"], "note-123")
	}
	if result["title"] != "Test Note" {
		t.Errorf("title = %v, want %v", result["title"], "Test Note")
	}
	if result["coEditing"] != true {
		t.Errorf("coEditing = %v, want %v", result["coEditing"], true)
	}
}

func TestCreateNoteInput(t *testing.T) {
	input := &CreateNoteInput{
		Title:     "New Note",
		Content:   "Content here",
		GroupIDs:  []string{"group-1", "group-2"},
		CoEditing: false,
		Draft:     true,
		FolderID:  "folder-1",
	}

	if input.Title != "New Note" {
		t.Errorf("Title = %v, want %v", input.Title, "New Note")
	}
	if len(input.GroupIDs) != 2 {
		t.Errorf("GroupIDs length = %v, want %v", len(input.GroupIDs), 2)
	}
	if input.Draft != true {
		t.Errorf("Draft = %v, want %v", input.Draft, true)
	}
}

func TestUpdateNoteInput(t *testing.T) {
	title := "Updated Title"
	content := "Updated Content"
	coEditing := true

	input := &UpdateNoteInput{
		ID:        "note-123",
		Title:     &title,
		Content:   &content,
		CoEditing: &coEditing,
	}

	if input.ID != "note-123" {
		t.Errorf("ID = %v, want %v", input.ID, "note-123")
	}
	if *input.Title != "Updated Title" {
		t.Errorf("Title = %v, want %v", *input.Title, "Updated Title")
	}
	if *input.CoEditing != true {
		t.Errorf("CoEditing = %v, want %v", *input.CoEditing, true)
	}
}
