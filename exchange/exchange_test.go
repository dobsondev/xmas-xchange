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
			exchanges := NewExchange(participants, 50)

			// Verify basic constraints
			if len(exchanges) != len(participants) {
				t.Errorf("Expected %d gift exchanges, got %d", len(participants), len(exchanges))
			}

			// Verify everyone is assigned as a giver
			givers := make(map[string]bool)
			for _, gift := range exchanges {
				givers[gift.Giver.Name] = true
			}
			if len(givers) != len(participants) {
				t.Errorf("Expected %d unique givers, got %d", len(participants), len(givers))
			}

			// Verify everyone is assigned as a receiver
			receivers := make(map[string]bool)
			for _, gift := range exchanges {
				receivers[gift.Receiver.Name] = true
			}
			if len(receivers) != len(participants) {
				t.Errorf("Expected %d unique receivers, got %d", len(participants), len(receivers))
			}

			// Verify no one gives to themselves
			for _, gift := range exchanges {
				if gift.Giver.Name == gift.Receiver.Name {
					t.Errorf("Participant %s is giving to themselves", gift.Giver.Name)
				}
			}

			// Verify restrictions are respected
			restrictionMap := make(map[string][]string)
			for _, p := range participants {
				restrictionMap[p.Name] = p.Restrictions
			}

			for _, gift := range exchanges {
				for _, restricted := range restrictionMap[gift.Giver.Name] {
					if gift.Receiver.Name == restricted {
						t.Errorf("Participant %s is giving to restricted person %s", gift.Giver.Name, gift.Receiver.Name)
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

	exchanges := NewExchange(participants, 50)

	if len(exchanges) != len(participants) {
		t.Errorf("Expected %d gift exchanges, got %d", len(participants), len(exchanges))
	}
}

func TestFilterByGiverName(t *testing.T) {
	exchanges := []Exchange{
		{Giver: Participant{Name: "Alice"}, Receiver: Participant{Name: "Bob"}},
		{Giver: Participant{Name: "Bob"}, Receiver: Participant{Name: "Alice"}},
	}

	filtered, err := FilterByGiverName(exchanges, "alice")
	if err != nil {
		t.Fatalf("FilterByGiverName returned error: %v", err)
	}
	if len(filtered) != 1 || filtered[0].Giver.Name != "Alice" {
		t.Errorf("Expected exactly one exchange for Alice, got %+v", filtered)
	}
}

func TestFilterByGiverName_CaseInsensitive(t *testing.T) {
	exchanges := []Exchange{
		{Giver: Participant{Name: "Alice"}, Receiver: Participant{Name: "Bob"}},
	}

	filtered, err := FilterByGiverName(exchanges, "ALICE")
	if err != nil {
		t.Fatalf("FilterByGiverName returned error: %v", err)
	}
	if len(filtered) != 1 {
		t.Errorf("Expected exactly one exchange, got %d", len(filtered))
	}
}

func TestFilterByGiverName_NoMatch(t *testing.T) {
	exchanges := []Exchange{
		{Giver: Participant{Name: "Alice"}, Receiver: Participant{Name: "Bob"}},
	}

	_, err := FilterByGiverName(exchanges, "Carol")
	if err == nil {
		t.Errorf("Expected an error for no match, got nil")
	}
}
