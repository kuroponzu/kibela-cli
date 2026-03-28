package kibela

import "time"

// Note represents a Kibela note.
type Note struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	ContentHTML   string     `json:"contentHtml,omitempty"`
	CoEditing     bool       `json:"coEditing"`
	PublishedAt   *time.Time `json:"publishedAt,omitempty"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	URL           string     `json:"url"`
	Path          string     `json:"path"`
	Author        *User      `json:"author,omitempty"`
	Groups        []Group    `json:"groups,omitempty"`
	Folders       []Folder   `json:"folders,omitempty"`
}

// User represents a Kibela user.
type User struct {
	ID       string `json:"id"`
	Account  string `json:"account"`
	RealName string `json:"realName"`
}

// Group represents a Kibela group.
type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Folder represents a Kibela folder.
type Folder struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
}

// CreateNoteParams represents input parameters for creating a note.
type CreateNoteParams struct {
	Title     string
	Content   string
	GroupIDs  []string
	CoEditing bool
	Draft     bool
	FolderID  string
}

// UpdateNoteParams represents input parameters for updating a note.
type UpdateNoteParams struct {
	ID        string
	Title     *string
	Content   *string
	CoEditing *bool
	Draft     *bool
}
