package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func printToAllLevels() {
	slog.Debug("test debug")
	slog.Info("test info")
	slog.Warn("test warn")
	slog.Error("test error")
}

func TestLoggerSetup(t *testing.T) {
	levels := make(map[string][]string)
	levels["debug"] = []string{"DEBUG", "INFO", "WARN", "ERROR"}
	levels["info"] = []string{"INFO", "WARN", "ERROR"}
	levels["warn"] = []string{"WARN", "ERROR"}
	levels["error"] = []string{"ERROR"}
	tmpDir := t.TempDir()

	t.Run("correct level", func(t *testing.T) {
		for level, expected := range levels {
			logPath := fmt.Sprintf("%s/test_correct_%s.log", tmpDir, level)
			logFile, readErr := os.Create(logPath)
			if readErr != nil {
				t.Fatalf("Cannot open temp file %s %s\n", logPath, readErr)
			}
			defer logFile.Close()

			SetupLogger(logFile, level)

			printToAllLevels()

			bytes, readErr := os.ReadFile(logPath)
			if readErr != nil {
				t.Fatalf("Could not read the test file %s: %s\n", logPath, readErr)
			}
			file := string(bytes)
			for _, output := range expected {
				if !strings.Contains(file, output) {
					t.Fatalf("Log is missing in the file %s: %s, at the level %s\n", logPath, output, level)
				}
			}
		}
	})

	t.Run("incorrect level", func(t *testing.T) {
		logPath := fmt.Sprintf("%s/test_incorrect.log", tmpDir)
		logFile, readErr := os.Create(logPath)
		if readErr != nil {
			t.Fatalf("Cannot open temp file %s %s\n", logPath, readErr)
		}
		defer logFile.Close()
		SetupLogger(logFile, "incorrect")

		printToAllLevels()

		bytes, readErr := os.ReadFile(logPath)
		if readErr != nil {
			t.Fatalf("Could not read the test file %s: %s\n", logPath, readErr)
		}
		file := string(bytes)
		if strings.Contains(file, "DEBUG") {
			t.Fatalf("Logging into debug detected: %s\n", file)
		}

		for _, output := range levels["info"] {
			if !strings.Contains(file, output) {
				t.Fatalf("Incorrect setup does not lead to printing at the info level: %s\n", file)
			}
		}

	})
}
