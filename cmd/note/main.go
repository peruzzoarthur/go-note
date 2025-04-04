package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	catppuccingo "github.com/catppuccin/go"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/peruzzoarthur/go-note/internal/file"
	"github.com/peruzzoarthur/go-note/internal/metadata"
)

func main() {

	// Ensure obsidian vault dir is set
	obsidianDir := os.Getenv("OBSIDIAN_VAULT")
	if obsidianDir == "" {
		fmt.Println("OBSIDIAN_VAULT environment variable not set")
		os.Exit(1)
	}

	// Ensure templates dir is set
	obsidianTemplatesDir := os.Getenv("OBSIDIAN_TEMPLATES")
	if obsidianTemplatesDir == "" {
		fmt.Println("OBSIDIAN_TEMPLATES environment variable not set")
		os.Exit(1)
	}

	// Get editor of choice (if empty default to nvim)
	editorPreference := os.Getenv("COLD_NOTE_EDITOR")
	if editorPreference == "" {
		editorPreference = "nvim"
	}

	// Initialize metadata with empty values
	meta := metadata.Metadata{
		Tags:    []string{},
		Aliases: []string{},
	}

	// Variables to store user selections
	var (
		filename     string
		selectedDir  string
		templateName string
		tagsInput    []string
		aliasesInput string
	)

	// Log styles
	greenStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(catppuccingo.Mocha.Green().Hex)).
		Bold(true)

	lavenderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(catppuccingo.Mocha.Lavender().Hex)).
		Bold(true)

	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(catppuccingo.Mocha.Red().Hex)).Bold(true)
	// Get directories for selection
	dirs, err := file.GetDirectories(obsidianDir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Get templates for selection
	templates, err := file.GetTemplates(obsidianTemplatesDir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Create directory options for the select menu
	dirOptions := make([]huh.Option[string], len(dirs))
	for i, dir := range dirs {
		dirOptions[i] = huh.NewOption(dir, dir)
	}

	// Create template options for the select menu
	templateOptions := make([]huh.Option[string], len(templates))
	for i, tmpl := range templates {
		templateOptions[i] = huh.NewOption(tmpl, tmpl)
	}

	// Load tags from JSON file
	tagsFilePath := filepath.Join(obsidianDir, "tags.json") // Adjust path as needed
	tagOptions, err := file.LoadTagsFromJSON(tagsFilePath)
	if err != nil {
		// If there's an error loading tags, fall back to default tags
		fmt.Printf(redStyle.Render("Warning: Could not load tags from JSON: %v\n"), err)
		os.Exit(1)
	}

	var catppuccin *huh.Theme = huh.ThemeCatppuccin()
	// Create the form
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Note Filename").
				Description("Enter a name for your note").
				Placeholder("brand-new-note").
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("filename cannot be empty")
					}
					return nil
				}).
				Value(&filename),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Directory").
				Description("Choose directory to save your note").
				Options(dirOptions...).
				Value(&selectedDir),

			huh.NewSelect[string]().
				Title("Select Template").
				Description("Choose a template for your note").
				Options(templateOptions...).
				Value(&templateName),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Tags").
				Description("Select tags for your note").
				Options(tagOptions...).
				Value(&tagsInput),

			huh.NewInput().
				Title("Aliases").
				Description("Enter comma-separated aliases").
				Placeholder("go notes,programming").
				Value(&aliasesInput),
		),
	).WithTheme(catppuccin)

	err = form.Run()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Process the form data
	meta.Title = filename

	// Tags are already processed as a slice
	meta.Tags = tagsInput

	// Process aliases
	if aliasesInput != "" {
		for _, alias := range strings.Split(aliasesInput, ",") {
			trimmedAlias := strings.TrimSpace(alias)
			if trimmedAlias != "" {
				meta.Aliases = append(meta.Aliases, trimmedAlias)
			}
		}
	}

	// Full directory path
	fullDirPath := filepath.Join(obsidianDir, selectedDir)

	// Full template path
	templatePath := filepath.Join(obsidianTemplatesDir, templateName)

	// Read template content
	templateContent, err := file.ReadTemplateContent(templatePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Create note creation function for the spinner
	createNote := func() {
		// Create the note file
		fullPath := filepath.Join(fullDirPath, filename+".md")
		if _, err := os.Stat(fullPath); err == nil {
			fmt.Printf("Error: file already exists: %s\n", fullPath)
			os.Exit(1)
		}

		file, err := os.Create(fullPath)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		content := metadata.FormatMetadata(templateContent, meta)
		if _, err := file.WriteString(content); err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			os.Exit(1)
		}
	}

	_ = spinner.New().
		Title("Creating your note...").
		Action(createNote).
		Run()

	createdFilePath := filepath.Join(fullDirPath, filename+".md")
	fmt.Println(greenStyle.Render(fmt.Sprintf("\nCreated note at %s", createdFilePath)))

	// Open the note with the selected editor
	var cmd *exec.Cmd

	switch editorPreference {
	case "nvim-zen":
		fmt.Println(lavenderStyle.Render("Opening note with Neovim (ZenMode)"))
		cmd = exec.Command("nvim", "+ normal ggzzi", createdFilePath, "-c", ":ZenMode")
	case "nvim":
		fmt.Println(lavenderStyle.Render("Opening note with Neovim"))
		cmd = exec.Command("nvim", "+ normal ggzzi", createdFilePath)
	case "vscode":
		fmt.Println(lavenderStyle.Render("Opening note with VSCode"))
		cmd = exec.Command("code", createdFilePath)
	default:
		fmt.Println(lavenderStyle.Render("Opening note with Neovim + ZenMode (default)"))
		cmd = exec.Command("nvim", "+ normal ggzzi", createdFilePath)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
