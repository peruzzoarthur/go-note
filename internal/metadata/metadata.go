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

	replacements := map[string]string{
		"{{date:YYYYMMDD}}":   now.Format("20060102"),
		"{{time:HHmm}}":       now.Format("1504"),
		"{{date:YYYY-MM-DD}}": now.Format("2006-01-02"),
		"{{title}}":           metadata.Title,
		"{{tags}}":            strings.Join(metadata.Tags, ", "),
		"{{alias}}":           strings.Join(metadata.Aliases, ", "),
	}

	for placeholder, value := range replacements {
		content = strings.ReplaceAll(content, placeholder, value)
	}

	return content
}

func ParseArgs(args []string) (string, Metadata) {
	metadata := Metadata{
		Tags:    make([]string, 0),
		Aliases: make([]string, 0),
	}

	filename := ""
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-t" || arg == "--tags":
			if i+1 < len(args) {
				i++
				metadata.Tags = strings.Split(args[i], ",")
				for i, tag := range metadata.Tags {
					metadata.Tags[i] = strings.TrimSpace(tag)
				}
			}
		case arg == "-a" || arg == "--aliases":
			if i+1 < len(args) {
				i++
				metadata.Aliases = strings.Split(args[i], ",")
				for i, alias := range metadata.Aliases {
					metadata.Aliases[i] = strings.TrimSpace(alias)
				}
			}
		default:
			if filename == "" {
				filename = arg
				metadata.Title = strings.ReplaceAll(arg, "-", " ")
			}
		}
	}

	return filename, metadata
}
