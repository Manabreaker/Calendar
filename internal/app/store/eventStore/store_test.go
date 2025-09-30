package eventStore

import (
	"encoding/json"
	"testing"
)

// helper формирует валидный JSON события с обязательными полями.
// Можно переопределить значения через extra.
func makeValidEventJSON(id string, extra map[string]any) []byte {
	m := map[string]any{
		"id":          id,
		"title":       "Title " + id,
		"date":        "2025-01-01",
		"owner_id":    "owner-" + id,
		"description": "desc",
	}
	for k, v := range extra {
		m[k] = v
	}
	b, _ := json.Marshal(m)
	return b
}

func TestCreateReadUpdateDelete(t *testing.T) {
	store := NewStore()

	// Create success
	if err := store.Create(makeValidEventJSON("1", nil)); err != nil {
		t.Fatalf("Create (ok) unexpected error: %v", err)
	}

	// Create duplicate
	if err := store.Create(makeValidEventJSON("1", nil)); err == nil {
		t.Fatalf("Create (duplicate) expected error, got nil")
	}

	// Read success
	data, err := store.Read("1")
	if err != nil {
		t.Fatalf("Read (ok) unexpected error: %v", err)
	}
	var obj map[string]any
	if json.Unmarshal(data, &obj) != nil || obj["id"] != "1" {
		t.Fatalf("Read returned wrong object: %s", string(data))
	}

	// Read missing
	if _, err := store.Read("missing"); err == nil {
		t.Fatalf("Read (missing) expected error, got nil")
	}

	// Update success (полный валидный объект)
	if err := store.Update(makeValidEventJSON("1", map[string]any{"title": "Updated Title"})); err != nil {
		t.Fatalf("Update (ok) unexpected error: %v", err)
	}
	data, err = store.Read("1")
	if err != nil {
		t.Fatalf("Read after update error: %v", err)
	}
	_ = json.Unmarshal(data, &obj)
	if obj["id"] != "1" || obj["title"] != "Updated Title" {
		t.Fatalf("Updated object mismatch: %+v", obj)
	}

	// Update missing
	if err := store.Update(makeValidEventJSON("nope", nil)); err == nil {
		t.Fatalf("Update (missing) expected error, got nil")
	}

	// Delete success
	if err := store.Delete("1"); err != nil {
		t.Fatalf("Delete (ok) unexpected error: %v", err)
	}
	if _, err := store.Read("1"); err == nil {
		t.Fatalf("Read after delete expected error, got nil")
	}

	// Delete missing
	if err := store.Delete("absent"); err == nil {
		t.Fatalf("Delete (missing) expected error, got nil")
	}
}

func TestCreateValidationFailures(t *testing.T) {
	store := NewStore()

	tests := []struct {
		name string
		buf  []byte
	}{
		{
			name: "missing title",
			buf:  makeValidEventJSON("v1", map[string]any{"title": ""}),
		},
		{
			name: "missing date",
			buf:  makeValidEventJSON("v2", map[string]any{"date": ""}),
		},
		{
			name: "missing owner_id",
			buf:  makeValidEventJSON("v3", map[string]any{"owner_id": ""}),
		},
	}

	for _, tc := range tests {
		if err := store.Create(tc.buf); err == nil {
			t.Fatalf("Create (%s) expected validation error, got nil", tc.name)
		}
	}
}

func TestCreateInvalidJSON(t *testing.T) {
	store := NewStore()
	if err := store.Create([]byte("{bad json")); err == nil {
		t.Fatalf("Create (invalid JSON) expected error, got nil")
	}
}
