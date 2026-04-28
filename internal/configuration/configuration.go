package configuration

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Config struct {
	APIKey      string `mapstructure:"api_key"`
	InstallPath string `mapstructure:"install_path"`
	LastVersion string `mapstructure:"last_version"`
}

// Dir returns the platform-appropriate config directory.
func Dir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		base := os.Getenv("APPDATA")
		if base == "" {
			return "", fmt.Errorf("APPDATA not set")
		}
		return filepath.Join(base, "stone"), nil
	default:
		base := os.Getenv("XDG_CONFIG_HOME")
		if base == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			base = filepath.Join(home, ".config")
		}
		return filepath.Join(base, "stone"), nil
	}
}

// Path returns the full path to the config file.
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.toml"), nil
}

// Load reads the config file into a Config struct.
// Returns an empty Config (not an error) if no file exists yet.
func Load() (*Config, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(dir)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return &cfg, nil
}

// Save writes the Config struct to disk, creating the directory if needed.
func Save(cfg *Config) error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	viper.Set("api_key", cfg.APIKey)
	viper.Set("install_path", cfg.InstallPath)
	viper.Set("last_version", cfg.LastVersion)

	path := filepath.Join(dir, "config.toml")
	if err := viper.WriteConfigAs(path); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}

// Validate checks that required fields are set.
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("api_key is not set — run: stone init")
	}
	if c.InstallPath == "" {
		return fmt.Errorf("install_path is not set — run: stone init")
	}
	return nil
}
