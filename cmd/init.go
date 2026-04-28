package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	config "github.com/iamllcoolray/stone-cli/internal/configuration"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set up stone and install utiLITI",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	existing, _ := config.Load()
	configPath, err := config.Path()
	if err != nil {
		return err
	}

	// warn if overwriting existing config
	if existing.InstallPath != "" || existing.APIKey != "" {
		fmt.Printf("Config already exists at %s\n", configPath)
		fmt.Print("Overwrite? [y/N] ")
		answer, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(answer)) != "y" {
			fmt.Println("Aborted. Use 'stone config' to update individual fields.")
			return nil
		}
	}

	fmt.Println("Initializing stone...")
	fmt.Println()
	fmt.Println("Get your API key at: https://itch.io/user/settings/api-keys")
	fmt.Println()

	apiKey := prompt(reader, "itch.io API key", existing.APIKey)
	if apiKey == "" {
		return fmt.Errorf("api_key cannot be empty")
	}

	defaultPath := existing.InstallPath
	if defaultPath == "" {
		home, _ := os.UserHomeDir()
		defaultPath = home
	}

	installPath := prompt(reader, "utiLITI install path", defaultPath)
	if installPath == "" {
		return fmt.Errorf("install_path cannot be empty")
	}

	cfg := &config.Config{
		APIKey:      apiKey,
		InstallPath: installPath,
		LastVersion: existing.LastVersion,
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	fmt.Println()
	fmt.Printf("Config saved to %s\n", configPath)
	fmt.Println()
	fmt.Println("  api_key      =", cfg.APIKey[:8]+"...")
	fmt.Println("  install_path =", cfg.InstallPath)
	fmt.Println()
	fmt.Println("Run 'stone update' to install utiLITI.")
	return nil
}

func prompt(r *bufio.Reader, question, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", question, defaultVal)
	} else {
		fmt.Printf("%s: ", question)
	}
	input, _ := r.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	return input
}
