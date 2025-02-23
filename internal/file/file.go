package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/peruzzoarthur/go-note/internal/metadata"
	"github.com/peruzzoarthur/go-note/internal/template"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func GetFilename() string {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter a filename: ")

		scanned := scanner.Scan()
		if scanned {
			text := strings.TrimSpace(scanner.Text())
			if len(text) > 0 {
				return text
			}
		}
	}
}

func SelectDir(obsidianDir string) (string, error) {
	entries, err := os.ReadDir(obsidianDir)
	if err != nil {
		return "", fmt.Errorf("error reading directory: %w", err)
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			if len(name) > 0 && (name[0] >= '0' && name[0] <= '9') {
				dirs = append(dirs, name)
			}
		}
	}

	if len(dirs) == 0 {
		return "", fmt.Errorf("no directories found in %s", obsidianDir)
	}

	fmt.Println("\nAvailable directories:")
	for i, dir := range dirs {
		fmt.Printf("%d: %s\n", i+1, dir)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nSelect dir number: ")

		scanned := scanner.Scan()

		if !scanned {
			return "", fmt.Errorf("invalid input")
		}

		text := scanner.Text()

		var selected int

		_, err := fmt.Sscanf(strings.TrimSpace(text), "%d", &selected)

		if err != nil || selected < 1 || selected > len(dirs) {
			fmt.Printf("Please enter a number between 1 and %d\n", len(dirs))
			continue
		}
		return filepath.Join(obsidianDir, dirs[selected-1]), nil
	}
}

func CreateNote(directory string, filename string, meta metadata.Metadata, templatesPath string) (string, error) {
	fullPath := filepath.Join(directory, filename+".md")
	if fileExists(fullPath) {
		return "", fmt.Errorf("file already exists: %s", fullPath)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %w", err)
	}

	defer file.Close()
	tmpl, err := template.ReadTemplate(templatesPath)
	if err != nil {
		return "", fmt.Errorf("error getting template %w", err)
	}

	content := metadata.FormatMetadata(tmpl, meta)
	if _, err := file.WriteString(content); err != nil {
		return "", fmt.Errorf("error writing to file: %w", err)
	}

	return fullPath, nil
}
