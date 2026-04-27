package cmd

import (
	"fmt"
	"os"

	config "github.com/iamllcoolray/stone-cli/internal/configuration"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stone",
	Short: "stone — utiLITI application manager",
	Long:  "stone checks for and installs updates for utiLITI application manager.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		skip := map[string]bool{
			"init":   true,
			"config": true,
			"remove": true,
		}
		if skip[cmd.Name()] {
			return nil
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		return cfg.Validate()
	},
}

func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
