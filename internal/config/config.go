package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

// Generated from https://mholt.github.io/json-to-go/
type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() Config {
	// Reads the json file found at ~/.gatorconfig.json and returns a Config struct.
	path, err := getConfigFilePath()
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(path, configFileName)

	dat, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = json.Unmarshal(dat, &config)
	if err != nil {
		log.Fatalf("config file has invalid format: %v", err)
	}

	return config
}

func (c Config) SetUser(user string) {
	// Writes the config struct to the JSON file after setting the current_user_name field.
	c.CurrentUserName = user

	if err := write(c); err != nil {
		log.Fatal(err)
	}
}

// helper functions
func getConfigFilePath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path, nil

}

func write(config Config) error {
	// Write the file back
	dat, err := json.Marshal(config)
	if err != nil {
		return err
	}

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	path = filepath.Join(path, configFileName)

	return os.WriteFile(path, dat, 0644)
}
