package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/peruzzoarthur/go-note/internal/file"
	"github.com/peruzzoarthur/go-note/internal/metadata"
)

func main() {
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
