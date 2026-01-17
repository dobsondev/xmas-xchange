package exchange

import (
	"fmt"
	"log/slog"
	"math/rand"
	"strings"

	. "github.com/dobsondev/xmas-xchange/participants"
)

type Exchange struct {
	GiftingInfo []GiftExchange
}

type GiftExchange struct {
	Giver    string
	Receiver string
}

func NewExchange(participants []Participant, maxAttempts int) Exchange {
	var err error
	var attempts int = 1

	for attempts <= maxAttempts {
		exchange := Exchange{
			GiftingInfo: []GiftExchange{},
		}

		success := true
		for _, giver := range participants {
			exchange, err = createSingleExchange(giver, participants, exchange)
			if err != nil {
				slog.Debug("Failed to create exchange, trying again...", "error", err, "attempt", attempts, "max attempts", maxAttempts)
				attempts++
				success = false
				break
			}
		}

		if success {
			return exchange
		}
	}

	panic("Could not create a valid exchange after multiple attempts. Please check participant restrictions!")
}

func createSingleExchange(giver Participant, participants []Participant, exchange Exchange) (Exchange, error) {
	receiver, err := determineReciever(giver, participants, exchange)
	if err != nil {
		return Exchange{}, err
	}

	giftExchange := GiftExchange{
		Giver:    giver.Name,
		Receiver: receiver.Name,
	}
	exchange.GiftingInfo = append(exchange.GiftingInfo, giftExchange)
	return exchange, nil
}

func determineReciever(giver Participant, participants []Participant, exchange Exchange) (Participant, error) {
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
	for _, giftExchange := range exchange.GiftingInfo {
		badNames = append(badNames, giftExchange.Receiver)
		slog.Debug("Already receiving", "giver", giftExchange.Giver, "receiver", giftExchange.Receiver)
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
func PrintExchange(exchange Exchange) {
	for _, giftExchange := range exchange.GiftingInfo {
		fmt.Printf("Giver: %s -> Receiver: %s\n", giftExchange.Giver, giftExchange.Receiver)
	}
}
