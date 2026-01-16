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

func NewExchange(participants []Participant) Exchange {
	exchange := Exchange{
		GiftingInfo: []GiftExchange{},
	}
	for _, giver := range participants {
		exchange = createSingleExchange(giver, participants, exchange)
	}
	return exchange
}

func createSingleExchange(giver Participant, participants []Participant, exchange Exchange) Exchange {
	receiver := determineReciever(giver, participants, exchange)
	giftExchange := GiftExchange{
		Giver:    giver.Name,
		Receiver: receiver.Name,
	}
	exchange.GiftingInfo = append(exchange.GiftingInfo, giftExchange)
	return exchange
}

func determineReciever(giver Participant, participants []Participant, exchange Exchange) Participant {
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

	randomIndex := rand.Intn(len(potentialReceivers))
	receiver := potentialReceivers[randomIndex]
	return receiver
}

func PrintExchange(exchange Exchange) {
	for _, giftExchange := range exchange.GiftingInfo {
		fmt.Printf("Giver: %s -> Receiver: %s\n", giftExchange.Giver, giftExchange.Receiver)
	}
}
