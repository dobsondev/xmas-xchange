package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	. "github.com/dobsondev/xmas-xchange/exchange"
	. "github.com/dobsondev/xmas-xchange/participants"
)

func main() {
	// Configure logging based on LOG_LEVEL environment variable
	logLevel := os.Getenv("LOG_LEVEL")
	level := getSlogLevel(logLevel)
	setSlogLevel(level)

	tomlFlagPtr := flag.String("toml", "participants.toml", "Path to the TOML file containing participants' data")
	dryRunFlagPtr := flag.Bool("dry-run", true, "If set, the program will not perform any actions, just simulate")
	maxAttemptsFlagPtr := flag.Int("max-attempts", 50, "Maximum number of attempts to create a valid exchange")
	flag.Parse()

	participants := GetParticipantsFromToml(tomlFlagPtr)
	dryRun := *dryRunFlagPtr

	if dryRun {
		fmt.Println("===\nDry run mode enabled. No sms messages will be sent.\n===")
	}

	maxAttempts := *maxAttemptsFlagPtr
	exchange := NewExchange(participants, maxAttempts)
	PrintExchange(exchange)
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
