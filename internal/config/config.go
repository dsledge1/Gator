package config

import (
	"encoding/json"
	"fmt"gh
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DB_URL string `json:"db_url"`
	User   string `json:"current_user_name"`
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	result, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(result, &cfg)
	if err != nil {
		return Config{}, err
	}
	fmt.Println(cfg)
	return cfg, nil
}

func (c *Config) SetUser(user string) {
	c.User = user
	write(*c)
}

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	config := filepath.Join(homedir, configFileName)
	return config, nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	json, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, json, 0666)
	if err != nil {
		return err
	}
	return nil
}
