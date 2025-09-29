package dataprocessor

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"importer/customerimporter"
)

// TODO: add tests for a correct error logging

func TestExportData(t *testing.T) {
	// Sample data for testing
	data := []customerimporter.DomainData{
		{Domain: "example.com", CustomerQuantity: 100},
		{Domain: "test.org", CustomerQuantity: 50},
	}

	t.Run("terminal output", func(t *testing.T) {
		// Capture stdout
		var buf bytes.Buffer
		// Redirect stdout to buffer
		stdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Run ExportData with empty outFile
		ExportData(data, "")

		// Restore stdout and read buffer
		w.Close()
		os.Stdout = stdout
		_, _ = buf.ReadFrom(r)

		// Expected output
		expected := "domain,number_of_customers\nexample.com,100\ntest.org,50\n"
		if got := buf.String(); got != expected {
			t.Errorf("Expected output:\n%s\nGot:\n%s", expected, got)
		}
	})

	t.Run("file output", func(t *testing.T) {
		// Create a temporary file
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "output.csv")

		// Run ExportData with file path
		ExportData(data, tmpFile)

		// Read the file
		content, err := os.ReadFile(tmpFile)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		// Expected CSV content
		expected := "domain,number_of_customers\nexample.com,100\ntest.org,50\n"
		if got := string(content); got != expected {
			t.Errorf("Expected file content:\n%s\nGot:\n%s", expected, got)
		}
	})
}
