package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config contains the Cloudflare account ID and API token.
// These fields are necessary to authenticate requests to the Cloudflare API.
type Config struct {
	AccountID string `json:"account_id"`
	Token     string `json:"token"`
}

var ConfigFile string // Path to the configuration file
var homeDir string    // User's home directory path
var err error         // Error variable for capturing function errors

// init is automatically called when the package is initialized.
// It sets up the paths for the configuration file based on the OS.
func init() {
	// Determine the user's home directory based on the operating system
	if runtime.GOOS == "windows" {
		homeDir = os.Getenv("USERPROFILE") // Get home directory on Windows
	} else {
		homeDir, err = os.UserHomeDir() // Get home directory on Unix-like systems
		if err != nil {
			panic(err) // Panic if unable to get home directory
		}
	}

	// Set the config file path in the user's home directory
	ConfigFile = filepath.Join(homeDir, ".aiCFtoken.json")
}

// LoadConfig loads the configuration from the user's home directory.
// It returns a Config struct populated with the account ID and token, or an error.
func LoadConfig() (Config, error) {
	// Early check for initialization error
	if err != nil {
		fmt.Printf("To use this service you need a Cloudflare API Token with WorkerAI Read and Edit Access from this Address %s\n", "https://dash.cloudflare.com/profile/api-tokens")
		return Config{}, err // Return error if initialization failed
	}

	// Read the config file
	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		// Return a specific error if the config file does not exist
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, errors.New("config file does not exist")
		}
		return Config{}, err // Return any other read error
	}

	var config Config
	// Parse the JSON config file into the Config struct
	if err = json.Unmarshal(data, &config); err != nil {
		return Config{}, err // Return error if JSON unmarshalling fails
	}

	return config, nil // Return the populated Config struct
}

// SaveConfig saves the given Config struct to the user's home directory.
// It returns an error if saving fails.
func SaveConfig(config Config) error {
	// Check if the required config fields are provided
	if config.AccountID == "" || config.Token == "" {
		return errors.New("config fields cannot be empty")
	}

	// Early check for initialization error
	if err != nil {
		return err // Return error if initialization failed
	}

	// Convert the Config struct to JSON
	data, err := json.Marshal(config)
	if err != nil {
		return err // Return error if JSON marshalling fails
	}

	// Write the JSON data to the config file
	err = os.WriteFile(ConfigFile, data, 0644)
	if err != nil {
		return err // Return error if writing to the file fails
	}

	// On Windows, set the file attribute to hidden
	if runtime.GOOS == "windows" {
		err = WinFSATTR(ConfigFile)
		if err != nil {
			return fmt.Errorf("failed to set file attributes: %w", err) // Return error if setting attributes fails
		}
	}

	return nil // Return nil if saving was successful
}
