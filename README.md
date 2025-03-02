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
```

## Example

```bash
echo "The quick brown fox jumps over the lazy dog" | ./wc
9
```

## Development

### Requirements
- Go 1.15 or newer

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
