package config

import (
	"encoding/json"
	"io"
)

// New reads the configuration from the file/Reader and parses it into a Config object.
func New(r io.Reader) (Config, error) {
	configBytes, err := io.ReadAll(r)
	if err != nil {
		return Config{}, err
	}

	var newConfig Config
	err = json.Unmarshal(configBytes, &newConfig)
	if err != nil {
		return Config{}, err
	}
	return newConfig, nil
}
