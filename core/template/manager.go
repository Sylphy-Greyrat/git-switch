package template

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sylphy/git-switch/core/config"
)

func ListTemplates() ([]string, error) {
	dir, err := config.DefaultConfigDir()
	if err != nil {
		return nil, err
	}
	templatesDir := filepath.Join(dir, "templates")
	entries, err := os.ReadDir(templatesDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read templates directory %s: %w", templatesDir, err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}
