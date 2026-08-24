package output

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteLocal_RoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "exchange.toml")
	data := []byte("hello = \"world\"\n")

	written, err := WriteLocal(path, data)
	if err != nil {
		t.Fatalf("WriteLocal returned error: %v", err)
	}
	if written != path {
		t.Errorf("Expected returned path %q, got %q", path, written)
	}

	readBack, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read back written file: %v", err)
	}
	if string(readBack) != string(data) {
		t.Errorf("Expected file contents %q, got %q", data, readBack)
	}
}

func TestWriteLocal_CreatesParentDirs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "dir", "exchange.toml")
	data := []byte("hello = \"world\"\n")

	if _, err := WriteLocal(path, data); err != nil {
		t.Fatalf("WriteLocal returned error: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("Expected file to exist at %q, got error: %v", path, err)
	}
}
