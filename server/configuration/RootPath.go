package configuration

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetRootPath returns the correct path of a folder, checking cwd or parent directory (1 level)
func GetRootPath(folderName string) (string, error) {
	// Check if the folder exists in the current directory
	if _, err := os.Stat(folderName); err == nil {
		return folderName, nil
	}
	candidate := filepath.Join("..", folderName)
	// Otherwise, check the parent directory
	if _, err := os.Stat(candidate); err == nil {
		return candidate, nil
	}
	// Otherwise, return an error
	return "", fmt.Errorf("folder %s not found", folderName)
}
