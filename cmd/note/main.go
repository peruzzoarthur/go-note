package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/peruzzoarthur/go-note/internal/file"
	"github.com/peruzzoarthur/go-note/internal/metadata"
)

func printHelp() {
	helpText := `
Go-Note - A command-line note-taking tool integrated with Obsidian

Usage:
  note [options] 

Environment Variables Required:
  OBSIDIAN_VAULT       Path to your Obsidian vault directory
  OBSIDIAN_TEMPLATES   Path to your Obsidian templates directory
  
Options:
  -h, -help            Show this help messaged
  -n, --name           Add name to the note (string)
  -t, --tags           Add tags to the note (comma-separated)
  -a, --aliases        Add aliases to the note (comma-separated)
  
Description:
  Go-Note helps you create and manage notes in your Obsidian vault.
  It allows you to:
  - Select a directory for your note
  - Choose a template from your templates
  - Create a new note with proper metadata
  - Automatically open the note in your text editor

Examples:
  note
  note -n 'my-new-note' -t 'golang,notes' -a 'go notes'

Note:
  Make sure your environment variables are properly set before running.
`
	fmt.Println(helpText)
	os.Exit(0)
}

func main() {
	var (
		name    string
		tags    string
		aliases string
	)

	// Define flags
	flag.StringVar(&name, "name", "", "Name of the note")
	flag.StringVar(&name, "n", "", "Name of the note (shorthand)")
	flag.StringVar(&tags, "tags", "", "Tags for the note (comma-separated)")
	flag.StringVar(&tags, "t", "", "Tags for the note (shorthand)")
	flag.StringVar(&aliases, "aliases", "", "Aliases for the note (comma-separated)")
	flag.StringVar(&aliases, "a", "", "Aliases for the note (shorthand)")
	help := flag.Bool("help", false, "Show help message")
	h := flag.Bool("h", false, "Show help message")

	flag.Parse()

	// Show help if requested
	if *help || *h {
		printHelp()
	}

	// Ensure OBSIDIAN_VAULT is set
	obsidianDir := os.Getenv("OBSIDIAN_VAULT")
	if obsidianDir == "" {
		fmt.Println("OBSIDIAN_VAULT environment variable not set")
		os.Exit(1)
	}

	obsidianTemplatesDir := os.Getenv("OBSIDIAN_TEMPLATES")
	if obsidianTemplatesDir == "" {
		fmt.Println("OBSIDIAN_TEMPLATES environment variable not set")
		os.Exit(1)
	}
	// Initialize metadata
	meta := metadata.Metadata{
		Tags:    []string{"tag1", "tag2", "tag3"},
		Aliases: []string{"aliases1", "aliases2"},
	}

	filename := strings.TrimSpace(name)

	if filename == "" {
		filename = file.GetFilename()
	}

	meta.Title = strings.ReplaceAll(filename, "-", " ")

	if tags != "" {
		meta.Tags = []string{}
		meta.Tags = append(meta.Tags, strings.Split(tags, ",")...)
	}

	if aliases != "" {
		meta.Aliases = []string{}
		meta.Aliases = append(meta.Aliases, strings.Split(aliases, ",")...)
	}

	// Select the directory and create the note
	selectedDir, err := file.SelectDir(obsidianDir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	createdFilePath, err := file.CreateNote(selectedDir, filename, meta, obsidianTemplatesDir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nCreated note at %s\n", createdFilePath)

	cmd := exec.Command("nvim", "+ normal ggzzi", createdFilePath, "-c", ":ZenMode")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Print("\nOpening note...\n")
	cmd.Run()
}
