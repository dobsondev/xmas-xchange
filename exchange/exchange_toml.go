package exchange

import (
	"os"

	"github.com/BurntSushi/toml"
)

type ExchangeConfig struct {
	Exchanges []Exchange `toml:"exchange"`
}

// EncodeExchangeToml marshals exchanges to TOML bytes.
func EncodeExchangeToml(exchanges []Exchange) ([]byte, error) {
	return toml.Marshal(ExchangeConfig{Exchanges: exchanges})
}

// DecodeExchangeToml unmarshals TOML bytes into exchanges.
func DecodeExchangeToml(data []byte) ([]Exchange, error) {
	var config ExchangeConfig
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config.Exchanges, nil
}

// ReadExchangeToml loads exchanges from a local TOML file previously saved by EncodeExchangeToml.
func ReadExchangeToml(path string) ([]Exchange, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return DecodeExchangeToml(data)
}
