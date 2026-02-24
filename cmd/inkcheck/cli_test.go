package main

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// buildCLI builds the inkcheck binary for testing.
func buildCLI(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "inkcheck")

	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build CLI: %v\n%s", err, output)
	}
	return binaryPath
}

// TestCLI_TextOutput tests the CLI with text output format.
func TestCLI_TextOutput(t *testing.T) {
	binary := buildCLI(t)
	testFile := filepath.Join("..", "..", "testdata", "argumentative_essay.md")

	cmd := exec.Command(binary, "-m", "sentence_length_variance", testFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI failed: %v\n%s", err, output)
	}

	// Output should contain the metric name
	if !strings.Contains(string(output), "sentence_length_variance") {
		t.Errorf("expected output to contain metric name, got: %s", output)
	}

	// Output should contain a number
	if !strings.Contains(string(output), ".") && !strings.Contains(string(output), "0") {
		t.Errorf("expected output to contain a numeric value, got: %s", output)
	}
}

// TestCLI_JSONOutput tests the CLI with JSON output format.
func TestCLI_JSONOutput(t *testing.T) {
	binary := buildCLI(t)
	testFile := filepath.Join("..", "..", "testdata", "technical_writing.md")

	cmd := exec.Command(binary, "-m", "paragraph_variance", "-format", "json", testFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI failed: %v\n%s", err, output)
	}

	// Parse JSON to verify it's valid
	var metrics []map[string]interface{}
	if err := json.Unmarshal(output, &metrics); err != nil {
		t.Fatalf("failed to parse JSON output: %v\noutput: %s", err, output)
	}

	// Should have at least one metric
	if len(metrics) == 0 {
		t.Error("expected at least one metric in JSON output")
	}

	// Verify metric structure
	metric := metrics[0]
	if _, ok := metric["name"]; !ok {
		t.Error("metric should have 'name' field")
	}
	if _, ok := metric["value"]; !ok {
		t.Error("metric should have 'value' field")
	}
}

// TestCLI_StdinInput tests reading from stdin.
func TestCLI_StdinInput(t *testing.T) {
	binary := buildCLI(t)
	testText := "This is a test sentence. This is another test sentence."

	cmd := exec.Command(binary, "-m", "sentence_length_variance")
	cmd.Stdin = strings.NewReader(testText)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI failed: %v\n%s", err, output)
	}

	if !strings.Contains(string(output), "sentence_length_variance") {
		t.Errorf("expected output to contain metric name, got: %s", output)
	}
}

// TestCLI_AllMetrics tests running multiple specific metrics.
// Note: Skips semantic metrics to avoid model download requirement.
func TestCLI_AllMetrics(t *testing.T) {
	binary := buildCLI(t)
	testFile := filepath.Join("..", "..", "testdata", "narrative_prose.md")

	// Test multiple metrics individually rather than all at once
	// to avoid semantic model requirement
	metrics := []string{
		"paragraph_variance",
		"sentence_length_variance",
		"hedging_analysis",
		"readability",
	}

	for _, metric := range metrics {
		cmd := exec.Command(binary, "-m", metric, testFile)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("CLI failed for %s: %v\n%s", metric, err, output)
		}

		if !strings.Contains(string(output), metric) {
			t.Errorf("expected output to contain %s", metric)
		}
	}
}

// TestCLI_MultipleFiles tests processing multiple files.
func TestCLI_MultipleFiles(t *testing.T) {
	binary := buildCLI(t)
	file1 := filepath.Join("..", "..", "testdata", "argumentative_essay.md")
	file2 := filepath.Join("..", "..", "testdata", "narrative_prose.md")

	cmd := exec.Command(binary, "-m", "paragraph_variance", file1, file2)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI failed: %v\n%s", err, output)
	}

	// Should contain both file paths in output
	outputStr := string(output)
	if !strings.Contains(outputStr, "argumentative_essay.md") {
		t.Error("expected output to contain first filename")
	}
	if !strings.Contains(outputStr, "narrative_prose.md") {
		t.Error("expected output to contain second filename")
	}
}

// TestCLI_InvalidMetric tests error handling for invalid metric names.
func TestCLI_InvalidMetric(t *testing.T) {
	binary := buildCLI(t)
	testFile := filepath.Join("..", "..", "testdata", "edge_cases.md")

	cmd := exec.Command(binary, "-m", "nonexistent_metric", testFile)
	output, err := cmd.CombinedOutput()

	// Should exit with non-zero status
	if err == nil {
		t.Error("expected non-zero exit code for invalid metric")
	}

	// Should contain error message
	if !strings.Contains(string(output), "unknown metric") {
		t.Errorf("expected error message about unknown metric, got: %s", output)
	}
}

// TestCLI_JSONMultipleMetrics tests JSON output with multiple metrics.
// Tests specific metrics to avoid semantic model requirement.
func TestCLI_JSONMultipleMetrics(t *testing.T) {
	binary := buildCLI(t)
	testFile := filepath.Join("..", "..", "testdata", "technical_writing.md")

	// Test a few different metrics with JSON output
	testMetrics := []string{
		"paragraph_variance",
		"hedging_analysis",
		"vocabulary_sophistication",
	}

	for _, metric := range testMetrics {
		cmd := exec.Command(binary, "-m", metric, "-format", "json", testFile)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("CLI failed for %s: %v\n%s", metric, err, output)
		}

		// Parse JSON
		var metrics []map[string]interface{}
		if err := json.Unmarshal(output, &metrics); err != nil {
			t.Fatalf("failed to parse JSON output for %s: %v\noutput: %s", metric, err, output)
		}

		// Should have the metric
		if len(metrics) != 1 {
			t.Errorf("expected 1 metric for %s, got %d", metric, len(metrics))
		}

		if metrics[0]["name"] != metric {
			t.Errorf("expected metric name %s, got %v", metric, metrics[0]["name"])
		}
	}
}

// TestCLI_EdgeCases tests CLI with edge case file.
func TestCLI_EdgeCases(t *testing.T) {
	binary := buildCLI(t)
	testFile := filepath.Join("..", "..", "testdata", "edge_cases.md")

	// Should not panic on edge cases
	cmd := exec.Command(binary, "-m", "sentence_length_variance", testFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI failed on edge cases: %v\n%s", err, output)
	}

	if len(output) == 0 {
		t.Error("expected some output for edge cases")
	}
}

// TestCLI_EmptyStdin tests handling of empty stdin input.
func TestCLI_EmptyStdin(t *testing.T) {
	binary := buildCLI(t)

	cmd := exec.Command(binary, "-m", "paragraph_variance")
	cmd.Stdin = bytes.NewReader([]byte{})

	output, err := cmd.CombinedOutput()
	// Should not panic, even if metrics are zero
	_ = err // May or may not error depending on implementation

	// Just verify we got some output (even if it's zeros)
	if len(output) == 0 {
		t.Error("expected some output even for empty input")
	}
}

// TestCLI_ConfigCommands tests config subcommands.
func TestCLI_ConfigCommands(t *testing.T) {
	binary := buildCLI(t)

	// Test config list
	cmd := exec.Command(binary, "config", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("config list failed: %v\n%s", err, output)
	}

	// Should contain config parameters
	outputStr := string(output)
	if !strings.Contains(outputStr, "redundancy_threshold") {
		t.Error("expected config list to show redundancy_threshold")
	}
}

// TestCLI_Help tests help output.
func TestCLI_Help(t *testing.T) {
	binary := buildCLI(t)

	cmd := exec.Command(binary, "-m", "help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("help command failed: %v\n%s", err, output)
	}

	// Should list available metrics
	outputStr := string(output)
	if !strings.Contains(outputStr, "paragraph_variance") {
		t.Error("expected help to list paragraph_variance metric")
	}
	if !strings.Contains(outputStr, "Available metrics") {
		t.Error("expected help to show 'Available metrics' section")
	}
}

// TestCLI_AnalyseSubcommand tests the explicit "analyse" subcommand.
func TestCLI_AnalyseSubcommand(t *testing.T) {
	binary := buildCLI(t)
	testFile := filepath.Join("..", "..", "testdata", "argumentative_essay.md")

	cmd := exec.Command(binary, "analyse", "-m", "paragraph_variance", testFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI analyse failed: %v\n%s", err, output)
	}

	if !strings.Contains(string(output), "paragraph_variance") {
		t.Errorf("expected output to contain metric name, got: %s", output)
	}
}

// TestCLI_AnalyzeSpelling tests the US spelling "analyze" also works.
func TestCLI_AnalyzeSpelling(t *testing.T) {
	binary := buildCLI(t)
	testFile := filepath.Join("..", "..", "testdata", "argumentative_essay.md")

	cmd := exec.Command(binary, "analyze", "-m", "paragraph_variance", testFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI analyze failed: %v\n%s", err, output)
	}

	if !strings.Contains(string(output), "paragraph_variance") {
		t.Errorf("expected output to contain metric name, got: %s", output)
	}
}

// TestCLI_SignatureText tests the signature subcommand with text output.
func TestCLI_SignatureText(t *testing.T) {
	binary := buildCLI(t)
	testFile := filepath.Join("..", "..", "testdata", "argumentative_essay.md")

	cmd := exec.Command(binary, "signature", testFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI signature failed: %v\n%s", err, output)
	}

	outputStr := string(output)

	// Should have 10 bullet points
	bulletCount := strings.Count(outputStr, "•")
	if bulletCount != 10 {
		t.Errorf("expected 10 bullet points, got %d\noutput: %s", bulletCount, outputStr)
	}

	// Should contain all axis names
	for _, name := range []string{"Formality", "Confidence", "Rhythm", "Economy", "Precision",
		"Coherence", "Vocabulary", "Stance", "Emotional Tone", "Temporal Orientation"} {
		if !strings.Contains(outputStr, name) {
			t.Errorf("expected output to contain %q", name)
		}
	}
}

// TestCLI_SignatureJSON tests the signature subcommand with JSON output.
func TestCLI_SignatureJSON(t *testing.T) {
	binary := buildCLI(t)
	testFile := filepath.Join("..", "..", "testdata", "technical_writing.md")

	cmd := exec.Command(binary, "signature", "-format", "json", testFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI signature JSON failed: %v\n%s", err, output)
	}

	// Parse JSON to verify it's valid
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("failed to parse signature JSON: %v\noutput: %s", err, output)
	}

	// Should have version, document, and signature fields
	if _, ok := result["version"]; !ok {
		t.Error("JSON should have 'version' field")
	}
	if _, ok := result["document"]; !ok {
		t.Error("JSON should have 'document' field")
	}
	if _, ok := result["signature"]; !ok {
		t.Error("JSON should have 'signature' field")
	}
}

// TestCLI_SignatureStdin tests the signature subcommand reading from stdin.
func TestCLI_SignatureStdin(t *testing.T) {
	binary := buildCLI(t)
	testText := "This is a test sentence. This is another test sentence. A third one follows."

	cmd := exec.Command(binary, "signature")
	cmd.Stdin = strings.NewReader(testText)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI signature stdin failed: %v\n%s", err, output)
	}

	bulletCount := strings.Count(string(output), "•")
	if bulletCount != 10 {
		t.Errorf("expected 10 bullet points, got %d", bulletCount)
	}
}

// TestCLI_SignatureMultipleFiles tests signature with multiple files.
func TestCLI_SignatureMultipleFiles(t *testing.T) {
	binary := buildCLI(t)
	file1 := filepath.Join("..", "..", "testdata", "argumentative_essay.md")
	file2 := filepath.Join("..", "..", "testdata", "narrative_prose.md")

	cmd := exec.Command(binary, "signature", file1, file2)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI signature multiple files failed: %v\n%s", err, output)
	}

	outputStr := string(output)

	// Should contain both filenames as headers
	if !strings.Contains(outputStr, "argumentative_essay.md") {
		t.Error("expected output to contain first filename")
	}
	if !strings.Contains(outputStr, "narrative_prose.md") {
		t.Error("expected output to contain second filename")
	}

	// Should have 20 bullet points (10 per file)
	bulletCount := strings.Count(outputStr, "•")
	if bulletCount != 20 {
		t.Errorf("expected 20 bullet points for 2 files, got %d", bulletCount)
	}
}
