package cmd

import (
	"fmt"

	"github.com/iamllcoolray/stone-cli/internal/api"
	config "github.com/iamllcoolray/stone-cli/internal/configuration"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if a new version of utiLITI is available",
	RunE:  runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func runCheck(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Println("Checking for updates...")

	client := api.New(cfg.APIKey)

	latest, err := client.FetchLatestVersion()
	if err != nil {
		return fmt.Errorf("fetching latest version: %w", err)
	}

	if cfg.LastVersion == "" {
		fmt.Printf("Latest version: %s\n", latest)
		fmt.Println("No version on record — run: stone update")
		return nil
	}

	if latest == cfg.LastVersion {
		fmt.Printf("Already up to date (%s)\n", latest)
		return nil
	}

	fmt.Printf("Update available: %s → %s\n", cfg.LastVersion, latest)
	fmt.Println("Run: stone update")
	return nil
}
