package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

const filePerm = 0644

type Config struct {
	path             string
	TelegramApiToken string `json:"telegram_api_token"`
	Database         string `json:"database"`
	TimeZone         string `json:"timezone"`
}

func JsonLoadFromFile(path string) (*Config, error) {
	var config Config = Config{path: path}
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err = config.Save(); err != nil {
				return nil, err
			}
			return &config, nil
		}
		return nil, err
	}
	defer file.Close()

	configBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if len(configBytes) == 0 {
		return &config, nil
	}

	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, errors.New("invalid JSON in config file: " + err.Error())
	}

	config.Save()
	return &config, nil
}

func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(c.path, data, filePerm)
	if err != nil {
		return err
	}

	return nil
}
