package cmd

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "xmas-xchange",
	Short: "Organize Secret Santa / gift exchange events",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logLevel := os.Getenv("LOG_LEVEL")
		level := getSlogLevel(logLevel)
		setSlogLevel(level)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func getSlogLevel(logLevel string) slog.Level {
	var level slog.Level

	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	return level
}

func setSlogLevel(level slog.Level) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	slog.SetDefault(slog.New(handler))
}
