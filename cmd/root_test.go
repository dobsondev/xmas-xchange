package cmd

import (
	"log/slog"
	"testing"
)

func TestGetSlogLevel_Debug(t *testing.T) {
	level := getSlogLevel("DEBUG")
	if level != slog.LevelDebug {
		t.Errorf("Expected LevelDebug, got %v", level)
	}
}

func TestGetSlogLevel_Info(t *testing.T) {
	level := getSlogLevel("INFO")
	if level != slog.LevelInfo {
		t.Errorf("Expected LevelInfo, got %v", level)
	}
}

func TestGetSlogLevel_Warn(t *testing.T) {
	level := getSlogLevel("WARN")
	if level != slog.LevelWarn {
		t.Errorf("Expected LevelWarn, got %v", level)
	}
}

func TestGetSlogLevel_Error(t *testing.T) {
	level := getSlogLevel("ERROR")
	if level != slog.LevelError {
		t.Errorf("Expected LevelError, got %v", level)
	}
}

func TestGetSlogLevel_CaseInsensitive(t *testing.T) {
	testCases := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"Debug", slog.LevelDebug},
		{"DeBuG", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"INFO", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"Warn", slog.LevelWarn},
		{"error", slog.LevelError},
		{"ErRoR", slog.LevelError},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			level := getSlogLevel(tc.input)
			if level != tc.expected {
				t.Errorf("For input '%s', expected %v, got %v", tc.input, tc.expected, level)
			}
		})
	}
}

func TestGetSlogLevel_DefaultToInfo(t *testing.T) {
	testCases := []string{
		"",
		"invalid",
		"TRACE",
		"WARNING",
		"random",
	}

	for _, input := range testCases {
		t.Run(input, func(t *testing.T) {
			level := getSlogLevel(input)
			if level != slog.LevelInfo {
				t.Errorf("For input '%s', expected default LevelInfo, got %v", input, level)
			}
		})
	}
}
