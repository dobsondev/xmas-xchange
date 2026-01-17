package participants

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetParticipantsFromToml_ValidFile(t *testing.T) {
	// Create a temporary TOML file
	content := `[[participant]]
name = "Alice"
number = "+15551000001"
restrictions = ["Bob"]

[[participant]]
name = "Bob"
number = "+15551000002"
restrictions = []

[[participant]]
name = "Carol"
number = "+15551000003"
restrictions = ["Alice", "Bob"]
`

	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test_participants.toml")
	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Test reading the file
	participants := GetParticipantsFromToml(&tempFile)

	// Verify the results
	if len(participants) != 3 {
		t.Errorf("Expected 3 participants, got %d", len(participants))
	}

	// Check first participant
	if participants[0].Name != "Alice" {
		t.Errorf("Expected first participant name to be 'Alice', got '%s'", participants[0].Name)
	}
	if participants[0].Number != "+15551000001" {
		t.Errorf("Expected first participant number to be '+15551000001', got '%s'", participants[0].Number)
	}
	if len(participants[0].Restrictions) != 1 || participants[0].Restrictions[0] != "Bob" {
		t.Errorf("Expected first participant restrictions to be ['Bob'], got %v", participants[0].Restrictions)
	}

	// Check second participant
	if participants[1].Name != "Bob" {
		t.Errorf("Expected second participant name to be 'Bob', got '%s'", participants[1].Name)
	}
	if len(participants[1].Restrictions) != 0 {
		t.Errorf("Expected second participant to have no restrictions, got %v", participants[1].Restrictions)
	}

	// Check third participant
	if participants[2].Name != "Carol" {
		t.Errorf("Expected third participant name to be 'Carol', got '%s'", participants[2].Name)
	}
	if len(participants[2].Restrictions) != 2 {
		t.Errorf("Expected third participant to have 2 restrictions, got %d", len(participants[2].Restrictions))
	}
}

func TestGetParticipantsFromToml_NoRestrictions(t *testing.T) {
	// Create a temporary TOML file with participants that have no restrictions
	content := `[[participant]]
name = "Dave"
number = "+15551000004"
restrictions = []

[[participant]]
name = "Eve"
number = "+15551000005"
restrictions = []
`

	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test_participants.toml")
	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	participants := GetParticipantsFromToml(&tempFile)

	if len(participants) != 2 {
		t.Errorf("Expected 2 participants, got %d", len(participants))
	}

	for i, p := range participants {
		if len(p.Restrictions) != 0 {
			t.Errorf("Expected participant %d to have no restrictions, got %v", i, p.Restrictions)
		}
	}
}

func TestGetParticipantsFromToml_SingleParticipant(t *testing.T) {
	// Create a temporary TOML file with a single participant
	content := `[[participant]]
name = "Frank"
number = "+15551000006"
restrictions = ["Grace", "Henry"]
`

	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test_participants.toml")
	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	participants := GetParticipantsFromToml(&tempFile)

	if len(participants) != 1 {
		t.Errorf("Expected 1 participant, got %d", len(participants))
	}

	if participants[0].Name != "Frank" {
		t.Errorf("Expected participant name to be 'Frank', got '%s'", participants[0].Name)
	}

	if len(participants[0].Restrictions) != 2 {
		t.Errorf("Expected 2 restrictions, got %d", len(participants[0].Restrictions))
	}
}

func TestGetParticipantsFromToml_EmptyFile(t *testing.T) {
	// Create an empty TOML file
	content := ``

	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test_participants.toml")
	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	participants := GetParticipantsFromToml(&tempFile)

	if len(participants) != 0 {
		t.Errorf("Expected 0 participants, got %d", len(participants))
	}
}
