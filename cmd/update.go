package cmd

import (
	"fmt"

	"github.com/iamllcoolray/stone-cli/internal/api"
	config "github.com/iamllcoolray/stone-cli/internal/configuration"
	"github.com/iamllcoolray/stone-cli/internal/updater"
	"github.com/spf13/cobra"
)

var (
	forceUpdate bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update utiLITI to the latest version",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().BoolVar(&forceUpdate, "force", false, "Force update even if already on latest version")
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
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

	// skip version check if already up to date unless --force
	if cfg.LastVersion != "" && latest == cfg.LastVersion && !forceUpdate {
		fmt.Printf("Already up to date (%s)\n", latest)
		return nil
	}

	if cfg.LastVersion != "" {
		fmt.Printf("Update available: %s → %s\n", cfg.LastVersion, latest)
	} else {
		fmt.Printf("Installing utiLITI %s...\n", latest)
	}

	upload, err := client.FetchPlatformUpload()
	if err != nil {
		return fmt.Errorf("fetching upload: %w", err)
	}

	downloadURL, err := client.FetchDownloadURL(upload.ID)
	if err != nil {
		return fmt.Errorf("fetching download url: %w", err)
	}

	u := updater.New(cfg.InstallPath, client.HTTPClient())
	if err := u.Run(downloadURL, latest); err != nil {
		return err
	}

	cfg.LastVersion = latest
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	return nil
}
