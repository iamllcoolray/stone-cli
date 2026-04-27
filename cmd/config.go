package cmd

import (
	config "github.com/iamllcoolray/stone-cli/internal/configuration"
	"github.com/spf13/cobra"
)

var (
	setInstallPath string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or update stone configuration",
	Long:  "Without flags prints the current config. Use flags to update individual fields.",
	RunE:  runConfig,
}

func init() {
	configCmd.Flags().StringVar(&setInstallPath, "install-path", "", "Set the local install path")
	rootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) error {
	config.Load()
	return nil
}

func printField(key, value string) {}
