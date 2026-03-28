package kibela

import "github.com/shurcooL/graphql"

// GetNoteByIDQuery is a GraphQL query to get a note by ID.
type GetNoteByIDQuery struct {
	Note struct {
		ID          graphql.ID
		Title       graphql.String
		Content     graphql.String
		ContentHTML graphql.String `graphql:"contentHtml"`
		CoEditing   graphql.Boolean
		PublishedAt *graphql.String
		UpdatedAt   graphql.String
		URL         graphql.String `graphql:"url"`
		Path        graphql.String
		Author      struct {
			ID       graphql.ID
			Account  graphql.String
			RealName graphql.String
		}
		Groups struct {
			Nodes []struct {
				ID   graphql.ID
				Name graphql.String
			}
		} `graphql:"groups(first: 10)"`
		Folders struct {
			Nodes []struct {
				ID       graphql.ID
				FullName graphql.String
			}
		} `graphql:"folders(first: 10)"`
	} `graphql:"note(id: $id)"`
}

// GetNoteByPathQuery is a GraphQL query to get a note by path.
type GetNoteByPathQuery struct {
	NoteFromPath struct {
		ID          graphql.ID
		Title       graphql.String
		Content     graphql.String
		ContentHTML graphql.String `graphql:"contentHtml"`
		CoEditing   graphql.Boolean
		PublishedAt *graphql.String
		UpdatedAt   graphql.String
		URL         graphql.String `graphql:"url"`
		Path        graphql.String
		Author      struct {
			ID       graphql.ID
			Account  graphql.String
			RealName graphql.String
		}
		Groups struct {
			Nodes []struct {
				ID   graphql.ID
				Name graphql.String
			}
		} `graphql:"groups(first: 10)"`
		Folders struct {
			Nodes []struct {
				ID       graphql.ID
				FullName graphql.String
			}
		} `graphql:"folders(first: 10)"`
	} `graphql:"noteFromPath(path: $path)"`
}

// CreateNoteMutation is a GraphQL mutation to create a note.
type CreateNoteMutation struct {
	CreateNote struct {
		Note struct {
			ID          graphql.ID
			Title       graphql.String
			Content     graphql.String
			CoEditing   graphql.Boolean
			PublishedAt *graphql.String
			UpdatedAt   graphql.String
			URL         graphql.String `graphql:"url"`
			Path        graphql.String
		}
	} `graphql:"createNote(input: $input)"`
}

// CreateNoteInputGQL is the GraphQL input type for creating a note.
type CreateNoteInputGQL struct {
	Title     graphql.String   `json:"title"`
	Content   graphql.String   `json:"content"`
	GroupIds  []graphql.ID     `json:"groupIds"`
	CoEditing graphql.Boolean  `json:"coEditing"`
	Draft     graphql.Boolean  `json:"draft"`
	FolderId  *graphql.ID      `json:"folderId,omitempty"`
}

// UpdateNoteMutation is a GraphQL mutation to update a note.
type UpdateNoteMutation struct {
	UpdateNote struct {
		Note struct {
			ID          graphql.ID
			Title       graphql.String
			Content     graphql.String
			CoEditing   graphql.Boolean
			PublishedAt *graphql.String
			UpdatedAt   graphql.String
			URL         graphql.String `graphql:"url"`
			Path        graphql.String
		}
	} `graphql:"updateNote(input: $input)"`
}

// UpdateNoteInputGQL is the GraphQL input type for updating a note.
type UpdateNoteInputGQL struct {
	ID        graphql.ID       `json:"id"`
	NewNote   *NoteContentInput `json:"newNote,omitempty"`
	BaseNote  *NoteContentInput `json:"baseNote,omitempty"`
	Draft     *graphql.Boolean `json:"draft,omitempty"`
	CoEditing *graphql.Boolean `json:"coEditing,omitempty"`
}

// NoteContentInput is the input type for note content.
type NoteContentInput struct {
	Title   graphql.String `json:"title"`
	Content graphql.String `json:"content"`
}
