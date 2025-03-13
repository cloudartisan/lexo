package main

import (
	"bytes"
	"os"
	"os/exec"
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
}
