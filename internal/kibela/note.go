package kibela

import (
	"context"
	"fmt"
	"time"

	"github.com/shurcooL/graphql"
)

// idToString converts a graphql.ID to a string.
func idToString(id graphql.ID) string {
	if s, ok := id.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", id)
}

// GetNoteByID retrieves a note by its ID.
func (c *Client) GetNoteByID(ctx context.Context, id string) (*Note, error) {
	var query GetNoteByIDQuery
	variables := map[string]interface{}{
		"id": graphql.ID(id),
	}

	if err := c.Query(ctx, &query, variables); err != nil {
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	return convertQueryNoteToNote(
		query.Note.ID,
		query.Note.Title,
		query.Note.Content,
		query.Note.ContentHTML,
		query.Note.CoEditing,
		query.Note.PublishedAt,
		query.Note.UpdatedAt,
		query.Note.URL,
		query.Note.Path,
		query.Note.Author.ID,
		query.Note.Author.Account,
		query.Note.Author.RealName,
		query.Note.Groups.Nodes,
		query.Note.Folders.Nodes,
	)
}

// GetNoteByPath retrieves a note by its path.
func (c *Client) GetNoteByPath(ctx context.Context, path string) (*Note, error) {
	var query GetNoteByPathQuery
	variables := map[string]interface{}{
		"path": graphql.String(path),
	}

	if err := c.Query(ctx, &query, variables); err != nil {
		return nil, fmt.Errorf("failed to get note by path: %w", err)
	}

	return convertQueryNoteToNote(
		query.NoteFromPath.ID,
		query.NoteFromPath.Title,
		query.NoteFromPath.Content,
		query.NoteFromPath.ContentHTML,
		query.NoteFromPath.CoEditing,
		query.NoteFromPath.PublishedAt,
		query.NoteFromPath.UpdatedAt,
		query.NoteFromPath.URL,
		query.NoteFromPath.Path,
		query.NoteFromPath.Author.ID,
		query.NoteFromPath.Author.Account,
		query.NoteFromPath.Author.RealName,
		query.NoteFromPath.Groups.Nodes,
		query.NoteFromPath.Folders.Nodes,
	)
}

// CreateNote creates a new note.
func (c *Client) CreateNote(ctx context.Context, input *CreateNoteInput) (*Note, error) {
	var mutation CreateNoteMutation

	groupIDs := make([]graphql.ID, len(input.GroupIDs))
	for i, id := range input.GroupIDs {
		groupIDs[i] = graphql.ID(id)
	}

	gqlInput := CreateNoteInputGQL{
		Title:     graphql.String(input.Title),
		Content:   graphql.String(input.Content),
		GroupIds:  groupIDs,
		CoEditing: graphql.Boolean(input.CoEditing),
		Draft:     graphql.Boolean(input.Draft),
	}

	if input.FolderID != "" {
		folderID := graphql.ID(input.FolderID)
		gqlInput.FolderId = &folderID
	}

	variables := map[string]interface{}{
		"input": gqlInput,
	}

	if err := c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	updatedAt, _ := time.Parse(time.RFC3339, string(mutation.CreateNote.Note.UpdatedAt))
	var publishedAt *time.Time
	if mutation.CreateNote.Note.PublishedAt != nil {
		t, _ := time.Parse(time.RFC3339, string(*mutation.CreateNote.Note.PublishedAt))
		publishedAt = &t
	}

	return &Note{
		ID:          idToString(mutation.CreateNote.Note.ID),
		Title:       string(mutation.CreateNote.Note.Title),
		Content:     string(mutation.CreateNote.Note.Content),
		CoEditing:   bool(mutation.CreateNote.Note.CoEditing),
		PublishedAt: publishedAt,
		UpdatedAt:   updatedAt,
		URL:         string(mutation.CreateNote.Note.URL),
		Path:        string(mutation.CreateNote.Note.Path),
	}, nil
}

// UpdateNote updates an existing note.
func (c *Client) UpdateNote(ctx context.Context, input *UpdateNoteInput) (*Note, error) {
	// First, get the current note to use as baseNote
	currentNote, err := c.GetNoteByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current note for update: %w", err)
	}

	var mutation UpdateNoteMutation

	// Prepare the new note content
	newTitle := currentNote.Title
	newContent := currentNote.Content
	if input.Title != nil {
		newTitle = *input.Title
	}
	if input.Content != nil {
		newContent = *input.Content
	}

	gqlInput := UpdateNoteInputGQL{
		ID: graphql.ID(input.ID),
		BaseNote: &NoteContentInput{
			Title:   graphql.String(currentNote.Title),
			Content: graphql.String(currentNote.Content),
		},
		NewNote: &NoteContentInput{
			Title:   graphql.String(newTitle),
			Content: graphql.String(newContent),
		},
	}

	if input.CoEditing != nil {
		coEditing := graphql.Boolean(*input.CoEditing)
		gqlInput.CoEditing = &coEditing
	}

	if input.Draft != nil {
		draft := graphql.Boolean(*input.Draft)
		gqlInput.Draft = &draft
	}

	variables := map[string]interface{}{
		"input": gqlInput,
	}

	if err := c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, fmt.Errorf("failed to update note: %w", err)
	}

	updatedAt, _ := time.Parse(time.RFC3339, string(mutation.UpdateNote.Note.UpdatedAt))
	var publishedAt *time.Time
	if mutation.UpdateNote.Note.PublishedAt != nil {
		t, _ := time.Parse(time.RFC3339, string(*mutation.UpdateNote.Note.PublishedAt))
		publishedAt = &t
	}

	return &Note{
		ID:          idToString(mutation.UpdateNote.Note.ID),
		Title:       string(mutation.UpdateNote.Note.Title),
		Content:     string(mutation.UpdateNote.Note.Content),
		CoEditing:   bool(mutation.UpdateNote.Note.CoEditing),
		PublishedAt: publishedAt,
		UpdatedAt:   updatedAt,
		URL:         string(mutation.UpdateNote.Note.URL),
		Path:        string(mutation.UpdateNote.Note.Path),
	}, nil
}

// Helper function to convert GraphQL query result to Note.
func convertQueryNoteToNote(
	id graphql.ID,
	title graphql.String,
	content graphql.String,
	contentHTML graphql.String,
	coEditing graphql.Boolean,
	publishedAtStr *graphql.String,
	updatedAtStr graphql.String,
	url graphql.String,
	path graphql.String,
	authorID graphql.ID,
	authorAccount graphql.String,
	authorRealName graphql.String,
	groupNodes []struct {
		ID   graphql.ID
		Name graphql.String
	},
	folderNodes []struct {
		ID       graphql.ID
		FullName graphql.String
	},
) (*Note, error) {
	updatedAt, _ := time.Parse(time.RFC3339, string(updatedAtStr))
	var publishedAt *time.Time
	if publishedAtStr != nil {
		t, _ := time.Parse(time.RFC3339, string(*publishedAtStr))
		publishedAt = &t
	}

	groups := make([]Group, len(groupNodes))
	for i, g := range groupNodes {
		groups[i] = Group{
			ID:   idToString(g.ID),
			Name: string(g.Name),
		}
	}

	folders := make([]Folder, len(folderNodes))
	for i, f := range folderNodes {
		folders[i] = Folder{
			ID:       idToString(f.ID),
			FullName: string(f.FullName),
		}
	}

	return &Note{
		ID:          idToString(id),
		Title:       string(title),
		Content:     string(content),
		ContentHTML: string(contentHTML),
		CoEditing:   bool(coEditing),
		PublishedAt: publishedAt,
		UpdatedAt:   updatedAt,
		URL:         string(url),
		Path:        string(path),
		Author: &User{
			ID:       idToString(authorID),
			Account:  string(authorAccount),
			RealName: string(authorRealName),
		},
		Groups:  groups,
		Folders: folders,
	}, nil
}
