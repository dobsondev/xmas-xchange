package participants

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type ParticipantsConfig struct {
	Participants []Participant `toml:"participant"`
}

type Participant struct {
	Name         string   `toml:"name"`
	Number       string   `toml:"number"`
	Restrictions []string `toml:"restrictions"`
}

func GetParticipantsFromToml(tomlFilePtr *string) []Participant {
	data, err := os.ReadFile(*tomlFilePtr)
	if err != nil {
		log.Fatal(err)
	}

	var participantsConfig ParticipantsConfig
	err = toml.Unmarshal(data, &participantsConfig)
	if err != nil {
		log.Fatal(err)
	}
	var participants []Participant = participantsConfig.Participants
	return participants
}

func PrintParticipants(participants []Participant) {
	for _, person := range participants {
		fmt.Printf("Name: %s, Number: %s, Restrictions: %v\n", person.Name, person.Number, person.Restrictions)
	}
}
