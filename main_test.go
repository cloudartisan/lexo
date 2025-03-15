package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3 word4\n")

	expected := 4
	actual := countWords(b)

	if actual != expected {
		t.Errorf("Expected %d, got %d.\n", expected, actual)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("line1\nline2\nline3\nline4\n")

	expected := 4
	actual := countLines(b)

	if actual != expected {
		t.Errorf("Expected %d, got %d.\n", expected, actual)
	}
}

func TestCountLinesWithoutTrailingNewline(t *testing.T) {
	b := bytes.NewBufferString("line1\nline2\nline3")

	expected := 3
	actual := countLines(b)

	if actual != expected {
		t.Errorf("Expected %d, got %d.\n", expected, actual)
	}
}

func TestCountChars(t *testing.T) {
	b := bytes.NewBufferString("hello")

	expected := 5
	actual := countChars(b)

	if actual != expected {
		t.Errorf("Expected %d, got %d.\n", expected, actual)
	}
}

func TestEmptyInput(t *testing.T) {
	b := bytes.NewBufferString("")

	// All counts should be 0 for empty input
	if countWords(b) != 0 {
		t.Errorf("Expected 0 words for empty input, got %d", countWords(b))
	}

	b.Reset()
	if countLines(b) != 0 {
		t.Errorf("Expected 0 lines for empty input, got %d", countLines(b))
	}

	b.Reset()
	if countChars(b) != 0 {
		t.Errorf("Expected 0 chars for empty input, got %d", countChars(b))
	}
}

// TestLOCFeatureExists checks that the --loc flag is properly registered
func TestLOCFeatureExists(t *testing.T) {
	if os.Getenv("SKIP_LOC_TEST") != "" {
		t.Skip("Skipping LOC test")
	}

	// We only check that the flag is recognized, not that it produces the correct output
	// since the output depends on the environment
	cmd := exec.Command("./wc", "--help")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("Failed to run help command: %v", err)
	}
	
	if !bytes.Contains(output, []byte("--loc")) {
		t.Errorf("Expected --loc flag to be listed in help output")
	}
	
	// Check for all required flags
	requiredFlags := []string{
		"-w", "--words",
		"-l", "--lines",
		"-c", "--chars",
		"--loc",
	}
	
	for _, flag := range requiredFlags {
		if !bytes.Contains(output, []byte(flag)) {
			t.Errorf("Expected %s flag to be listed in help output", flag)
		}
	}
}


// TestParseFlags tests the flag parsing function
func TestParseFlags(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()
	
	tests := []struct {
		name     string
		args     []string
		expected *Config
	}{
		{
			name: "default",
			args: []string{"wc"},
			expected: &Config{
				Word: true,
			},
		},
		{
			name: "lines short flag",
			args: []string{"wc", "-l"},
			expected: &Config{
				Line: true,
			},
		},
		{
			name: "lines long flag",
			args: []string{"wc", "--lines"},
			expected: &Config{
				Line: true,
			},
		},
		{
			name: "chars short flag",
			args: []string{"wc", "-c"},
			expected: &Config{
				Char: true,
			},
		},
		{
			name: "chars long flag",
			args: []string{"wc", "--chars"},
			expected: &Config{
				Char: true,
			},
		},
		{
			name: "words short flag",
			args: []string{"wc", "-w"},
			expected: &Config{
				Word: true,
			},
		},
		{
			name: "words long flag",
			args: []string{"wc", "--words"},
			expected: &Config{
				Word: true,
			},
		},
		{
			name: "loc flag",
			args: []string{"wc", "--loc"},
			expected: &Config{
				LOC:   true,
				Paths: []string{"."},
			},
		},
		{
			name: "loc flag with paths",
			args: []string{"wc", "--loc", "dir1", "dir2"},
			expected: &Config{
				LOC:   true,
				Paths: []string{"dir1", "dir2"},
			},
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set args for this test
			os.Args = tc.args
			
			// Create config and parse flags
			var errBuf bytes.Buffer
			cfg := &Config{
				ErrorOutput: &errBuf,
			}
			
			// Skip help check (which would exit in real usage)
			if len(tc.args) > 1 && (tc.args[1] == "-h" || tc.args[1] == "--help") {
				cfg.Word = true
			} else {
				ParseFlags(cfg)
			}
			
			// Check flags
			if cfg.Word != tc.expected.Word {
				t.Errorf("Word flag: expected %v, got %v", tc.expected.Word, cfg.Word)
			}
			if cfg.Line != tc.expected.Line {
				t.Errorf("Line flag: expected %v, got %v", tc.expected.Line, cfg.Line)
			}
			if cfg.Char != tc.expected.Char {
				t.Errorf("Char flag: expected %v, got %v", tc.expected.Char, cfg.Char)
			}
			if cfg.LOC != tc.expected.LOC {
				t.Errorf("LOC flag: expected %v, got %v", tc.expected.LOC, cfg.LOC)
			}
			
			// Check paths if LOC is true
			if cfg.LOC {
				if len(cfg.Paths) != len(tc.expected.Paths) {
					t.Errorf("Paths: expected %v, got %v", tc.expected.Paths, cfg.Paths)
				} else {
					for i, p := range cfg.Paths {
						if p != tc.expected.Paths[i] {
							t.Errorf("Path %d: expected %s, got %s", i, tc.expected.Paths[i], p)
						}
					}
				}
			}
		})
	}
}

// TestCommandLineFlags tests the operation of all command-line flags
func TestCommandLineFlags(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TEST") != "" {
		t.Skip("Skipping integration test")
	}

	// Build the executable
	if err := exec.Command("go", "build", "-o", "wc.test").Run(); err != nil {
		t.Fatalf("Failed to build test executable: %v", err)
	}
	defer os.Remove("wc.test")
	
	// Test cases for different flags
	tests := []struct {
		name     string
		input    string
		flags    []string
		expected string
	}{
		{"words default", "one two three", []string{}, "3"},
		{"words short", "one two three", []string{"-w"}, "3"},
		{"words long", "one two three", []string{"--words"}, "3"},
		{"lines short", "line1\nline2\nline3", []string{"-l"}, "3"},
		{"lines long", "line1\nline2\nline3", []string{"--lines"}, "3"},
		{"chars short", "hello", []string{"-c"}, "5"},
		{"chars long", "hello", []string{"--chars"}, "5"},
		{"empty input", "", []string{}, "0"},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("./wc.test", tc.flags...)
			stdin, err := cmd.StdinPipe()
			if err != nil {
				t.Fatalf("Failed to get stdin pipe: %v", err)
			}
			
			go func() {
				defer stdin.Close()
				io.WriteString(stdin, tc.input)
			}()
			
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Command failed: %v", err)
			}
			
			actual := strings.TrimSpace(string(output))
			if actual != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

// TestErrorHandling checks if the program properly handles errors
func TestErrorHandling(t *testing.T) {
	// Test case for when scc is not installed
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", "/nonexistent")
	defer os.Setenv("PATH", oldPath)
	
	// Create a test configuration
	var errBuf bytes.Buffer
	cfg := &Config{
		LOC:         true,
		Paths:       []string{"."},
		ErrorOutput: &errBuf,
	}
	
	// Run the program (should fail)
	err := Run(cfg)
	
	if err == nil {
		t.Errorf("Expected error when scc is not installed")
	}
	
	if !strings.Contains(err.Error(), "scc is not installed") {
		t.Errorf("Expected error message to mention scc not being installed, got: %s", err.Error())
	}
}

// TestConfig tests the configuration functions
func TestConfig(t *testing.T) {
	// Test NewDefaultConfig
	cfg := NewDefaultConfig()
	if cfg.Input != os.Stdin || cfg.Output != os.Stdout || cfg.ErrorOutput != os.Stderr {
		t.Errorf("NewDefaultConfig didn't set proper defaults")
	}
	
	// Test Run with different configurations
	tests := []struct {
		name     string
		input    string
		config   *Config
		expected string
	}{
		{
			name:  "word count",
			input: "word1 word2 word3",
			config: &Config{
				Word:  true,
				Input: nil, // will be set in the test
			},
			expected: "3",
		},
		{
			name:  "line count",
			input: "line1\nline2\nline3",
			config: &Config{
				Line:  true,
				Input: nil, // will be set in the test
			},
			expected: "3",
		},
		{
			name:  "char count",
			input: "hello",
			config: &Config{
				Char:  true,
				Input: nil, // will be set in the test
			},
			expected: "5",
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up input and output buffers
			inBuf := bytes.NewBufferString(tc.input)
			var outBuf bytes.Buffer
			
			tc.config.Input = inBuf
			tc.config.Output = &outBuf
			
			// Run with this configuration
			err := Run(tc.config)
			if err != nil {
				t.Fatalf("Run returned error: %v", err)
			}
			
			// Check output
			actual := strings.TrimSpace(outBuf.String())
			if actual != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

// TestMockLOC tests the countLinesOfCode function with a mock scc command
func TestMockLOC(t *testing.T) {
	// Skip if we can't create test files
	tempDir, err := os.MkdirTemp("", "wc-test")
	if err != nil {
		t.Skip("Couldn't create temp directory for test")
	}
	defer os.RemoveAll(tempDir)
	
	// Create a fake scc output script
	fakeScc := filepath.Join(tempDir, "scc")
	sccScript := `#!/bin/sh
echo '[{"Name":"Go","Code":100,"Comment":20,"Blank":10,"Complexity":5,"Count":3,"WeightedComplex":15}]'
`
	
	if err := os.WriteFile(fakeScc, []byte(sccScript), 0755); err != nil {
		t.Skipf("Failed to write fake scc script: %v", err)
	}
	
	// Temporarily add our fake scc to the PATH
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", tempDir+string(os.PathListSeparator)+oldPath)
	defer os.Setenv("PATH", oldPath)
	
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	// Run countLinesOfCode
	err = countLinesOfCode([]string{"."})
	
	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	
	if err != nil {
		t.Fatalf("countLinesOfCode failed: %v", err)
	}
	
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := strings.TrimSpace(buf.String())
	
	// Check that it printed the expected count
	if output != "100" {
		t.Errorf("Expected output '100', got '%s'", output)
	}
}
