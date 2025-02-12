# About Go-note

Go-note is a helper module for creating Obsidian notes with pre-defined templates and metadata availability.

### Key Features

#### 1. **Command-Line Interface**

- Flexible note creation with flags for name (-n), tags (-t), and aliases (-a)
- Help command available (-h, --help)
- Interactive prompts when flags aren't provided

#### 2. **Environment Integration**

- Uses environment variables for configuration
- OBSIDIAN_VAULT: Points to the main vault directory
- OBSIDIAN_TEMPLATES: Points to templates directory

#### 3. **Directory Management**

- Smart directory selection
- Lists and allows selection of numbered directories
- Strictly focuses on directories starting with numbers (a common Obsidian organization pattern)

#### 4. **Template System**

- Dynamic template selection
- Supports multiple markdown templates
- Template variable substitution including:
  - Dates (current, yesterday, tomorrow)
  - Time stamps
  - Title
  - Tags
  - Aliases

#### 5. **Metadata Handling**

- Automatic metadata generation
- Support for YAML frontmatter
- Customizable tags and aliases
- Automatic date and time stamping

#### 6. **Editor Integration**

- Direct integration with Neovim
- Opens in Zen Mode for focused writing
- Cursor positioned for immediate editing

### Technical Implementation

#### 1. **Project Structure**

- Modular design with separate packages:
  - `file`: File operations and directory handling
  - `metadata`: Metadata generation and formatting
  - `template`: Template selection and processing
- Clean separation of concerns

#### 2. **Code Organization**

- Main package for CLI interface
- Internal packages for core functionality
- Modular and maintainable structure

### User Experience

- Simple command-line interface
- Interactive when needed
- Streamlined workflow from creation to editing
- Flexible configuration through environment variables

### Benefits

- Fast note creation
- Consistent note structure
- Integration with existing tools (Obsidian, Neovim)
- Customizable to personal workflow

### Technical Requirements

- Go 1.21
- Neovim
- Obsidian
- Proper environment configuration

# Install and run

To use this program your need to first clone [this repository](https://github.com/peruzzoarthur/go-note).

```bash

git clone https://github.com/peruzzoarthur/go-note.git

```

Then change into the cloned project, run build and move the binary file to your binaries path.

```bash
cd go-note
go build -o note cmd/note/main.go
sudo mv note /usr/local/bin/
```

In order to find the Obsidian vault and be able to import the templates for the .md files, declare the environment variables in your shell config file.

```bash
echo 'export OBSIDIAN_VAULT="/home/user/vault/path"' >> ~/.zshrc # adjust filename for your shell config file
echo 'export OBSIDIAN_TEMPLATES="/home/user/templates/path"' >> ~/.zshrc
```

Done! Go-note is now installed and configured. It can be run using:

```bash

note [options]

# Options:

#   -h, -help            Show this help messaged
#   -n, --name           Add name to the note (string)
#   -t, --tags           Add tags to the note (comma-separated)
#   -a, --aliases        Add aliases to the note (comma-separated)

# Examples:
#   note
#   note -n 'my-new-note' -t 'golang,notes' -a 'go notes'
```

> [!WARNING]
> This module strictly focuses on directories starting with numbers (a common Obsidian organization pattern).
> In other words: 'the directories where you create your notes must start with a number'. e.g. '00-inbox', '1-projects'
