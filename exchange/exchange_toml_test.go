package exchange

import (
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"

	. "github.com/dobsondev/xmas-xchange/participants"
)

func testExchanges() []Exchange {
	return []Exchange{
		{
			Giver:    Participant{Name: "Alice", Number: "+15551000001", Restrictions: []string{"Bob"}},
			Receiver: Participant{Name: "Bob", Number: "+15551000002", Restrictions: []string{}},
		},
		{
			Giver:    Participant{Name: "Bob", Number: "+15551000002", Restrictions: []string{}},
			Receiver: Participant{Name: "Alice", Number: "+15551000001", Restrictions: []string{"Bob"}},
		},
	}
}

func TestEncodeExchangeToml_RoundTrip(t *testing.T) {
	exchanges := testExchanges()

	data, err := EncodeExchangeToml(exchanges)
	if err != nil {
		t.Fatalf("EncodeExchangeToml returned error: %v", err)
	}

	var config ExchangeConfig
	if err := toml.Unmarshal(data, &config); err != nil {
		t.Fatalf("toml.Unmarshal returned error: %v", err)
	}

	if !reflect.DeepEqual(exchanges, config.Exchanges) {
		t.Errorf("Round-tripped exchanges do not match original.\nOriginal: %+v\nRead back: %+v", exchanges, config.Exchanges)
	}
}
