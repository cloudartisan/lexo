package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/abadojack/whatlanggo"
)

func countWords(r io.Reader) int {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)

	wc := 0
	for scanner.Scan() {
		wc++
	}

	return wc
}

// WordFrequency represents a word and its frequency count
type WordFrequency struct {
	Word  string
	Count int
}

// analyzeWordFrequency counts the frequency of each word in the text
// and returns the results sorted by frequency (highest first) or alphabetically
func analyzeWordFrequency(r io.Reader, sortByCount bool, limit int) ([]WordFrequency, error) {
	// If limit is 0 or negative, set a reasonable default
	if limit <= 0 {
		limit = 10
	}

	// Create a scanner to read words
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)

	// Use a map to count word frequencies
	wordCounts := make(map[string]int)

	// Process each word
	for scanner.Scan() {
		word := scanner.Text()
		
		// Convert to lowercase for case-insensitive counting
		word = strings.ToLower(word)
		
		// Remove any punctuation at the start or end of the word
		word = strings.Trim(word, ".,;:!?\"'()[]{}")
		
		// Skip empty strings after trimming
		if word == "" {
			continue
		}
		
		// Increment the word count
		wordCounts[word]++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Convert map to slice for sorting
	var frequencies []WordFrequency
	for word, count := range wordCounts {
		frequencies = append(frequencies, WordFrequency{Word: word, Count: count})
	}

	// Sort the frequencies
	if sortByCount {
		// Sort by count (descending) with alphabetical tiebreaker
		sort.Slice(frequencies, func(i, j int) bool {
			if frequencies[i].Count == frequencies[j].Count {
				return frequencies[i].Word < frequencies[j].Word
			}
			return frequencies[i].Count > frequencies[j].Count
		})
	} else {
		// Sort alphabetically
		sort.Slice(frequencies, func(i, j int) bool {
			return frequencies[i].Word < frequencies[j].Word
		})
	}

	// Apply limit
	if limit > 0 && limit < len(frequencies) {
		frequencies = frequencies[:limit]
	}

	return frequencies, nil
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

// CodeStats holds statistics about code in a file or directory
type CodeStats struct {
	Total     int // Total lines
	Code      int // Lines of code (non-blank, non-comment)
	Comments  int // Comment lines
	Blank     int // Blank lines
	Files     int // Number of files processed
}

// countLinesOfCode counts lines of code in files or directories without external dependencies
func countLinesOfCode(paths []string) error {
	// Set of directories to skip
	skipDirs := map[string]bool{
		".git":         true,
		".hg":          true,
		".svn":         true,
		"node_modules": true,
		".idea":        true,
		".vscode":      true,
		"target":       true,
		"build":        true,
		"dist":         true,
		"bin":          true,
		"obj":          true,
	}

	// Set of file extensions to consider as code
	codeExtensions := map[string]bool{
		".go":    true,
		".java":  true,
		".js":    true,
		".ts":    true,
		".jsx":   true,
		".tsx":   true,
		".py":    true,
		".c":     true,
		".cpp":   true,
		".h":     true,
		".hpp":   true,
		".cs":    true,
		".rb":    true,
		".php":   true,
		".scala": true,
		".rs":    true,
		".swift": true,
		".sh":    true,
		".bat":   true,
		".ps1":   true,
		".html":  true,
		".css":   true,
		".scss":  true,
		".sql":   true,
		".kt":    true,
		".kts":   true,
		".ex":    true,
		".exs":   true,
		".md":    true,
	}

	// Initialize statistics
	stats := CodeStats{}

	// If no paths provided, use current directory
	if len(paths) == 0 {
		paths = []string{"."}
	}

	// Process each path
	for _, path := range paths {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to get file info for %s: %w", path, err)
		}

		if fileInfo.IsDir() {
			// Process directory recursively
			err = processDirectory(path, skipDirs, codeExtensions, &stats)
			if err != nil {
				return err
			}
		} else {
			// Process single file
			fileStats, err := processFile(path)
			if err != nil {
				return err
			}
			
			// Only count it if it has a recognized extension
			ext := strings.ToLower(path[strings.LastIndexByte(path, '.')+1:])
			if _, ok := codeExtensions["."+ext]; ok || len(ext) == 0 || ext == path {
				stats.Total += fileStats.Total
				stats.Code += fileStats.Code
				stats.Comments += fileStats.Comments
				stats.Blank += fileStats.Blank
				stats.Files++
			}
		}
	}

	// Print the code count
	fmt.Println(stats.Code)
	
	return nil
}

// processDirectory processes a directory recursively
func processDirectory(dirPath string, skipDirs map[string]bool, codeExtensions map[string]bool, stats *CodeStats) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	for _, entry := range entries {
		entryName := entry.Name()
		entryPath := dirPath + "/" + entryName

		// Skip hidden files and directories
		if strings.HasPrefix(entryName, ".") {
			continue
		}

		if entry.IsDir() {
			// Skip directories in the ignore list
			if skipDirs[entryName] {
				continue
			}

			// Process subdirectory recursively
			err = processDirectory(entryPath, skipDirs, codeExtensions, stats)
			if err != nil {
				return err
			}
		} else {
			// Check if it's a code file based on extension
			ext := strings.ToLower(entryName[strings.LastIndexByte(entryName, '.')+1:])
			if _, ok := codeExtensions["."+ext]; !ok {
				continue
			}

			// Process code file
			fileStats, err := processFile(entryPath)
			if err != nil {
				// Just skip problematic files
				continue
			}

			stats.Total += fileStats.Total
			stats.Code += fileStats.Code
			stats.Comments += fileStats.Comments
			stats.Blank += fileStats.Blank
			stats.Files++
		}
	}

	return nil
}

// processFile counts lines of code, comments, and blank lines in a single file
func processFile(filePath string) (CodeStats, error) {
	stats := CodeStats{}

	file, err := os.Open(filePath)
	if err != nil {
		return stats, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	isMultilineComment := false
	
	// Get file extension to determine comment syntax
	ext := strings.ToLower(filePath[strings.LastIndexByte(filePath, '.')+1:])
	
	// This is a simplified approach - in a full implementation, you'd want
	// a more robust language detection mechanism
	for scanner.Scan() {
		line := scanner.Text()
		stats.Total++
		
		// Trimmed line for blank line detection
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			stats.Blank++
			continue
		}
		
		// Detect comments based on file extension
		// This is a simplified approach - a real implementation would be more thorough
		switch ext {
		case "go", "c", "cpp", "java", "js", "ts", "cs", "swift", "kt":
			// Handle C-style comments
			if isMultilineComment {
				stats.Comments++
				if strings.Contains(line, "*/") {
					isMultilineComment = false
				}
				continue
			}
			
			if strings.HasPrefix(trimmedLine, "//") {
				stats.Comments++
				continue
			}
			
			if strings.HasPrefix(trimmedLine, "/*") {
				isMultilineComment = true
				stats.Comments++
				if strings.Contains(line, "*/") {
					isMultilineComment = false
				}
				continue
			}
			
		case "py", "rb":
			// Handle Python/Ruby style comments
			if strings.HasPrefix(trimmedLine, "#") {
				stats.Comments++
				continue
			}
			
		case "sh", "bash":
			// Handle shell script comments
			if strings.HasPrefix(trimmedLine, "#") {
				stats.Comments++
				continue
			}
			
		// Add more languages as needed
		}
		
		// If not a comment or blank line, count as code
		stats.Code++
	}

	if err := scanner.Err(); err != nil {
		return stats, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	return stats, nil
}

// Config holds the configuration for the program
type Config struct {
	LOC                bool
	Line               bool
	Char               bool
	Word               bool
	DetectLanguage     bool
	ShowLanguageName   bool
	FrequencyAnalysis  bool
	FrequencyLimit     int
	SortByCount        bool
	Paths              []string
	Input              io.Reader
	Output             io.Writer
	ErrorOutput        io.Writer
}

// NewDefaultConfig creates a default configuration
func NewDefaultConfig() *Config {
	return &Config{
		Input:          os.Stdin,
		Output:         os.Stdout,
		ErrorOutput:    os.Stderr,
		FrequencyLimit: 10, // Default to showing top 10 words
	}
}

// ParseFlags parses command-line flags and updates the configuration
func ParseFlags(cfg *Config) {
	// Check for help flag manually
	for _, arg := range os.Args[1:] {
		if arg == "-h" || arg == "--help" {
			fmt.Fprintf(cfg.ErrorOutput, "Usage: %s [flags] [path...]\n\n", os.Args[0])
			fmt.Fprintf(cfg.ErrorOutput, "Text and code analysis utility for counting, language detection, and more.\n")
			fmt.Fprintf(cfg.ErrorOutput, "By default, counts words from stdin.\n\n")
			fmt.Fprintf(cfg.ErrorOutput, "Options:\n")
			fmt.Fprintf(cfg.ErrorOutput, "  -w, --words       Count words (default behavior)\n")
			fmt.Fprintf(cfg.ErrorOutput, "  -l, --lines       Count lines instead of words\n")
			fmt.Fprintf(cfg.ErrorOutput, "  -c, --chars       Count characters instead of words\n")
			fmt.Fprintf(cfg.ErrorOutput, "      --loc         Count lines of code in specified paths or current directory\n")
			fmt.Fprintf(cfg.ErrorOutput, "      --lang        Detect language of text in specified files or stdin\n")
			fmt.Fprintf(cfg.ErrorOutput, "      --lang-name   Show human-readable language name (implies --lang)\n")
			fmt.Fprintf(cfg.ErrorOutput, "      --freq        Analyze word frequency\n")
			fmt.Fprintf(cfg.ErrorOutput, "      --sort-count  Sort frequency by count (default is alphabetical)\n")
			fmt.Fprintf(cfg.ErrorOutput, "      --limit N     Limit frequency results to top N words\n")
			fmt.Fprintf(cfg.ErrorOutput, "  -h, --help        Show this help message\n")
			os.Exit(0)
		}
	}
	
	// Define flags
	var loc bool
	var l, c, w bool
	var lang, langName bool
	var freq, sortByCount bool
	var limit int
	var paths []string
	
	// Process args to handle GNU-style long options
	for i := 0; i < len(os.Args[1:]); i++ {
		arg := os.Args[1:][i]
		
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
		case "--freq":
			freq = true
			continue
		case "--sort-count":
			sortByCount = true
			continue
		case "--limit":
			// Check if there's a next argument for the limit value
			if i+1 < len(os.Args[1:]) {
				// Try to parse the next argument as a number
				if n, err := fmt.Sscanf(os.Args[1:][i+1], "%d", &limit); n == 1 && err == nil {
					// Skip the next arg since we've consumed it
					i++
					continue
				}
			}
			// If we can't parse a number, use the default limit
			continue
		}
		
		// Handle non-flag arguments (paths for all operations)
		if !strings.HasPrefix(arg, "-") {
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
	cfg.FrequencyAnalysis = freq
	cfg.SortByCount = sortByCount
	if limit > 0 {
		cfg.FrequencyLimit = limit
	}
	
	// Set default behavior to match wc: if no counting flags are specified, show lines, words, and chars
	if !w && !l && !c && !loc && !lang && !freq {
		cfg.Line = true
		cfg.Word = true 
		cfg.Char = true
	} else {
		cfg.Word = w
	}
	
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
	
	// If we're doing frequency analysis, handle that
	if cfg.FrequencyAnalysis {
		// Check if paths are provided
		if len(cfg.Paths) > 0 {
			// Process each file
			for _, path := range cfg.Paths {
				if err := processFileForFrequency(path, cfg); err != nil {
					return err
				}
			}
			return nil
		}
		
		// No paths, process stdin
		return processReaderForFrequency(cfg.Input, cfg)
	}
	
	// Handle standard counting options
	// Check if paths are provided for standard counting
	if len(cfg.Paths) > 0 {
		// Process each file
		totalLines, totalWords, totalChars := 0, 0, 0
		showTotal := len(cfg.Paths) > 1 && cfg.Line && cfg.Word && cfg.Char
		
		for _, path := range cfg.Paths {
			lines, words, chars, err := processFileForCounting(path, cfg)
			if err != nil {
				return err
			}
			
			// If we're doing a wc-like output with multiple files, we need to track totals
			if showTotal {
				totalLines += lines
				totalWords += words
				totalChars += chars
			}
		}
		
		// Display totals for multiple files
		if showTotal {
			FormatLikeWC(cfg.Output, totalLines, totalWords, totalChars, "total")
		}
		
		return nil
	}
	
	// No paths, process stdin for standard counting
	// Read all input into a buffer to allow multiple passes
	inputData, err := io.ReadAll(cfg.Input)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	
	// If default behavior (like wc), show all three counts
	if cfg.Line && cfg.Word && cfg.Char {
		lineCount := countLines(bytes.NewReader(inputData))
		wordCount := countWords(bytes.NewReader(inputData))
		charCount := countChars(bytes.NewReader(inputData))
		
		// Format output like wc: lines words chars
		FormatLikeWC(cfg.Output, lineCount, wordCount, charCount, "")
		return nil
	}
	
	// Otherwise handle individual flags
	var count int
	switch {
	case cfg.Line:
		count = countLines(bytes.NewReader(inputData))
	case cfg.Char:
		count = countChars(bytes.NewReader(inputData))
	case cfg.Word:
		count = countWords(bytes.NewReader(inputData))
	}
	
	// Match wc's spacing for output without a filename (no trailing space)
	fmt.Fprintf(cfg.Output, "%8d", count)
	fmt.Fprintln(cfg.Output)
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

// FormatLikeWC formats counts exactly like the wc utility
func FormatLikeWC(w io.Writer, lineCount, wordCount, charCount int, path string) {
	// Exact format string to match wc output
	// The key is to use the spacing for consistent results
	if path == "" {
		// No extra space at the end for stdin
		fmt.Fprintf(w, "%8d %7d %7d", lineCount, wordCount, charCount)
	} else {
		// With path
		fmt.Fprintf(w, "%8d %7d %7d %s", lineCount, wordCount, charCount, path)
	}
	// Use Fprintln to add the newline exactly like wc does
	fmt.Fprintln(w)
}

// processFileForCounting handles standard counting operations for a specific file
// returns lineCount, wordCount, charCount, and error
func processFileForCounting(path string, cfg *Config) (int, int, int, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()
	
	// Read the file contents to handle multiple passes
	fileContents, err := io.ReadAll(file)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to read file %s: %w", path, err) 
	}
	
	// Set up various counts
	var lineCount, wordCount, charCount int
	
	// If default behavior (like wc), show all three counts
	if cfg.Line && cfg.Word && cfg.Char {
		lineCount = countLines(bytes.NewReader(fileContents))
		wordCount = countWords(bytes.NewReader(fileContents))
		charCount = countChars(bytes.NewReader(fileContents))
		
		// Use our wc-like formatter
		FormatLikeWC(cfg.Output, lineCount, wordCount, charCount, path)
		return lineCount, wordCount, charCount, nil
	}
	
	// Otherwise handle individual flags
	var count int
	switch {
	case cfg.Line:
		count = countLines(bytes.NewReader(fileContents))
		lineCount = count
	case cfg.Char:
		count = countChars(bytes.NewReader(fileContents))
		charCount = count
	case cfg.Word:
		count = countWords(bytes.NewReader(fileContents))
		wordCount = count
	}
	
	// Print with filename, using the same spacing as wc
	fmt.Fprintf(cfg.Output, "%8d %s\n", count, path)
	
	return lineCount, wordCount, charCount, nil
}

// processFileForFrequency handles word frequency analysis for a specific file
func processFileForFrequency(path string, cfg *Config) error {
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
	return processReaderForFrequency(file, cfg)
}

// processReaderForFrequency handles word frequency analysis for any io.Reader
func processReaderForFrequency(r io.Reader, cfg *Config) error {
	// Analyze word frequency
	frequencies, err := analyzeWordFrequency(r, cfg.SortByCount, cfg.FrequencyLimit)
	if err != nil {
		return fmt.Errorf("failed to analyze word frequency: %w", err)
	}
	
	// Determine the longest word to format output nicely
	maxWordLen := 0
	for _, wf := range frequencies {
		if len(wf.Word) > maxWordLen {
			maxWordLen = len(wf.Word)
		}
	}
	
	// Print header
	if cfg.SortByCount {
		fmt.Fprintf(cfg.Output, "Word frequency (sorted by count):\n")
	} else {
		fmt.Fprintf(cfg.Output, "Word frequency (sorted alphabetically):\n")
	}
	
	// Print a separator line
	fmt.Fprintf(cfg.Output, "%s  %s\n", strings.Repeat("-", maxWordLen), "------")
	
	// Print the results in a nicely formatted two-column layout
	for _, wf := range frequencies {
		fmt.Fprintf(cfg.Output, "%-*s  %6d\n", maxWordLen, wf.Word, wf.Count)
	}
	
	return nil
}

// Allow os.Exit to be mocked in tests
var osExit = os.Exit

func main() {
	// Create default configuration
	cfg := NewDefaultConfig()
	
	// Parse command-line flags
	ParseFlags(cfg)
	
	// Run the program
	if err := Run(cfg); err != nil {
		fmt.Fprintf(cfg.ErrorOutput, "Error: %v\n", err)
		osExit(1)
	}
}
