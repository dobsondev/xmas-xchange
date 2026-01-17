package exchange

import (
	"testing"

	. "github.com/dobsondev/xmas-xchange/participants"
)

func TestNewExchange_MultipleAttempts(t *testing.T) {
	// Create a realistic set of participants with restrictions
	participants := []Participant{
		{Name: "Alice", Number: "+15551000001", Restrictions: []string{"Bob"}},
		{Name: "Bob", Number: "+15551000002", Restrictions: []string{}},
		{Name: "Carol", Number: "+15551000003", Restrictions: []string{"Alice"}},
		{Name: "Dave", Number: "+15551000004", Restrictions: []string{"Eve"}},
		{Name: "Eve", Number: "+15551000005", Restrictions: []string{"Dave"}},
		{Name: "Frank", Number: "+15551000006", Restrictions: []string{"Grace"}},
		{Name: "Grace", Number: "+15551000007", Restrictions: []string{"Frank"}},
	}

	// Run 15 exchanges to ensure the retry logic works consistently
	numAttempts := 15
	for i := 0; i < numAttempts; i++ {
		t.Run("Exchange attempt", func(t *testing.T) {
			// This should not panic
			exchange := NewExchange(participants, 50)

			// Verify basic constraints
			if len(exchange.GiftingInfo) != len(participants) {
				t.Errorf("Expected %d gift exchanges, got %d", len(participants), len(exchange.GiftingInfo))
			}

			// Verify everyone is assigned as a giver
			givers := make(map[string]bool)
			for _, gift := range exchange.GiftingInfo {
				givers[gift.Giver] = true
			}
			if len(givers) != len(participants) {
				t.Errorf("Expected %d unique givers, got %d", len(participants), len(givers))
			}

			// Verify everyone is assigned as a receiver
			receivers := make(map[string]bool)
			for _, gift := range exchange.GiftingInfo {
				receivers[gift.Receiver] = true
			}
			if len(receivers) != len(participants) {
				t.Errorf("Expected %d unique receivers, got %d", len(participants), len(receivers))
			}

			// Verify no one gives to themselves
			for _, gift := range exchange.GiftingInfo {
				if gift.Giver == gift.Receiver {
					t.Errorf("Participant %s is giving to themselves", gift.Giver)
				}
			}

			// Verify restrictions are respected
			restrictionMap := make(map[string][]string)
			for _, p := range participants {
				restrictionMap[p.Name] = p.Restrictions
			}

			for _, gift := range exchange.GiftingInfo {
				for _, restricted := range restrictionMap[gift.Giver] {
					if gift.Receiver == restricted {
						t.Errorf("Participant %s is giving to restricted person %s", gift.Giver, gift.Receiver)
					}
				}
			}
		})
	}
}

func TestNewExchange_ImpossibleConstraints(t *testing.T) {
	// Create participants where it's impossible to satisfy constraints
	participants := []Participant{
		{Name: "John", Number: "+15559000001", Restrictions: []string{"Jane"}},
		{Name: "Jane", Number: "+15559000002", Restrictions: []string{"John"}},
	}

	// This should panic after max attempts
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for impossible constraints, but none occurred")
		}
	}()

	NewExchange(participants, 10)
}

func TestNewExchange_NoRestrictions(t *testing.T) {
	// Simple case with no restrictions
	participants := []Participant{
		{Name: "Mike", Number: "+15558000001", Restrictions: []string{}},
		{Name: "Nancy", Number: "+15558000002", Restrictions: []string{}},
		{Name: "Oscar", Number: "+15558000003", Restrictions: []string{}},
	}

	exchange := NewExchange(participants, 50)

	if len(exchange.GiftingInfo) != len(participants) {
		t.Errorf("Expected %d gift exchanges, got %d", len(participants), len(exchange.GiftingInfo))
	}
}
