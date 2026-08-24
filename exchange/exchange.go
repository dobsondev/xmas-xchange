package exchange

import (
	"fmt"
	"log/slog"
	"math/rand"
	"strings"

	. "github.com/dobsondev/xmas-xchange/participants"
)

type Exchange struct {
	Giver    Participant `toml:"giver"`
	Receiver Participant `toml:"receiver"`
}

func NewExchange(participants []Participant, maxAttempts int) []Exchange {
	var err error
	var attempts int = 1

	for attempts <= maxAttempts {
		exchanges := []Exchange{}

		success := true
		for _, giver := range participants {
			exchanges, err = createSingleExchange(giver, participants, exchanges)
			if err != nil {
				slog.Debug("Failed to create exchange, trying again...", "error", err, "attempt", attempts, "max attempts", maxAttempts)
				attempts++
				success = false
				break
			}
		}

		if success {
			return exchanges
		}
	}

	panic("Could not create a valid exchange after multiple attempts. Please check participant restrictions!")
}

func createSingleExchange(giver Participant, participants []Participant, exchanges []Exchange) ([]Exchange, error) {
	receiver, err := determineReceiver(giver, participants, exchanges)
	if err != nil {
		return nil, err
	}

	exchange := Exchange{
		Giver:    giver,
		Receiver: receiver,
	}
	exchanges = append(exchanges, exchange)
	return exchanges, nil
}

func determineReceiver(giver Participant, participants []Participant, exchanges []Exchange) (Participant, error) {
	potentialReceivers := participants

	// Can't give to themselves
	badNames := []string{giver.Name}
	slog.Debug("Self restriction", "giver", giver.Name)
	// Can't give to a restricted person
	for _, restrictedNames := range giver.Restrictions {
		badNames = append(badNames, restrictedNames)
		slog.Debug("Restriction", "giver", giver.Name, "restricted", restrictedNames)
	}
	// Can't give to someone who's already receiving a gift
	for _, exchange := range exchanges {
		badNames = append(badNames, exchange.Receiver.Name)
		slog.Debug("Already receiving", "giver", exchange.Giver.Name, "receiver", exchange.Receiver.Name)
	}

	// Filter out bad names
	filteredReceivers := []Participant{}
	for _, participant := range potentialReceivers {
		isBad := false
		for _, badName := range badNames {
			if strings.EqualFold(participant.Name, badName) {
				isBad = true
				break
			}
		}
		if !isBad {
			filteredReceivers = append(filteredReceivers, participant)
			slog.Debug("Valid receiver", "giver", giver.Name, "receiver", participant.Name)
		}
	}
	potentialReceivers = filteredReceivers

	// Check if there are no valid receivers
	if len(potentialReceivers) == 0 {
		slog.Debug("No valid receiver found", "giver", giver.Name)
		return Participant{}, fmt.Errorf("no valid receiver found for %s", giver.Name)
	}

	randomIndex := rand.Intn(len(potentialReceivers))
	receiver := potentialReceivers[randomIndex]
	return receiver, nil
}

// PrintExchange outputs the gift exchange assignments to stdout.
// This function is not unit tested as it only performs console output formatting.
func PrintExchange(exchanges []Exchange) {
	for _, exchange := range exchanges {
		fmt.Printf("Giver: %s -> Receiver: %s\n", exchange.Giver.Name, exchange.Receiver.Name)
	}
}
