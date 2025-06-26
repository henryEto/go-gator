package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url   string `json:"db_url"`
	Username string `json:"current_user_name"`
}

func Read() (*Config, error) {
	cfg := Config{}
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.Username = userName
	err := cfg.write()
	return err
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(home, configFileName)
	return configPath, nil
}

func (cfg *Config) write() error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}
	return nil
}
