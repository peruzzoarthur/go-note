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

	args := os.Args[1:]
	var filename string
	var meta metadata.Metadata

	if len(args) == 0 {
		filename = file.GetFilename()
		meta = metadata.Metadata{
			Title: strings.ReplaceAll(filename, "-", " "),
		}
	} else {
		filename, meta = metadata.ParseArgs(args)
		if filename == "" {
			fmt.Println("Error: Filename is required")
			fmt.Println("Usage: note <filename> [-t|--tags 'tag1,tag2'] [-a|--aliases 'alias1,alias2']")
			fmt.Println("Example: note my-new-note -t 'golang,notes' -a 'go notes,notes'")
			os.Exit(1)
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
