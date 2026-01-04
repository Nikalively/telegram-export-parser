package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSingleFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "result.json")
	jsonData := `[
		{
			"id": 1,
			"from_id": "123456",
			"text": "Hello",
			"date": "2026-01-01T00:00:00Z",
			"chat": {"id": 100},
			"entities": [{"type": "mention", "text": "@user", "offset": 0, "length": 5}]
		}
	]`
	if err := os.WriteFile(filePath, []byte(jsonData), 0644); err != nil {
		t.Fatal(err)
	}

	events, err := ParseSingleFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}
	if events[0].ID != 1 || events[0].FromID != "123456" {
		t.Errorf("Event data mismatch")
	}
}

func TestMergeFromFolders(t *testing.T) {
	tempDir := t.TempDir()
	folder1 := filepath.Join(tempDir, "export1")
	folder2 := filepath.Join(tempDir, "export2")
	os.MkdirAll(folder1, 0755)
	os.MkdirAll(folder2, 0755)

	jsonData1 := `[{"id": 1, "from_id": "123", "text": "Msg1", "date": "2026-01-01T00:00:00Z", "chat": {"id": 100}}]`
	jsonData2 := `[{"id": 2, "from_id": "456", "text": "Msg2", "date": "2026-01-02T00:00:00Z", "chat": {"id": 100}}]`

	os.WriteFile(filepath.Join(folder1, "result.json"), []byte(jsonData1), 0644)
	os.WriteFile(filepath.Join(folder2, "result.json"), []byte(jsonData2), 0644)

	merged, err := MergeFromFolders([]string{folder1, folder2})
	if err != nil {
		t.Fatal(err)
	}
	if len(merged.Events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(merged.Events))
	}
	if !merged.Events[0].Date.Before(merged.Events[1].Date) {
		t.Errorf("Events not sorted by date")
	}
}
