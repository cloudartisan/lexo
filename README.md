# lexo - Text Analysis Utility

A versatile command-line utility for text and code analysis. Evolved from the classic Unix `wc` command with enhanced features for modern text processing, including language detection and word frequency analysis.

## Installation

### Prerequisites

First, install Go on your system:

#### macOS
```bash
# Using Homebrew
brew install go

# Or download from the official website
# https://golang.org/dl/
```

#### Linux
```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install golang

# Fedora
sudo dnf install golang

# Or download from the official website
# https://golang.org/dl/
```

#### Windows
```bash
# Using Chocolatey
choco install golang

# Or download the installer from
# https://golang.org/dl/
```

### Building the Application
Once Go is installed:

```bash
go build
```

This will create a `lexo` executable in the current directory.

### Installation
To install lexo globally on your system:

```bash
go install
```

## Usage

```bash
# Count words from stdin (default)
lexo

# Count words from a file
cat file.txt | lexo

# Count words from a command output
echo "hello world" | lexo

# Count words explicitly
lexo -w
lexo --words

# Count lines instead of words
lexo -l
lexo --lines

# Count characters instead of words
lexo -c
lexo --chars

# Count lines of code in current directory
lexo --loc

# Count lines of code in a specific directory
lexo --loc /path/to/project

# Count lines of code in multiple directories
lexo --loc dir1 dir2 dir3

# Detect language of text (from stdin)
lexo --lang < file.txt

# Detect language with human-readable name (from stdin)
lexo --lang-name < file.txt

# Detect language of specific file
lexo --lang file.txt

# Detect language of multiple files
lexo --lang file1.txt file2.txt

# Detect language with human-readable name
lexo --lang-name file.txt

# Detect language and count words
lexo --lang -w file.txt

# Analyze word frequency (alphabetical order)
lexo --freq file.txt

# Analyze word frequency, sorted by count (most frequent first)
lexo --freq --sort-count file.txt

# Limit frequency results to top N words
lexo --freq --sort-count --limit 5 file.txt

# Analyze multiple files
lexo --freq file1.txt file2.txt
```

## Examples

```bash
# Counting words (default)
echo "The quick brown fox jumps over the lazy dog" | lexo
9

# Counting lines
echo -e "line 1\nline 2\nline 3" | lexo -l
3

# Counting characters
echo "hello" | lexo -c
5

# Counting lines of code
lexo --loc .
125

# Detecting language
echo "The quick brown fox jumps over the lazy dog" | lexo --lang
Language: en-US

# Detecting language with human-readable name
echo "El zorro marrón rápido salta sobre el perro perezoso" | lexo --lang-name
Language: Spanish (Spain)

# Detecting language and counting words
echo "Le renard brun rapide saute par-dessus le chien paresseux" | lexo --lang -w
Language: fr
Count: 10

# Detecting language of specific files
lexo --lang main.go README.md
main.go:
Language: fr
README.md:
Language: en-US

# Detecting language with human-readable name from file
lexo --lang-name README.md
Language: English (US)

# Word frequency analysis
lexo --freq --sort-count --limit 5 README.md
Word frequency (sorted by count):
-----------  ------
the          15
to           12
language     8
of           7
for          6

# Alphabetical word frequency
lexo --freq --limit 5 README.md
Word frequency (sorted alphabetically):
-----------  ------
a            5
about        1
accurately   1
add          1
all          3
```

## Dependencies

This tool has the following dependencies:

### Whatlanggo

The `--lang` and `--lang-name` features use the [whatlanggo](https://github.com/abadojack/whatlanggo) library for language detection, which supports over 80 languages. This dependency is managed through Go modules and doesn't require separate installation.

## Development

### Requirements
- Go 1.16 or newer

### Build and Test

```bash
# Build the application
go build

# Run all tests
go test ./...

# Run specific test
go test -run TestCountWords

# Check test coverage
go test -cover ./...

# Generate HTML coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Test Coverage

The project maintains comprehensive test coverage to ensure reliability and stability. Tests cover all major functionality including:

- Word, line, and character counting
- Language detection
- Word frequency analysis
- Command-line flag parsing
- Error handling

## License

This project is open source.
