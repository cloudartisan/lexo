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
# Count words from stdin
./wc

# Count words from a file
cat file.txt | ./wc

# Count words from a command output
echo "hello world" | ./wc

# Count lines of code in current directory
./wc --loc

# Count lines of code in a specific directory
./wc --loc /path/to/project

# Count lines of code in multiple directories
./wc --loc dir1 dir2 dir3
```

## Examples

```bash
# Counting words
echo "The quick brown fox jumps over the lazy dog" | ./wc
9

# Counting lines of code
./wc --loc .
125
```

## Dependencies

For the `--loc` feature to work, you need to have the `scc` tool installed:

```bash
go install github.com/boyter/scc@latest
```

This will install the [scc (Sloc, Cloc and Code)](https://github.com/boyter/scc) tool, which provides accurate code counting with support for many languages and smart handling of exclusion patterns.

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
