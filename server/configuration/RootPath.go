package configuration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	ProjectPath  = "./projects"
	PistacheRoot = "./pistache"
	TemplateRoot = "./front"
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

func init() {
	log.Println("Loading paths...")
	var err error
	PistacheRoot, err = GetRootPath("pistache")
	if err != nil {
		log.Println("'pistache' folder could not be found. pistache will be disabled")
	} else {
		log.Printf("Pistache root set to %s", PistacheRoot)
	}
	if os.Getenv("MODE") == "dev" {
		TemplateRoot = "../front/"
	}
	log.Println("Template root set to", TemplateRoot)
	ProjectPath, err = GetRootPath("projects")
	if err != nil {
		log.Println("'projects' folder could not be found. project management will be disabled")
	} else {
		log.Printf("Monitoring projects from %s", ProjectPath)
	}
}
