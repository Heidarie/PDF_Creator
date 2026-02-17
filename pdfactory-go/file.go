package pdfactory

import (
	"os"
	"path/filepath"
)

func writeFile(path string, data []byte) error {
	if path == "" {
		return os.ErrInvalid
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
