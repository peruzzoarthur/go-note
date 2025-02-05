package metadata

import (
	"strings"
	"time"
)

type Metadata struct {
	Title   string
	Tags    []string
	Aliases []string
}

func FormatMetadata(content string, metadata Metadata) string {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)

	replacements := map[string]string{
		"{{date:YYYYMMDD}}":        now.Format("20060102"),
		"{{time:HHmm}}":            now.Format("1504"),
		"{{date:YYYY-MM-DD}}":      now.Format("2006-01-02"),
		"{{yesterday:YYYY-MM-DD}}": yesterday.Format("2006-01-02"),
		"{{tomorrow:YYYY-MM-DD}}":  tomorrow.Format("2006-01-02"),
		"{{title}}":                metadata.Title,
		"{{tags}}":                 strings.Join(metadata.Tags, ", "),
		"{{alias}}":                strings.Join(metadata.Aliases, ", "),
	}

	for placeholder, value := range replacements {
		content = strings.ReplaceAll(content, placeholder, value)
	}

	return content
}

// func ParseArgs(filename, tags, aliases string) (string, Metadata) {
// 	// Initialize metadata
// 	metadata := Metadata{
// 		Tags:    make([]string, 0),
// 		Aliases: make([]string, 0),
// 	}
//
// 	// Handle filename
// 	if filename == "" && flag.NArg() > 0 {
// 		// If no -n/--name flag, use the first positional argument
// 		filename = flag.Arg(0)
// 	}
//
// 	if filename == "" {
// 		fmt.Println("Error: filename is required. Use -n/--name flag or provide it as an argument")
// 		flag.Usage()
// 		os.Exit(1)
// 	}
//
// 	// Set title (replace hyphens with spaces)
// 	metadata.Title = strings.ReplaceAll(filename, "-", " ")
//
// 	// Process tags if provided
// 	if tags != "" {
// 		metadata.Tags = strings.Split(tags, ",")
// 		for i, tag := range metadata.Tags {
// 			metadata.Tags[i] = strings.TrimSpace(tag)
// 		}
// 	}
//
// 	// Process aliases if provided
// 	if aliases != "" {
// 		metadata.Aliases = strings.Split(aliases, ",")
// 		for i, alias := range metadata.Aliases {
// 			metadata.Aliases[i] = strings.TrimSpace(alias)
// 		}
// 	}
//
// 	return filename, metadata
// }
// func ParseArgs(args []string) (string, Metadata) {
// 	metadata := Metadata{
// 		Tags:    make([]string, 0),
// 		Aliases: make([]string, 0),
// 	}
//
// 	filename := ""
// 	for i := 0; i < len(args); i++ {
// 		arg := args[i]
// 		switch {
// 		case arg == "-t" || arg == "--tags":
// 			if i+1 < len(args) {
// 				i++
// 				metadata.Tags = strings.Split(args[i], ",")
// 				for i, tag := range metadata.Tags {
// 					metadata.Tags[i] = strings.TrimSpace(tag)
// 				}
// 			}
// 		case arg == "-a" || arg == "--aliases":
// 			if i+1 < len(args) {
// 				i++
// 				metadata.Aliases = strings.Split(args[i], ",")
// 				for i, alias := range metadata.Aliases {
// 					metadata.Aliases[i] = strings.TrimSpace(alias)
// 				}
// 			}
// 		case arg == "-n" || arg == "--name":
// 			if i+1 < len(args) {
// 				i++
// 				filename = arg
// 				metadata.Title = strings.ReplaceAll(arg, "-", " ")
// 			} else {
// 				fmt.Print("Pass correct number of names to -n flag")
// 				os.Exit(1)
// 			}
// 		default:
// 			if filename == "" {
// 				filename = arg
// 				metadata.Title = strings.ReplaceAll(arg, "-", " ")
// 			}
// 		}
// 	}
//
// 	return filename, metadata
// }
