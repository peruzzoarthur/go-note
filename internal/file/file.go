package file

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/peruzzoarthur/go-note/internal/metadata"
	"github.com/peruzzoarthur/go-note/internal/template"
)

func GetFilename() string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter a filename: ")
		filename, _ := reader.ReadString('\n')
		trimmed := strings.TrimSpace(filename)

		if len(trimmed) > 0 {
			return trimmed
		}

		fmt.Println("Error: The filename cannot be empty. Please insert a proper value.")
	}
}

func SelectDir(zetDir string) (string, error) {
	entries, err := os.ReadDir(zetDir)
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
		return "", fmt.Errorf("no directories found in %s", zetDir)
	}

	fmt.Println("\nAvailable directories:")
	for i, dir := range dirs {
		fmt.Printf("%d: %s\n", i+1, dir)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nSelect directory number: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("error reading input: %w", err)
		}

		var selection int
		_, err = fmt.Sscanf(strings.TrimSpace(input), "%d", &selection)
		if err != nil || selection < 1 || selection > len(dirs) {
			fmt.Printf("Please enter a number between 1 and %d\n", len(dirs))
			continue
		}

		return filepath.Join(zetDir, dirs[selection-1]), nil
	}
}

func CreateNote(directory string, filename string, meta metadata.Metadata) error {
	if len(filename) == 0 {
		return fmt.Errorf("please insert a filename")
	}

	fullPath := filepath.Join(directory, filename+".md")

	if template.FileExists(fullPath) {
		return fmt.Errorf("file already exists: %s", fullPath)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	obsidianTemplates := os.Getenv("OBSIDIAN_TEMPLATES")
	if obsidianTemplates == "" {
		fmt.Println("OBSIDIAN_TEMPLATES environment variable not set")
		os.Exit(1)
	}

	tmpl, err := template.ReadTemplate(obsidianTemplates)
	if err != nil {
		return fmt.Errorf("error getting template %w", err)
	}

	content := metadata.FormatMetadata(tmpl, meta)

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	cmd := exec.Command("nvim", "+ normal ggzzi", fullPath, "-c", ":ZenMode")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
