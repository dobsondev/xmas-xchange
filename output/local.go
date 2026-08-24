package output

import (
	"os"
	"path/filepath"
)

// WriteLocal writes data to path, creating parent directories as needed, and returns path.
func WriteLocal(path string, data []byte) (string, error) {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return "", err
		}
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}

	return path, nil
}
