package dataprocessor

// import (
// 	"bytes"
// 	"container/heap"
// 	"io"
// 	"os"
// 	"path/filepath"
// 	"testing"
//
// 	"importer/customerimporter"
// )
//
// func TestExportData(t *testing.T) {
// 	// Sample data for testing
// 	data := customerimporter.PriorityQueue{
// 		&customerimporter.DomainData{Domain: "example.com", CustomerQuantity: 100, Index: 1},
// 		&customerimporter.DomainData{Domain: "test.org", CustomerQuantity: 50, Index: 0},
// 	}
// 	heap.Init(&data)
//
// 	t.Run("terminal output", func(t *testing.T) {
// 		// Capture stdout
// 		oldStdout := os.Stdout
// 		r, w, err := os.Pipe()
// 		if err != nil {
// 			t.Fatalf("Failed to create pipe: %v", err)
// 		}
// 		os.Stdout = w
//
// 		// Run ExportData with empty outFile
// 		ExportData(data, "")
//
// 		// Restore stdout and close writer
// 		w.Close()
// 		os.Stdout = oldStdout
//
// 		// Read the output
// 		output := new(bytes.Buffer)
// 		_, err = io.Copy(output, r)
// 		if err != nil {
// 			t.Fatalf("Failed to copy output: %v", err)
// 		}
//
// 		// Expected output
// 		expected := "domain,number_of_customers\ntest.org,50\nexample.com,100\n"
// 		if got := output.String(); got != expected {
// 			t.Errorf("Expected output:\n%s\nGot:\n%s", expected, got)
// 		}
// 	})
//
// 	t.Run("file output", func(t *testing.T) {
// 		// Create a temporary file
// 		tmpDir := t.TempDir()
// 		tmpFile := filepath.Join(tmpDir, "output.csv")
//
// 		// Run ExportData with file path
// 		ExportData(data, tmpFile)
//
// 		// Read the file
// 		content, err := os.ReadFile(tmpFile)
// 		if err != nil {
// 			t.Fatalf("Failed to read output file: %v", err)
// 		}
//
// 		// Expected CSV content
// 		expected := "domain,number_of_customers\nexample.com,100\ntest.org,50\n"
// 		if got := string(content); got != expected {
// 			t.Errorf("Expected file content:\n%s\nGot:\n%s", expected, got)
// 		}
// 	})
// }
