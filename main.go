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

	"github.com/abadojack/whatlanggo"
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

// detectLanguage tries to detect the language of the text
// and returns the language tag (e.g., en-US, es, fr) and a human-readable name
func detectLanguage(r io.Reader) (string, string, error) {
	// We need to read the text into memory to process it
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	
	// Read all the text (up to a reasonable limit)
	// This gives better accuracy than just a small sample
	scanner := bufio.NewScanner(tee)
	scanner.Split(bufio.ScanWords)
	
	var sample strings.Builder
	wordCount := 0
	const maxWords = 1000 // Reasonable limit to avoid memory issues with very large files
	
	for scanner.Scan() && wordCount < maxWords {
		if wordCount > 0 {
			sample.WriteString(" ")
		}
		sample.WriteString(scanner.Text())
		wordCount++
	}
	
	if err := scanner.Err(); err != nil {
		return "", "", fmt.Errorf("error reading text: %w", err)
	}
	
	// If we didn't get any words, we can't detect the language
	if wordCount == 0 {
		return "und", "Unknown", nil
	}
	
	// Use whatlanggo for accurate language detection
	text := sample.String()
	// No special options needed - the default algorithm is already quite good
	info := whatlanggo.Detect(text)
	
	// Get the ISO language code
	langTag := info.Lang.Iso6391()
	
	// Get the English name of the language
	langName := info.Lang.String()
	
	// If the language is unknown, fall back to a sensible default
	if langTag == "" {
		return "und", "Unknown", nil
	}
	
	// For certain languages with common regional variants, add region code
	// This is just an example - in a real system this would be more sophisticated
	switch langTag {
	case "en":
		// For demo purposes, we'll mark English as US English
		// A more sophisticated implementation might infer the region from the text
		langTag = "en-US"
		langName = "English (US)"
	case "es":
		// For demo purposes, we'll mark Spanish as Spanish from Spain
		langTag = "es-ES"
		langName = "Spanish (Spain)"
	case "pt":
		// For demo purposes, we'll mark Portuguese as Brazilian Portuguese
		langTag = "pt-BR"
		langName = "Portuguese (Brazil)"
	case "zh":
		// For demo purposes, we'll mark Chinese as Simplified Chinese
		langTag = "zh-CN"
		langName = "Chinese (Simplified)"
	}
	
	return langTag, langName, nil
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
	LOC                bool
	Line               bool
	Char               bool
	Word               bool
	DetectLanguage     bool
	ShowLanguageName   bool
	Paths              []string
	Input              io.Reader
	Output             io.Writer
	ErrorOutput        io.Writer
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
			fmt.Fprintf(cfg.ErrorOutput, "  -w, --words       Count words (default behavior)\n")
			fmt.Fprintf(cfg.ErrorOutput, "  -l, --lines       Count lines instead of words\n")
			fmt.Fprintf(cfg.ErrorOutput, "  -c, --chars       Count characters instead of words\n")
			fmt.Fprintf(cfg.ErrorOutput, "      --loc         Count lines of code in specified paths or current directory\n")
			fmt.Fprintf(cfg.ErrorOutput, "      --lang        Detect language of text in specified files or stdin\n")
			fmt.Fprintf(cfg.ErrorOutput, "      --lang-name   Show human-readable language name (implies --lang)\n")
			fmt.Fprintf(cfg.ErrorOutput, "  -h, --help        Show this help message\n")
			os.Exit(0)
		}
	}
	
	// Define flags
	var loc bool
	var l, c, w bool
	var lang, langName bool
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
		case "--lang":
			lang = true
			continue
		case "--lang-name":
			lang = true
			langName = true
			continue
		}
		
		// Handle non-flag arguments (paths for --loc or --lang)
		if (loc || lang) && !strings.HasPrefix(arg, "-") {
			paths = append(paths, arg)
			continue
		}
	}
	
	// Update the configuration
	cfg.LOC = loc
	cfg.Line = l
	cfg.Char = c
	cfg.DetectLanguage = lang
	cfg.ShowLanguageName = langName
	cfg.Word = w || (!cfg.Line && !cfg.Char && !cfg.LOC && !cfg.DetectLanguage)
	
	// Set paths
	if len(paths) > 0 {
		cfg.Paths = paths
	} else if loc || lang {
		// Default to current directory for --loc (consistent with existing behavior),
		// but don't default for language detection (will use stdin)
		if loc {
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
	
	// If we're detecting language, we need to handle the special case
	if cfg.DetectLanguage {
		// Check if paths are provided
		if len(cfg.Paths) > 0 {
			// Process each file
			for _, path := range cfg.Paths {
				if err := processFileForLanguage(path, cfg); err != nil {
					return err
				}
			}
			return nil
		}
		
		// No paths, process stdin
		return processReaderForLanguage(cfg.Input, cfg)
	}
	
	// Handle standard counting options
	// Check if paths are provided for standard counting
	if len(cfg.Paths) > 0 {
		// Process each file
		for _, path := range cfg.Paths {
			if err := processFileForCounting(path, cfg); err != nil {
				return err
			}
		}
		return nil
	}
	
	// No paths, process stdin for standard counting
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

// processFileForLanguage handles language detection for a specific file
func processFileForLanguage(path string, cfg *Config) error {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()
	
	// If multiple files, print the filename
	if len(cfg.Paths) > 1 {
		fmt.Fprintf(cfg.Output, "%s:\n", path)
	}
	
	// Process the file
	return processReaderForLanguage(file, cfg)
}

// processReaderForLanguage handles language detection for any io.Reader
func processReaderForLanguage(r io.Reader, cfg *Config) error {
	// Create a buffer to allow reading the input twice
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	
	// First pass: detect language
	langTag, langName, err := detectLanguage(tee)
	if err != nil {
		return fmt.Errorf("failed to detect language: %w", err)
	}
	
	// Second pass: handle standard counting options if requested
	var count int
	var needsCount bool
	switch {
	case cfg.Line:
		count = countLines(&buf)
		needsCount = true
	case cfg.Char:
		count = countChars(&buf)
		needsCount = true
	case cfg.Word:
		count = countWords(&buf)
		needsCount = true
	}
	
	// Print language info
	if cfg.ShowLanguageName {
		fmt.Fprintf(cfg.Output, "Language: %s\n", langName)
	} else {
		fmt.Fprintf(cfg.Output, "Language: %s\n", langTag)
	}
	
	// Print count if needed
	if needsCount {
		fmt.Fprintf(cfg.Output, "Count: %d\n", count)
	}
	
	return nil
}

// processFileForCounting handles standard counting operations for a specific file
func processFileForCounting(path string, cfg *Config) error {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()
	
	// Count based on the selected option
	var count int
	switch {
	case cfg.Line:
		count = countLines(file)
	case cfg.Char:
		count = countChars(file)
	case cfg.Word:
		count = countWords(file)
	}
	
	// If multiple files, print the filename with the count
	if len(cfg.Paths) > 1 {
		fmt.Fprintf(cfg.Output, "%s: %d\n", path, count)
	} else {
		fmt.Fprintln(cfg.Output, count)
	}
	
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
