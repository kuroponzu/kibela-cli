package kibela

import "github.com/shurcooL/graphql"

// GetNoteByIDQuery is a GraphQL query to get a note by ID.
type GetNoteByIDQuery struct {
	Note struct {
		ID          graphql.ID      `graphql:"id"`
		Title       graphql.String  `graphql:"title"`
		Content     graphql.String  `graphql:"content"`
		ContentHTML graphql.String  `graphql:"contentHtml"`
		Coediting   graphql.Boolean `graphql:"coediting"`
		PublishedAt *graphql.String `graphql:"publishedAt"`
		UpdatedAt   graphql.String  `graphql:"updatedAt"`
		URL         graphql.String  `graphql:"url"`
		Path        graphql.String  `graphql:"path"`
		Author      struct {
			ID       graphql.ID     `graphql:"id"`
			Account  graphql.String `graphql:"account"`
			RealName graphql.String `graphql:"realName"`
		} `graphql:"author"`
		Groups struct {
			Nodes []struct {
				ID   graphql.ID     `graphql:"id"`
				Name graphql.String `graphql:"name"`
			}
		} `graphql:"groups(first: 10)"`
		Folders struct {
			Nodes []struct {
				ID       graphql.ID     `graphql:"id"`
				FullName graphql.String `graphql:"fullName"`
			}
		} `graphql:"folders(first: 10)"`
	} `graphql:"note(id: $id)"`
}

// GetNoteByPathQuery is a GraphQL query to get a note by path.
type GetNoteByPathQuery struct {
	NoteFromPath struct {
		ID          graphql.ID      `graphql:"id"`
		Title       graphql.String  `graphql:"title"`
		Content     graphql.String  `graphql:"content"`
		ContentHTML graphql.String  `graphql:"contentHtml"`
		Coediting   graphql.Boolean `graphql:"coediting"`
		PublishedAt *graphql.String `graphql:"publishedAt"`
		UpdatedAt   graphql.String  `graphql:"updatedAt"`
		URL         graphql.String  `graphql:"url"`
		Path        graphql.String  `graphql:"path"`
		Author      struct {
			ID       graphql.ID     `graphql:"id"`
			Account  graphql.String `graphql:"account"`
			RealName graphql.String `graphql:"realName"`
		} `graphql:"author"`
		Groups struct {
			Nodes []struct {
				ID   graphql.ID     `graphql:"id"`
				Name graphql.String `graphql:"name"`
			}
		} `graphql:"groups(first: 10)"`
		Folders struct {
			Nodes []struct {
				ID       graphql.ID     `graphql:"id"`
				FullName graphql.String `graphql:"fullName"`
			}
		} `graphql:"folders(first: 10)"`
	} `graphql:"noteFromPath(path: $path)"`
}

// CreateNoteMutation is a GraphQL mutation to create a note.
type CreateNoteMutation struct {
	CreateNote struct {
		Note struct {
			ID          graphql.ID      `graphql:"id"`
			Title       graphql.String  `graphql:"title"`
			Content     graphql.String  `graphql:"content"`
			Coediting   graphql.Boolean `graphql:"coediting"`
			PublishedAt *graphql.String `graphql:"publishedAt"`
			UpdatedAt   graphql.String  `graphql:"updatedAt"`
			URL         graphql.String  `graphql:"url"`
			Path        graphql.String  `graphql:"path"`
		} `graphql:"note"`
	} `graphql:"createNote(input: $input)"`
}

// CreateNoteInput is the GraphQL input type for creating a note.
// Note: The type name must match the GraphQL input type name exactly.
type CreateNoteInput struct {
	Title     graphql.String  `json:"title"`
	Content   graphql.String  `json:"content"`
	GroupIds  []graphql.ID    `json:"groupIds"`
	Coediting graphql.Boolean `json:"coediting"`
	Draft     graphql.Boolean `json:"draft"`
	FolderId  *graphql.ID     `json:"folderId,omitempty"`
}

// UpdateNoteMutation is a GraphQL mutation to update a note.
type UpdateNoteMutation struct {
	UpdateNote struct {
		Note struct {
			ID          graphql.ID      `graphql:"id"`
			Title       graphql.String  `graphql:"title"`
			Content     graphql.String  `graphql:"content"`
			Coediting   graphql.Boolean `graphql:"coediting"`
			PublishedAt *graphql.String `graphql:"publishedAt"`
			UpdatedAt   graphql.String  `graphql:"updatedAt"`
			URL         graphql.String  `graphql:"url"`
			Path        graphql.String  `graphql:"path"`
		} `graphql:"note"`
	} `graphql:"updateNote(input: $input)"`
}

// UpdateNoteInput is the GraphQL input type for updating a note.
// Note: The type name must match the GraphQL input type name exactly.
type UpdateNoteInput struct {
	ID        graphql.ID        `json:"id"`
	NewNote   *NoteContentInput `json:"newNote,omitempty"`
	BaseNote  *NoteContentInput `json:"baseNote,omitempty"`
	Draft     *graphql.Boolean  `json:"draft,omitempty"`
	Coediting *graphql.Boolean  `json:"coediting,omitempty"`
}

// NoteContentInput is the input type for note content.
type NoteContentInput struct {
	Title   graphql.String `json:"title"`
	Content graphql.String `json:"content"`
}

// GetGroupsQuery is a GraphQL query to get groups.
type GetGroupsQuery struct {
	Groups struct {
		Edges []struct {
			Node struct {
				ID         graphql.ID      `graphql:"id"`
				Name       graphql.String  `graphql:"name"`
				IsDefault  graphql.Boolean `graphql:"isDefault"`
				IsArchived graphql.Boolean `graphql:"isArchived"`
			}
		}
	} `graphql:"groups(first: $first)"`
}
