package main

import (
	"flag"
	"fmt"

	. "github.com/dobsondev/xmas-xchange/exchange"
	. "github.com/dobsondev/xmas-xchange/participants"
)

func main() {
	tomlFlagPtr := flag.String("toml", "participants.toml", "Path to the TOML file containing participants' data")
	dryRunFlagPtr := flag.Bool("dry-run", true, "If set, the program will not perform any actions, just simulate")
	flag.Parse()

	participants := GetParticipantsFromToml(tomlFlagPtr)
	dryRun := *dryRunFlagPtr

	if dryRun {
		fmt.Println("===\nDry run mode enabled. No sms messages will be sent.\n===")
	}

	fmt.Printf("Number of participants: %d\n", len(participants))
	PrintParticipants(participants)

	exchange := NewExchange(participants)
	PrintExchange(exchange)
}
