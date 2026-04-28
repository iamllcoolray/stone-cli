package cmd

import (
	"fmt"

	config "github.com/iamllcoolray/stone-cli/internal/configuration"
	"github.com/spf13/cobra"
)

var (
	setAPIKey      string
	setInstallPath string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or update stone configuration",
	Long:  "Without flags prints the current config. Use flags to update individual fields.",
	RunE:  runConfig,
}

func init() {
	configCmd.Flags().StringVar(&setAPIKey, "api-key", "", "Set the itch.io API key")
	configCmd.Flags().StringVar(&setInstallPath, "install-path", "", "Set the local install path")
	rootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// no flags — print current config
	if !cmd.Flags().Changed("api-key") && !cmd.Flags().Changed("install-path") {
		configPath, _ := config.Path()
		fmt.Printf("Config file: %s\n\n", configPath)

		if cfg.APIKey == "" && cfg.InstallPath == "" {
			fmt.Println("No config found. Run 'stone init' to get started.")
			return nil
		}

		// mask api key for display
		maskedKey := cfg.APIKey
		if len(maskedKey) > 8 {
			maskedKey = maskedKey[:8] + "..."
		}
		printField("api_key", maskedKey)
		printField("install_path", cfg.InstallPath)
		printField("last_version", cfg.LastVersion)
		return nil
	}

	// apply flag updates
	if cmd.Flags().Changed("api-key") {
		cfg.APIKey = setAPIKey
	}
	if cmd.Flags().Changed("install-path") {
		cfg.InstallPath = setInstallPath
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	fmt.Println("Config updated.")
	fmt.Println()
	maskedKey := cfg.APIKey
	if len(maskedKey) > 8 {
		maskedKey = maskedKey[:8] + "..."
	}
	printField("api_key", maskedKey)
	printField("install_path", cfg.InstallPath)
	printField("last_version", cfg.LastVersion)
	return nil
}

func printField(key, value string) {
	if value == "" {
		fmt.Printf("  %-16s (not set)\n", key)
	} else {
		fmt.Printf("  %-16s %s\n", key, value)
	}
}
