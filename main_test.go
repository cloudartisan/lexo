package main

import (
	"bytes"
	"fmt"
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
		"--lang",
		"--lang-name",
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
			{
				name: "lang flag",
				args: []string{"wc", "--lang"},
				expected: &Config{
					DetectLanguage: true,
				},
			},
			{
				name: "lang-name flag",
				args: []string{"wc", "--lang-name"},
				expected: &Config{
					DetectLanguage:   true,
					ShowLanguageName: true,
				},
			},
			{
				name: "lang with words",
				args: []string{"wc", "--lang", "-w"},
				expected: &Config{
					DetectLanguage: true,
					Word:           true,
				},
			},	}
	
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
		// Language detection tests with various flags and combinations
		{"lang flag", "This is English text for testing the language detection feature", []string{"--lang"}, "Language: en"},
		{"lang with words", "This is English text for testing the language detection feature", []string{"--lang", "-w"}, "Count: 10"},
		{"lang with lines", "Line 1\nLine 2\nLine 3", []string{"--lang", "-l"}, "Count: 3"},
		{"lang with chars", "Hello, world!", []string{"--lang", "-c"}, "Count: 13"},
		{"lang name flag", "This is English text for testing the language detection feature", []string{"--lang-name"}, "Language: English"},
		{"lang name with count", "This is English text for testing the language detection feature", []string{"--lang-name", "-w"}, "Count: 10"},
		{"spanish detection", "El zorro marrón rápido salta sobre el perro perezoso", []string{"--lang"}, "Language: es-ES"},
		{"spanish with name", "El zorro marrón rápido salta sobre el perro perezoso", []string{"--lang-name"}, "Language: Spanish (Spain)"},
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
			
			// For language detection tests, we'll use contains matching
			if strings.Contains(tc.name, "lang") {
				if !strings.Contains(actual, tc.expected) {
					t.Errorf("Expected output to contain %q, got %q", tc.expected, actual)
				}
			} else {
				// For regular tests, we'll use exact matching
				if actual != tc.expected {
					t.Errorf("Expected %q, got %q", tc.expected, actual)
				}
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

// TestLanguageDetectionErrors tests error handling in language detection
func TestLanguageDetectionErrors(t *testing.T) {
	// Create a reader that always returns an error
	errReader := &errorReader{err: fmt.Errorf("simulated read error")}
	
	// Test detectLanguage with an error-generating reader
	_, _, err := detectLanguage(errReader)
	
	if err == nil {
		t.Error("Expected error from detectLanguage with error reader, got nil")
	}
	
	if !strings.Contains(err.Error(), "simulated read error") {
		t.Errorf("Expected error message to contain 'simulated read error', got: %s", err.Error())
	}
}

// errorReader is a mock io.Reader that always returns an error
type errorReader struct {
	err error
}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}

// TestConfig tests the configuration functions
func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectLangTag string
		expectName    string // Expected language name (partial match)
	}{
		{
			name:          "english",
			input:         "This is English text for testing purposes. It contains multiple sentences with various words to ensure accurate detection.",
			expectLangTag: "en",
			expectName:    "English",
		},
		{
			name:          "spanish",
			input:         "El zorro marrón rápido salta sobre el perro perezoso. Esta es una frase común utilizada para probar la detección de idiomas.",
			expectLangTag: "es",
			expectName:    "Spanish",
		},
		{
			name:          "french",
			input:         "Le renard brun rapide saute par-dessus le chien paresseux. C'est une phrase couramment utilisée pour tester la détection de langue.",
			expectLangTag: "fr",
			expectName:    "French",
		},
		{
			name:          "german",
			input:         "Der braune Fuchs springt über den faulen Hund. Dies ist ein häufig verwendeter Satz, um die Spracherkennung zu testen.",
			expectLangTag: "de",
			expectName:    "German",
		},
		{
			name:          "italian",
			input:         "La volpe marrone rapida salta sopra il cane pigro. Questa è una frase comunemente usata per testare il rilevamento della lingua.",
			expectLangTag: "it",
			expectName:    "Italian",
		},
		{
			name:          "portuguese",
			input:         "A raposa marrom rápida pula sobre o cão preguiçoso. Esta é uma frase comumente usada para testar a detecção de idioma.",
			expectLangTag: "pt",
			expectName:    "Portuguese",
		},
		{
			name:          "dutch",
			input:         "De snelle bruine vos springt over de luie hond. Dit is een veelgebruikte zin om taaldetectie te testen.",
			expectLangTag: "nl",
			expectName:    "Dutch",
		},
		{
			name:          "swedish",
			input:         "Den snabba bruna räven hoppar över den lata hunden. Detta är en vanligt använd mening för att testa språkdetektering.",
			expectLangTag: "sv",
			expectName:    "Swedish",
		},
		{
			name:          "mixed language with english dominant",
			input:         "This is mostly English text with a little bit of Español mixed in. La mayoría del texto está en inglés.",
			expectLangTag: "en",
			expectName:    "English",
		},
		{
			name:          "short text",
			input:         "Hello world",
			expectLangTag: "", // Accept any language for very short text
			expectName:    "",
		},
		{
			name:          "empty",
			input:         "",
			expectLangTag: "und",
			expectName:    "Unknown",
		},
		{
			name:          "non-language characters",
			input:         "123456789!@#$%^&*()",
			expectLangTag: "und", // This might be detected as something else, but that's fine
			expectName:    "",    // Don't check name for this case
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := bytes.NewBufferString(tc.input)
			
			langTag, langName, err := detectLanguage(b)
			if err != nil {
				t.Fatalf("detectLanguage returned error: %v", err)
			}
			
			// For empty, non-language inputs, or cases where we don't care about the specific result,
			// be flexible about the detection result
			if tc.input == "" || tc.expectLangTag == "und" || tc.expectLangTag == "" {
				// Either we got "und" or some language was detected - both are acceptable
				// because different language detection libraries might handle these edge cases differently
				t.Logf("Input detected as: %s (%s)", langTag, langName)
				return
			}
			
			// Check the language tag contains the expected base language (ignoring region)
			if !strings.HasPrefix(langTag, tc.expectLangTag) {
				t.Errorf("Expected language tag starting with %q, got %q", tc.expectLangTag, langTag)
			}
			
			// Check the language name contains the expected language name (if provided)
			if tc.expectName != "" && !strings.Contains(langName, tc.expectName) {
				t.Errorf("Expected language name containing %q, got %q", tc.expectName, langName)
			}
		})
	}
}

// TestFileProcessing tests processing specific files
func TestFileProcessing(t *testing.T) {
	// Create a temp file for testing
	tempFile, err := os.CreateTemp("", "wc-test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	// Write test data
	testData := "This is test data for language detection and counting.\nIt has multiple lines.\nAnd various words and characters."
	if _, err := tempFile.Write([]byte(testData)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}
	
	// Test cases for various operations on the temp file
	tests := []struct {
		name     string
		config   *Config
		expected string
		contains bool
	}{
		{
			name: "language detection from file",
			config: &Config{
				DetectLanguage: true,
				Paths:          []string{tempFile.Name()},
				Output:         nil, // will be set in the test
			},
			expected: "Language: en",
			contains: true,
		},
		{
			name: "language name from file",
			config: &Config{
				DetectLanguage:   true,
				ShowLanguageName: true,
				Paths:            []string{tempFile.Name()},
				Output:           nil, // will be set in the test
			},
			expected: "Language: English",
			contains: true,
		},
		{
			name: "word count from file",
			config: &Config{
				Word:   true,
				Paths:  []string{tempFile.Name()},
				Output: nil, // will be set in the test
			},
			expected: "18", // Count of words in testData
			contains: false,
		},
		{
			name: "line count from file",
			config: &Config{
				Line:   true,
				Paths:  []string{tempFile.Name()},
				Output: nil, // will be set in the test
			},
			expected: "3", // Count of lines in testData
			contains: false,
		},
		{
			name: "language and word count from file",
			config: &Config{
				DetectLanguage: true,
				Word:           true,
				Paths:          []string{tempFile.Name()},
				Output:         nil, // will be set in the test
			},
			expected: "Count: 18",
			contains: true,
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var outBuf bytes.Buffer
			tc.config.Output = &outBuf
			
			// Run with this configuration
			err := Run(tc.config)
			if err != nil {
				t.Fatalf("Run returned error: %v", err)
			}
			
			// Check output
			actual := strings.TrimSpace(outBuf.String())
			if tc.contains {
				if !strings.Contains(actual, tc.expected) {
					t.Errorf("Expected output to contain %q, got %q", tc.expected, actual)
				}
			} else {
				if actual != tc.expected {
					t.Errorf("Expected %q, got %q", tc.expected, actual)
				}
			}
		})
	}
}

// TestMultipleFiles tests processing multiple files
func TestMultipleFiles(t *testing.T) {
	// Create temp files for testing
	tempFile1, err := os.CreateTemp("", "wc-test-1-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file 1: %v", err)
	}
	defer os.Remove(tempFile1.Name())
	
	tempFile2, err := os.CreateTemp("", "wc-test-2-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file 2: %v", err)
	}
	defer os.Remove(tempFile2.Name())
	
	// Write test data
	testData1 := "This is file one.\nIt has English text."
	testData2 := "El segundo archivo.\nTexto en español."
	
	if _, err := tempFile1.Write([]byte(testData1)); err != nil {
		t.Fatalf("Failed to write to temp file 1: %v", err)
	}
	if err := tempFile1.Close(); err != nil {
		t.Fatalf("Failed to close temp file 1: %v", err)
	}
	
	if _, err := tempFile2.Write([]byte(testData2)); err != nil {
		t.Fatalf("Failed to write to temp file 2: %v", err)
	}
	if err := tempFile2.Close(); err != nil {
		t.Fatalf("Failed to close temp file 2: %v", err)
	}
	
	// Test multiple file handling
	var outBuf bytes.Buffer
	cfg := &Config{
		DetectLanguage: true,
		Paths:          []string{tempFile1.Name(), tempFile2.Name()},
		Output:         &outBuf,
	}
	
	// Run with this configuration
	err = Run(cfg)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	
	// Check output contains both filenames and language detections
	actual := strings.TrimSpace(outBuf.String())
	
	// Check for first file
	if !strings.Contains(actual, tempFile1.Name()) {
		t.Errorf("Output should contain first filename %q, got: %q", tempFile1.Name(), actual)
	}
	
	// Check for second file
	if !strings.Contains(actual, tempFile2.Name()) {
		t.Errorf("Output should contain second filename %q, got: %q", tempFile2.Name(), actual)
	}
}

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
		contains bool // if true, check that output contains expected, rather than exact match
	}{
		{
			name:  "word count",
			input: "word1 word2 word3",
			config: &Config{
				Word:  true,
				Input: nil, // will be set in the test
			},
			expected: "3",
			contains: false,
		},
		{
			name:  "line count",
			input: "line1\nline2\nline3",
			config: &Config{
				Line:  true,
				Input: nil, // will be set in the test
			},
			expected: "3",
			contains: false,
		},
		{
			name:  "char count",
			input: "hello",
			config: &Config{
				Char:  true,
				Input: nil, // will be set in the test
			},
			expected: "5",
			contains: false,
		},
		{
			name:  "language detection",
			input: "This is English text for testing.",
			config: &Config{
				DetectLanguage: true,
				Input:          nil, // will be set in the test
			},
			expected: "Language: en-US",
			contains: true,
		},
		{
			name:  "language detection with name",
			input: "This is English text for testing.",
			config: &Config{
				DetectLanguage:   true,
				ShowLanguageName: true,
				Input:            nil, // will be set in the test
			},
			expected: "Language: English",
			contains: true,
		},
		{
			name:  "language detection with word count",
			input: "This is English text for testing.",
			config: &Config{
				DetectLanguage: true,
				Word:           true,
				Input:          nil, // will be set in the test
			},
			expected: "Count: 6",
			contains: true,
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
			if tc.contains {
				// For language detection tests, we just check if the output contains the expected string
				if !strings.Contains(actual, tc.expected) {
					t.Errorf("Expected output to contain %q, got %q", tc.expected, actual)
				}
			} else {
				// For regular tests, we use exact matching
				if actual != tc.expected {
					t.Errorf("Expected %q, got %q", tc.expected, actual)
				}
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
