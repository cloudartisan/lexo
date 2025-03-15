# wc - Word Counter

A simple command-line utility that counts words from standard input. Learning
exercise for Go.

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

This will create a `wc` executable in the current directory.

## Usage

```bash
# Count words from stdin (default)
./wc

# Count words from a file
cat file.txt | ./wc

# Count words from a command output
echo "hello world" | ./wc

# Count words explicitly
./wc -w
./wc --words

# Count lines instead of words
./wc -l
./wc --lines

# Count characters instead of words
./wc -c
./wc --chars

# Count lines of code in current directory
./wc --loc

# Count lines of code in a specific directory
./wc --loc /path/to/project

# Count lines of code in multiple directories
./wc --loc dir1 dir2 dir3

# Detect language of text (from stdin)
./wc --lang < file.txt

# Detect language with human-readable name (from stdin)
./wc --lang-name < file.txt

# Detect language of specific file
./wc --lang file.txt

# Detect language of multiple files
./wc --lang file1.txt file2.txt

# Detect language with human-readable name
./wc --lang-name file.txt

# Detect language and count words
./wc --lang -w file.txt
```

## Examples

```bash
# Counting words (default)
echo "The quick brown fox jumps over the lazy dog" | ./wc
9

# Counting lines
echo -e "line 1\nline 2\nline 3" | ./wc -l
3

# Counting characters
echo "hello" | ./wc -c
5

# Counting lines of code
./wc --loc .
125

# Detecting language
echo "The quick brown fox jumps over the lazy dog" | ./wc --lang
Language: en

# Detecting language with human-readable name
echo "El zorro marrón rápido salta sobre el perro perezoso" | ./wc --lang-name
Language: Spanish

# Detecting language and counting words
echo "Le renard brun rapide saute par-dessus le chien paresseux" | ./wc --lang -w
Language: fr
Count: 10

# Detecting language of specific files
./wc --lang main.go README.md
main.go:
Language: en-US
README.md:
Language: en-US

# Detecting language with human-readable name from file
./wc --lang-name README.md
Language: English (US)
```

## Dependencies

This tool has the following dependencies:

### SCC (Sloc, Cloc and Code)

For the `--loc` feature to work, you need to have the `scc` tool installed:

```bash
go install github.com/boyter/scc@latest
```

This will install the [scc (Sloc, Cloc and Code)](https://github.com/boyter/scc) tool, which provides accurate code counting with support for many languages and smart handling of exclusion patterns.

### Whatlanggo

The `--lang` and `--lang-name` features use the [whatlanggo](https://github.com/abadojack/whatlanggo) library for language detection, which supports over 80 languages. This dependency is managed through Go modules and doesn't require separate installation.

## Development

### Requirements
- Go 1.16 or newer
- scc tool (for the --loc feature)

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
```

## License

This project is open source.
