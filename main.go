package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// LanguageSummary represents a summary of a language from scc JSON output
type LanguageSummary struct {
	Name            string `json:"Name"`
	Code            int    `json:"Code"`
	Comment         int    `json:"Comment"`
	Blank           int    `json:"Blank"`
	Complexity      int    `json:"Complexity"`
	Count           int    `json:"Count"`
	WeightedComplex int    `json:"WeightedComplexity"`
}

func countWords(r io.Reader) int {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)

	wc := 0
	for scanner.Scan() {
		wc++
	}

	return wc
}

func countLines(r io.Reader) int {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	lc := 0
	for scanner.Scan() {
		lc++
	}

	return lc
}

func countChars(r io.Reader) int {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanRunes)

	cc := 0
	for scanner.Scan() {
		cc++
	}

	return cc
}

func countLinesOfCode(paths []string) error {
	// Build exclusion pattern
	excludes := []string{
		"--exclude-dir=target",
		"--exclude-dir=node_modules",
		"--exclude-dir=.git",
		"--exclude-dir=.idea",
		"--exclude-dir=.vscode",
		"--exclude-dir=build",
		"--exclude-dir=dist",
		"--exclude-dir=bin",
		"--exclude-dir=obj",
	}

	// Prepare scc command
	args := append([]string{
		"--format=json",
	}, excludes...)
	args = append(args, paths...)

	// Execute scc command
	cmd := exec.Command("scc", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	// Check if scc is installed
	if _, err := exec.LookPath("scc"); err != nil {
		return fmt.Errorf("scc is not installed. Please install it with 'go install github.com/boyter/scc@latest'")
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run scc: %w", err)
	}

	// Parse JSON output
	var summaries []LanguageSummary
	if err := json.Unmarshal(out.Bytes(), &summaries); err != nil {
		return fmt.Errorf("failed to parse scc output: %w", err)
	}

	// Calculate total lines of code
	var total int
	for _, summary := range summaries {
		total += summary.Code
	}

	fmt.Println(total)
	return nil
}

// Config holds the configuration for the program
type Config struct {
	LOC          bool
	Line         bool
	Char         bool
	Word         bool
	Paths        []string
	Input        io.Reader
	Output       io.Writer
	ErrorOutput  io.Writer
}

// NewDefaultConfig creates a default configuration
func NewDefaultConfig() *Config {
	return &Config{
		Input:       os.Stdin,
		Output:      os.Stdout,
		ErrorOutput: os.Stderr,
	}
}

// ParseFlags parses command-line flags and updates the configuration
func ParseFlags(cfg *Config) {
	// Check for help flag manually
	for _, arg := range os.Args[1:] {
		if arg == "-h" || arg == "--help" {
			fmt.Fprintf(cfg.ErrorOutput, "Usage: %s [flags] [path...]\n\n", os.Args[0])
			fmt.Fprintf(cfg.ErrorOutput, "Count words, lines, characters, or lines of code.\n")
			fmt.Fprintf(cfg.ErrorOutput, "By default, counts words from stdin.\n\n")
			fmt.Fprintf(cfg.ErrorOutput, "Options:\n")
			fmt.Fprintf(cfg.ErrorOutput, "  -w, --words     Count words (default behavior)\n")
			fmt.Fprintf(cfg.ErrorOutput, "  -l, --lines     Count lines instead of words\n")
			fmt.Fprintf(cfg.ErrorOutput, "  -c, --chars     Count characters instead of words\n")
			fmt.Fprintf(cfg.ErrorOutput, "      --loc       Count lines of code instead of words\n")
			fmt.Fprintf(cfg.ErrorOutput, "  -h, --help      Show this help message\n")
			os.Exit(0)
		}
	}
	
	// Define flags
	var loc bool
	var l, c, w bool
	var paths []string
	
	// Process args to handle GNU-style long options
	for _, arg := range os.Args[1:] {
		// Process flags
		switch arg {
		case "--loc":
			loc = true
			continue
		case "-l", "--lines":
			l = true
			continue
		case "-c", "--chars":
			c = true
			continue
		case "-w", "--words":
			w = true
			continue
		}
		
		// Handle non-flag arguments (paths for --loc)
		if loc && !strings.HasPrefix(arg, "-") {
			paths = append(paths, arg)
			continue
		}
	}
	
	// Update the configuration
	cfg.LOC = loc
	cfg.Line = l
	cfg.Char = c
	cfg.Word = w || (!cfg.Line && !cfg.Char && !cfg.LOC)
	
	// Set paths for LOC feature
	if loc {
		if len(paths) > 0 {
			cfg.Paths = paths
		} else {
			cfg.Paths = []string{"."}
		}
	}
}

// Run executes the program with the given configuration
func Run(cfg *Config) error {
	// LOC flag takes precedence
	if cfg.LOC {
		if err := countLinesOfCode(cfg.Paths); err != nil {
			return err
		}
		return nil
	}
	
	// Handle standard counting options
	var count int
	switch {
	case cfg.Line:
		count = countLines(cfg.Input)
	case cfg.Char:
		count = countChars(cfg.Input)
	case cfg.Word:
		count = countWords(cfg.Input)
	}
	
	fmt.Fprintln(cfg.Output, count)
	return nil
}

func main() {
	// Create default configuration
	cfg := NewDefaultConfig()
	
	// Parse command-line flags
	ParseFlags(cfg)
	
	// Run the program
	if err := Run(cfg); err != nil {
		fmt.Fprintf(cfg.ErrorOutput, "Error: %v\n", err)
		os.Exit(1)
	}
}
