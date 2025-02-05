package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/peruzzoarthur/go-note/internal/file"
	"github.com/peruzzoarthur/go-note/internal/metadata"
)

func printHelp() {
	helpText := `
Go-Note - A command-line note-taking tool integrated with Obsidian

Usage:
  note [options] <filename>

Environment Variables Required:
  OBSIDIAN_VAULT       Path to your Obsidian vault directory
  OBSIDIAN_TEMPLATES   Path to your Obsidian templates directory
  
Options:
  -h, -help            Show this help message
  -t, --tags           Add tags to the note (comma-separated)
  -a, --aliases        Add aliases to the note (comma-separated)
  
Description:
  Go-Note helps you create and manage notes in your Obsidian vault.
  It allows you to:
  - Select a directory for your note
  - Choose a template from your Obsidian templates
  - Create a new note with proper metadata
  - Automatically open the note in Neovim with ZenMode

Examples:
  go-note my-new-note
  go-note my-new-note -t "golang,notes" -a "go notes,programming"

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

	// Define all flags
	flag.StringVar(&name, "name", "", "Name of the note")
	flag.StringVar(&name, "n", "", "Name of the note (shorthand)")
	flag.StringVar(&tags, "tags", "", "Tags for the note (comma-separated)")
	flag.StringVar(&tags, "t", "", "Tags for the note (shorthand)")
	flag.StringVar(&aliases, "aliases", "", "Aliases for the note (comma-separated)")
	flag.StringVar(&aliases, "a", "", "Aliases for the note (shorthand)")
	help := flag.Bool("help", false, "Show help message")
	h := flag.Bool("h", false, "Show help message")

	flag.Parse()
	if *help || *h {
		printHelp()
	}

	obsidianDir := os.Getenv("OBSIDIAN_VAULT")
	if obsidianDir == "" {
		fmt.Println("OBSIDIAN_VAULT environment variable not set")
		os.Exit(1)
	}

	var filename string
	var meta metadata.Metadata

	// Handle filename from flags or interactive input
	if name != "" || flag.NArg() > 0 {
		// Use filename from flags or first argument
		if name != "" {
			filename = name
		} else {
			filename = flag.Arg(0)
		}
		meta = metadata.Metadata{
			Title:   strings.ReplaceAll(filename, "-", " "),
			Tags:    make([]string, 0),
			Aliases: make([]string, 0),
		}

		// Process tags and aliases
		if tags != "" {
			meta.Tags = strings.Split(tags, ",")
			for i, tag := range meta.Tags {
				meta.Tags[i] = strings.TrimSpace(tag)
			}
		}

		if aliases != "" {
			meta.Aliases = strings.Split(aliases, ",")
			for i, alias := range meta.Aliases {
				meta.Aliases[i] = strings.TrimSpace(alias)
			}
		}
	} else {
		// If no filename provided via flags or args, get it interactively
		filename = file.GetFilename()
		meta = metadata.Metadata{
			Title:   strings.ReplaceAll(filename, "-", " "),
			Tags:    []string{"tag1", "tag2"},
			Aliases: []string{"aliases1", "aliases2"},
		}
	}

	selectedDir, err := file.SelectDir(obsidianDir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := file.CreateNote(selectedDir, filename, meta); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
