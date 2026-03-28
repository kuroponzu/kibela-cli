package kibela

import (
	"encoding/json"
	"testing"
)

func TestGroup_JSONMarshal(t *testing.T) {
	group := &Group{
		ID:         "group-123",
		Name:       "Test Group",
		IsDefault:  true,
		IsArchived: false,
	}

	data, err := json.Marshal(group)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if result["id"] != "group-123" {
		t.Errorf("id = %v, want %v", result["id"], "group-123")
	}
	if result["name"] != "Test Group" {
		t.Errorf("name = %v, want %v", result["name"], "Test Group")
	}
	if result["isDefault"] != true {
		t.Errorf("isDefault = %v, want %v", result["isDefault"], true)
	}
	// isArchived is omitted when false due to omitempty
	if _, exists := result["isArchived"]; exists {
		t.Errorf("isArchived should be omitted when false, got %v", result["isArchived"])
	}
}

func TestGroup_JSONMarshal_Archived(t *testing.T) {
	group := &Group{
		ID:         "group-456",
		Name:       "Archived Group",
		IsDefault:  false,
		IsArchived: true,
	}

	data, err := json.Marshal(group)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if result["isArchived"] != true {
		t.Errorf("isArchived = %v, want %v", result["isArchived"], true)
	}
}

func TestGetGroupsQuery_Structure(t *testing.T) {
	// Verify the query structure can be instantiated
	var query GetGroupsQuery

	// Access the nested structure to ensure it compiles correctly
	_ = query.Groups.Edges
}
