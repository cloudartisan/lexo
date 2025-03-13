package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
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

func main() {
	// Define command line flags
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [--loc] [path...]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Count words from stdin or lines of code in specified paths.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	
	locFlag := flag.Bool("loc", false, "Count lines of code instead of words")
	flag.Parse()

	if *locFlag {
		args := flag.Args()
		paths := []string{"."}
		if len(args) > 0 {
			paths = args
		}
		
		if err := countLinesOfCode(paths); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println(countWords(os.Stdin))
	}
}
