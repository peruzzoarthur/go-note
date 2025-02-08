package template

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func selectTemplate(templatesDir string) (string, error) {
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return "", fmt.Errorf("error reading templates directory: %w", err)
	}

	var templates []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			templates = append(templates, entry.Name())
		}
	}

	if len(templates) == 0 {
		return "", fmt.Errorf("no template files found in %s", templatesDir)
	}

	fmt.Println("\nAvailable templates:")
	for i, tmpl := range templates {
		fmt.Printf("%d: %s\n", i+1, tmpl)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nSelect template number: ")

		scanned := scanner.Scan()
		if !scanned {
			return "", fmt.Errorf("invalid input")
		}

		input := scanner.Text()
		var selection int
		_, err = fmt.Sscanf(strings.TrimSpace(input), "%d", &selection)
		if err != nil || selection < 1 || selection > len(templates) {
			fmt.Printf("Please enter a number between 1 and %d\n", len(templates))
			continue
		}

		return filepath.Join(templatesDir, templates[selection-1]), nil
	}
}

func ReadTemplate(templatesDir string) (string, error) {
	// Expand any environment variables in the path
	expandedPath := os.ExpandEnv(templatesDir)

	// If path starts with ~, replace with home directory
	if strings.HasPrefix(expandedPath, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("error getting home directory: %w", err)
		}
		expandedPath = filepath.Join(home, expandedPath[2:])
	}

	if !FileExists(expandedPath) {
		return "", fmt.Errorf("templates directory not found at: %s", expandedPath)
	}

	templatePath, err := selectTemplate(expandedPath)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("error reading template: %w", err)
	}

	return string(content), nil
}
